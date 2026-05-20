package setup

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"app/src/model"
	"gorm.io/gorm"
)

func RenderGuidedSetupPageStepwise(db *gorm.DB, refreshCreds func() (kitsuHost, botToken, guildID, webhookURL string), r *http.Request) string {
	lang := currentLang(r)
	diag := localizeSetupDiagnostics(lang, BuildSetupDiagnostics(db, refreshCreds))
	kitsuHost, _, fallbackGuildID, _ := refreshCreds()
	detectedHost := publicKitsuHostnameFromRequest(r, kitsuHost)
	currentStep := guidedCurrentStep(diag)

	projectOptions, defaultProjectID := buildGuidedProjectOptions(db)
	if defaultProjectID == "" {
		defaultProjectID = strings.TrimSpace(diag.AppliedProjectID)
	}
	if defaultProjectID == "" {
		defaultProjectID = strings.TrimSpace(diag.VerifiedProjectID)
	}

	discordServerName := fallbackGuildID
	discordServerDetail := t(lang, "Guild ID を入力するとここに通知先サーバーが表示されます。", "Enter a Guild ID to resolve the Discord Server here.")
	if fallbackGuildID != "" {
		discordServerDetail = "Guild ID: " + fallbackGuildID
	}
	if strings.TrimSpace(diag.Discord.Detail) != "" && diag.Discord.Status == SetupOK {
		discordServerName = firstNonEmpty(diag.Discord.Detail, fallbackGuildID)
	}

	step1Summary := guidedKitsuSummary(lang, diag)
	step2Summary := guidedDiscordSummary(lang, diag)
	step3Summary := guidedProjectSummary(lang, diag)
	step4Summary := t(lang, "Checker / user mapping は最初の通知成功後に追加すれば十分です。", "Checker / user mapping can wait until after the first successful notification.")
	if currentStep < 4 {
		step4Summary = t(lang, "Step 3 が完了するとここが見えるようになります。", "This step appears after Step 3 is complete.")
	}

	step3Mode := "plan"
	if diag.NotificationVerified {
		step3Mode = "done"
	} else if diag.ProjectSetupApplied {
		step3Mode = "test"
	}

	topNotice := diag.NextAction
	if len(diag.Warnings) > 0 {
		topNotice = diag.Warnings[0]
	}

	body := fmt.Sprintf(`
<style>
.guided-shell{display:grid;grid-template-columns:260px minmax(0,1fr);gap:18px;align-items:start}
.guided-nav,.guided-main{display:flex;flex-direction:column;gap:14px}
.guided-overview,.guided-step{padding:16px;border-radius:16px;border:1px solid var(--line);background:rgba(255,255,255,.03)}
.guided-step.active{border-color:rgba(255,141,72,.48);background:rgba(255,141,72,.06)}
.guided-step.done{border-color:rgba(142,207,139,.34);background:rgba(142,207,139,.05)}
.guided-step.locked{opacity:.6}
.guided-head{display:flex;justify-content:space-between;gap:12px;align-items:flex-start;flex-wrap:wrap}
.guided-head h2,.guided-head h3{margin:0}
.guided-kicker{color:var(--muted);margin-top:8px}
.guided-pill{padding:4px 10px;border-radius:999px;background:rgba(255,255,255,.08);font-size:.75rem;font-weight:700}
.guided-nav-list{display:flex;flex-direction:column;gap:10px}
.guided-nav-item{padding:12px;border-radius:14px;border:1px solid var(--line);background:rgba(255,255,255,.02)}
.guided-nav-item.active{border-color:rgba(255,141,72,.45);background:rgba(255,141,72,.06)}
.guided-nav-item.done{border-color:rgba(142,207,139,.34);background:rgba(142,207,139,.05)}
.guided-nav-item.locked{opacity:.58}
.guided-nav-title{display:flex;justify-content:space-between;gap:10px;align-items:flex-start}
.guided-nav-title strong{display:block}
.guided-nav-note{margin-top:6px;color:var(--muted);font-size:.9rem;line-height:1.4}
.guided-note{color:var(--muted);font-size:.92rem;line-height:1.5}
.guided-form{display:grid;gap:10px;max-width:720px}
.guided-form label{display:block;font-size:.78rem;color:var(--muted-2);margin-bottom:6px}
.guided-form input,.guided-form select{width:100%%;padding:10px 12px;border-radius:12px;border:1px solid var(--line);background:rgba(255,255,255,.04);color:var(--text)}
.guided-actions{display:flex;gap:10px;flex-wrap:wrap}
.guided-banner{padding:12px 14px;border-radius:14px;border:1px solid var(--line);background:rgba(255,255,255,.02);margin-top:12px}
.guided-banner.ok{border-color:rgba(142,207,139,.35);background:rgba(142,207,139,.06)}
.guided-banner.warn{border-color:rgba(255,200,80,.3);background:rgba(255,200,80,.06)}
.guided-banner.error{border-color:rgba(255,106,80,.34);background:rgba(255,106,80,.06)}
.guided-subpanel{padding:14px;border-radius:14px;border:1px solid var(--line);background:rgba(255,255,255,.02)}
.guided-subpanel.hidden{display:none}
.guided-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:12px}
.guided-field{padding:12px;border-radius:12px;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06)}
.guided-field strong{display:block;margin-bottom:4px}
.guided-warning-list{margin:10px 0 0;padding-left:18px;color:#ffc850}
.guided-inline{margin-top:10px}
.guided-result{padding:10px 12px;border-radius:12px;margin-top:10px;border:1px solid var(--line);background:rgba(255,255,255,.02)}
.guided-result.ok{border-color:rgba(142,207,139,.35);color:#8ecf8b}
.guided-result.warn{border-color:rgba(255,200,80,.3);color:#ffc850}
.guided-result.error{border-color:rgba(255,106,80,.34);color:#ff6a50}
.guided-server-box{padding:12px;border-radius:12px;border:1px solid var(--line);background:rgba(255,255,255,.02);margin-bottom:12px}
.guided-server-name{font-weight:700}
.guided-server-sub{color:var(--muted);font-size:.88rem;margin-top:4px}
.guided-preview-list{display:grid;gap:8px;margin-top:12px}
.guided-preview-item{padding:12px;border-radius:12px;border:1px solid rgba(255,255,255,.06);background:rgba(255,255,255,.03)}
.guided-preview-item strong{display:block;margin-bottom:4px}
.hidden{display:none!important}
@media (max-width: 980px){.guided-shell{grid-template-columns:1fr}}
</style>
<div class="section-card glass">
  <div class="guided-head">
    <div>
      <h2>%s</h2>
      <p class="guided-kicker">%s</p>
    </div>
    <div style="display:flex;gap:10px;align-items:center;flex-wrap:wrap">
      <div class="guided-pill">%s</div>
      <a class="btn-ghost" style="font-size:.82rem;padding:5px 12px" href="%s">%s</a>
      <a class="btn-ghost" style="font-size:.82rem;padding:5px 12px" href="%s">%s</a>
    </div>
  </div>
</div>
<div class="guided-banner warn">%s</div>
<div class="guided-shell">
  <aside class="guided-nav">
    <div class="guided-overview">
      <div class="guided-head" style="margin-bottom:12px">
        <h3>%s</h3>
        <div class="guided-pill">%s</div>
      </div>
      <div class="guided-nav-list">
        %s
      </div>
      <p class="guided-note" style="margin-top:14px;padding-top:12px;border-top:1px solid rgba(255,255,255,.07);font-size:.8rem">%s</p>
    </div>
  </aside>
  <main class="guided-main">
    %s
  </main>
</div>
<script>
(function(){
  function esc(v){
    return String(v || '').replace(/[&<>"']/g, function(ch){
      return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch];
    });
  }
  function show(el){ if(el) el.classList.remove('hidden'); }
  function hide(el){ if(el) el.classList.add('hidden'); }
  function setResult(id, kind, html){
    var el = document.getElementById(id);
    if(!el) return;
    el.className = 'guided-result ' + kind;
    el.innerHTML = html;
    show(el);
  }
  function reloadSoon(){
    window.setTimeout(function(){ window.location.reload(); }, 500);
  }

  function setSubmitting(btn, loading){
    if(!btn) return;
    btn.disabled = loading;
    btn.style.opacity = loading ? '0.65' : '';
    if(loading){ btn.dataset.orig = btn.textContent; btn.textContent = '...'; }
    else if(btn.dataset.orig){ btn.textContent = btn.dataset.orig; }
  }

  var step1Form = document.getElementById('guidedStep1Form');
  if(step1Form){
    step1Form.addEventListener('submit', async function(ev){
      ev.preventDefault();
      var btn = document.getElementById('step1SubmitBtn');
      setSubmitting(btn, true);
      var payload = {};
      new FormData(step1Form).forEach(function(value, key){ payload[key] = value; });
      var res = await fetch(%q, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      var data = await res.json();
      if(data.authenticated){
        setResult('guidedStep1Result', 'ok', %q);
        reloadSoon();
      } else {
        setSubmitting(btn, false);
        setResult('guidedStep1Result', 'error', esc(data.error || %q));
      }
    });
  }

  var step2Form = document.getElementById('guidedStep2Form');
  if(step2Form){
    step2Form.addEventListener('submit', async function(ev){
      ev.preventDefault();
      var btn = document.getElementById('step2SubmitBtn');
      setSubmitting(btn, true);
      var payload = {};
      new FormData(step2Form).forEach(function(value, key){ payload[key] = value; });
      var res = await fetch(%q, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      var data = await res.json();
      if(data.bot_valid && data.guild_valid && data.permissions && data.permissions.manage_channels && data.permissions.manage_webhooks){
        setResult('guidedStep2Result', 'ok', esc((data.bot_name || 'Bot') + ' / ' + (data.guild_name || payload.guild_id)));
        reloadSoon();
      } else if(data.bot_valid && data.guild_valid){
        setSubmitting(btn, false);
        setResult('guidedStep2Result', 'warn', esc(data.error || %q));
      } else {
        setSubmitting(btn, false);
        setResult('guidedStep2Result', 'error', esc(data.error || %q));
      }
    });
  }

  var previewForm = document.getElementById('guidedProjectPreviewForm');
  var previewPanel = document.getElementById('guidedProjectPreviewPanel');
  var previewResult = document.getElementById('guidedProjectPreviewResult');
  var applyPanel = document.getElementById('guidedProjectApplyPanel');
  var testPanel = document.getElementById('guidedProjectTestPanel');
  var planPanel = document.getElementById('guidedProjectPlanPanel');
  var completePanel = document.getElementById('guidedProjectCompletePanel');
  var selectedProjectField = document.getElementById('guidedProjectSelected');
  var selectedProjectNameField = document.getElementById('guidedProjectSelectedName');
  var selectedGuildField = document.getElementById('guidedProjectSelectedGuild');

  function showPlan(){ show(planPanel); hide(previewPanel); hide(previewResult); hide(applyPanel); hide(testPanel); hide(completePanel); }
  function showPreview(){ hide(planPanel); show(previewPanel); hide(previewResult); hide(applyPanel); hide(testPanel); hide(completePanel); }
  function showApply(){ hide(planPanel); hide(previewPanel); hide(previewResult); show(applyPanel); show(testPanel); hide(completePanel); }
  function showDone(){ hide(planPanel); hide(previewPanel); hide(previewResult); hide(applyPanel); hide(testPanel); show(completePanel); }

  if(previewForm){
    previewForm.addEventListener('submit', async function(ev){
      ev.preventDefault();
      var payload = {template:'cg'};
      new FormData(previewForm).forEach(function(value, key){ payload[key] = value; });
      var res = await fetch(%q, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      var data = await res.json();
      if(!data.ok){
        setResult('guidedProjectPreviewResult', 'error', esc(data.error || %q));
        show(previewPanel);
        return;
      }
      if(selectedProjectField){ selectedProjectField.value = data.project_id || payload.project_id || ''; }
      if(selectedProjectNameField){ selectedProjectNameField.textContent = data.project_name || ''; }
      if(selectedGuildField){ selectedGuildField.textContent = (data.discord_server_name || '') + (data.guild_id ? ' / Guild ID: ' + data.guild_id : ''); }
      var html = '';
      html += '<div class="guided-grid">';
      html += '<div class="guided-field"><strong>%s</strong><div>' + esc(data.project_name || payload.project_id || '') + '</div></div>';
      html += '<div class="guided-field"><strong>%s</strong><div>' + esc(data.discord_server_name || data.guild_id || '') + '</div></div>';
      html += '<div class="guided-field"><strong>%s</strong><div>' + esc(data.category_name || '') + '</div></div>';
      html += '<div class="guided-field"><strong>%s</strong><div>' + esc(String((data.channels_to_create || []).length)) + '</div></div>';
      html += '<div class="guided-field"><strong>%s</strong><div>' + esc(String(data.webhooks_to_create || 0)) + '</div></div>';
      html += '</div>';
      html += '<div class="guided-preview-list">';
      (data.channels_to_create || []).forEach(function(item){
        html += '<div class="guided-preview-item"><strong>#' + esc(item.name) + '</strong><div class="guided-note">' + esc((item.task_types || []).join(', ')) + '</div></div>';
      });
      html += '</div>';
      if(data.warnings && data.warnings.length){
        html += '<ul class="guided-warning-list">' + data.warnings.map(function(item){ return '<li>' + esc(item) + '</li>'; }).join('') + '</ul>';
      }
      html += '<div class="guided-banner warn" style="margin-top:12px"><strong>%s</strong><div class="guided-note" style="margin-top:6px">%s</div></div>';
      html += '<div class="guided-actions" style="margin-top:12px">';
      html += '<button id="guidedConfirmCreateBtn" class="btn" type="button">%s</button>';
      html += '<button id="guidedBackToPlanBtn" class="btn-ghost" type="button">%s</button>';
      html += '</div>';
      previewPanel.innerHTML = html;
      showPreview();
      var confirmBtn = document.getElementById('guidedConfirmCreateBtn');
      var backBtn = document.getElementById('guidedBackToPlanBtn');
      if(backBtn){
        backBtn.addEventListener('click', function(){ showPlan(); });
      }
      if(confirmBtn){
        confirmBtn.addEventListener('click', async function(){
          var applyPayload = {
            project_id: payload.project_id,
            guild_id: data.guild_id || payload.guild_id || '',
            template: 'cg',
            language: payload.language || 'ja'
          };
          var applyRes = await fetch(%q, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(applyPayload)});
          var applyData = await applyRes.json();
          if(applyData.ok){
            var summary = '<div class="guided-grid">';
            summary += '<div class="guided-field"><strong>%s</strong><div>' + esc(applyData.project_name || payload.project_id || '') + '</div></div>';
            summary += '<div class="guided-field"><strong>%s</strong><div>' + esc(applyData.discord_server_name || data.discord_server_name || applyData.guild_id || '') + '</div></div>';
            summary += '<div class="guided-field"><strong>%s</strong><div>' + esc(String(applyData.channels_created || 0)) + '</div></div>';
            summary += '<div class="guided-field"><strong>%s</strong><div>' + esc(String(applyData.webhooks_created || 0)) + '</div></div>';
            summary += '<div class="guided-field"><strong>%s</strong><div>' + esc(applyData.safe_to_retry ? 'yes' : 'no') + '</div></div>';
            summary += '</div>';
            if((applyData.channels_created || 0) !== (data.channels_to_create || []).length || (applyData.webhooks_created || 0) !== (data.webhooks_to_create || 0)){
              summary += '<div class="guided-banner warn">%s</div>';
            }
            if(applyData.warnings && applyData.warnings.length){
              summary += '<ul class="guided-warning-list">' + applyData.warnings.map(function(item){ return '<li>' + esc(item) + '</li>'; }).join('') + '</ul>';
            }
            if(applyData.lines && applyData.lines.length){
              summary += '<details style="margin-top:10px"><summary>%s</summary><pre>' + esc(applyData.lines.join('\n')) + '</pre></details>';
            }
            applyPanel.innerHTML = summary;
            showApply();
            var testProjectInput = document.querySelector('#guidedProjectTestForm input[name="project_id"]');
            if(testProjectInput){ testProjectInput.value = applyData.project_id || payload.project_id || ''; }
            return;
          }
          setResult('guidedProjectApplyResult', 'error', esc(applyData.error || %q));
          hide(planPanel);
          hide(previewPanel);
          hide(previewResult);
          show(applyPanel);
          hide(testPanel);
          hide(completePanel);
        });
      }
    });
  }

  var testForm = document.getElementById('guidedProjectTestForm');
  if(testForm){
    testForm.addEventListener('submit', async function(ev){
      ev.preventDefault();
      var payload = {};
      new FormData(testForm).forEach(function(value, key){ payload[key] = value; });
      var res = await fetch(%q, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      var data = await res.json();
      if(data.ok){
        setResult('guidedProjectTestResult', 'ok', %q + ' ' + esc(data.project_name || payload.project_id || ''));
        if(completePanel){
          showDone();
        }
        reloadSoon();
      } else {
        setResult('guidedProjectTestResult', 'error', esc(data.error || %q));
      }
    });
  }
})();
</script>`,
		esc(t(lang, "Guided Setup", "Guided Setup")),
		esc(t(lang, "最初の通知が届くまでを、1 step ずつ順番に進みます。詳しい診断や編集は Manual Setup に残します。", "This flow walks you step by step to the first delivered notification. Detailed diagnostics stay in Manual Setup.")),
		esc(t(lang, "Current Step", "Current Step")+": "+guidedStepTitle(lang, currentStep)),
		withLang("/bot/admin/setup", r),
		esc(t(lang, "Manual Setup →", "Manual Setup →")),
		withLang("/bot/setup-wizard", r),
		esc(t(lang, "← Entry", "← Entry")),
		esc(topNotice),
		esc(t(lang, "Setup Flow", "Setup Flow")),
		esc(func() string {
			if diag.SetupComplete {
				return "complete"
			}
			return "in progress"
		}()),
		func() string {
			var nav strings.Builder
			for _, c := range diag.Env {
				if c.Status == SetupError {
					nav.WriteString(navStepHTML(lang, 0, currentStep, t(lang, "必要な環境変数が不足しています", "Missing required environment variables"), t(lang, "システム確認", "System Check")))
					break
				}
			}
			nav.WriteString(navStepHTML(lang, 1, currentStep, step1Summary, t(lang, "Kitsu Connection", "Kitsu Connection")))
			nav.WriteString(navStepHTML(lang, 2, currentStep, step2Summary, t(lang, "Discord Bot", "Discord Bot")))
			nav.WriteString(navStepHTML(lang, 3, currentStep, step3Summary, t(lang, "Project Setup", "Project Setup")))
			if currentStep >= 4 {
				nav.WriteString(navStepHTMLOptional(lang, 4, currentStep, step4Summary, t(lang, "Mapping (Optional)", "Mapping (Optional)")))
			}
			return nav.String()
		}(),
		esc(t(lang,
			"Project = Kitsuの制作単位 / Guild = Discordの通知先サーバー",
			"Project = Kitsu production unit / Guild = Discord notification server")),
		renderActiveStepBody(lang, diag, currentStep, detectedHost, discordServerName, discordServerDetail, defaultProjectID, projectOptions, step3Mode, fallbackGuildID, r),
		withLang("/api/setup/test-kitsu", r),
		esc(t(lang, "Kitsu に接続できました。", "Kitsu connection succeeded.")),
		esc(t(lang, "Kitsu 接続に失敗しました。", "Kitsu connection failed.")),
		withLang("/api/setup/test-discord", r),
		esc(t(lang, "Discord Bot と Server を確認できました。", "Discord bot and server are ready.")),
		esc(t(lang, "Discord Bot の確認に失敗しました。", "Discord validation failed.")),
		withLang("/api/setup/preview-project", r),
		esc(t(lang, "preview を作成できませんでした。", "Preview could not be created.")),
		esc(t(lang, "Project", "Project")),
		esc(t(lang, "Discord Server", "Discord Server")),
		esc(t(lang, "Category", "Category")),
		esc(t(lang, "Channels", "Channels")),
		esc(t(lang, "Webhooks", "Webhooks")),
		esc(t(lang, "確認前の注意", "Before you confirm")),
		esc(t(lang, "Preview ではまだ Discord に何も作成していません。次の Confirm and Create で category / channels / webhooks の作成を開始します。Preview 内容を確認してから進んでください。失敗時は best-effort で rollback を試みますが保証ではありません。必要なら残った Discord リソースを手動で確認または削除してから再試行してください。", "The preview has not created anything in Discord yet. The next Confirm and Create action starts creating the category, channels, and webhooks. Review the previewed plan before you continue. If provisioning fails, KitsuSync attempts best-effort rollback, but it is not guaranteed. You may need to manually check or delete leftover Discord resources before retrying.")),
		esc(t(lang, "Confirm and Create", "Confirm and Create")),
		esc(t(lang, "Back", "Back")),
		withLang("/api/setup/apply-project", r),
		esc(t(lang, "Project", "Project")),
		esc(t(lang, "Discord Server", "Discord Server")),
		esc(t(lang, "Channels created", "Channels created")),
		esc(t(lang, "Webhooks created", "Webhooks created")),
		esc(t(lang, "Safe to retry", "Safe to retry")),
		esc(t(lang, "Preview と execute の結果に差がありました。execute の結果を正として表示しています。", "Preview and execute differed. Showing the execute result as the source of truth.")),
		esc(t(lang, "Project setup execution failed.", "Project setup execution failed.")),
		esc(t(lang, "Raw details", "Raw details")),
		withLang("/api/setup/test-notification", r),
		esc(t(lang, "Project Setup complete", "Project Setup complete")),
		esc(t(lang, "Test notification failed.", "Test notification failed.")),
	)

	return adminPage(lang, t(lang, "Guided Setup", "Guided Setup"), r, body)
}

func guidedCurrentStep(diag SetupDiagnostics) int {
	for _, c := range diag.Env {
		if c.Status == SetupError {
			return 0
		}
	}
	if diag.Kitsu.Status != SetupOK {
		return 1
	}
	if diag.Discord.Status != SetupOK {
		return 2
	}
	if !diag.NotificationVerified {
		return 3
	}
	return 4
}

func guidedStepTitle(lang string, step int) string {
	switch step {
	case 0:
		return t(lang, "システム確認", "System Check")
	case 1:
		return t(lang, "Kitsu Connection", "Kitsu Connection")
	case 2:
		return t(lang, "Discord Bot", "Discord Bot")
	case 3:
		return t(lang, "Project Setup", "Project Setup")
	case 4:
		return t(lang, "Mapping (Optional)", "Mapping (Optional)")
	default:
		return t(lang, "Setup", "Setup")
	}
}

// navStepHTMLOptional renders a step nav item with an "optional" pill when not yet reached.
func navStepHTMLOptional(lang string, step, current int, summary, title string) string {
	class := "guided-nav-item locked"
	pill := t(lang, "任意", "optional")
	if step < current {
		class = "guided-nav-item done"
		pill = "done"
	} else if step == current {
		class = "guided-nav-item active"
		pill = "active"
	}
	return fmt.Sprintf(`<div class="%s"><div class="guided-nav-title"><strong>%d. %s</strong><span class="guided-pill" style="opacity:.8">%s</span></div><div class="guided-nav-note">%s</div></div>`,
		class, step, html.EscapeString(title), html.EscapeString(pill), html.EscapeString(summary))
}

func navStepHTML(lang string, step, current int, summary, title string) string {
	class := "guided-nav-item locked"
	pill := t(lang, "前のStepを完了", "complete previous step")
	if step < current {
		class = "guided-nav-item done"
		pill = "✓ done"
	} else if step == current {
		class = "guided-nav-item active"
		pill = "▶ " + t(lang, "進行中", "active")
	}
	return fmt.Sprintf(`<div class="%s"><div class="guided-nav-title"><strong>%d. %s</strong><span class="guided-pill">%s</span></div><div class="guided-nav-note">%s</div></div>`,
		class, step, html.EscapeString(title), html.EscapeString(pill), html.EscapeString(summary))
}

func guidedKitsuSummary(lang string, diag SetupDiagnostics) string {
	switch diag.Kitsu.Status {
	case SetupOK:
		return t(lang, "Kitsu に接続できました。", "Kitsu connected.")
	case SetupWarn:
		return firstNonEmpty(diag.Kitsu.Detail, diag.Kitsu.Fix, diag.Kitsu.Summary)
	case SetupError:
		return firstNonEmpty(diag.Kitsu.Fix, diag.Kitsu.Detail, diag.Kitsu.Summary)
	default:
		return t(lang, "Kitsu 接続を確認します。", "Check Kitsu connection.")
	}
}

func guidedDiscordSummary(lang string, diag SetupDiagnostics) string {
	switch diag.Discord.Status {
	case SetupOK:
		return firstNonEmpty(diag.Discord.Detail, diag.Discord.Summary, t(lang, "Discord Bot は準備できています。", "Discord Bot is ready."))
	case SetupWarn:
		return firstNonEmpty(diag.Discord.Fix, diag.Discord.Detail, diag.Discord.Summary)
	case SetupError:
		return firstNonEmpty(diag.Discord.Fix, diag.Discord.Detail, diag.Discord.Summary)
	default:
		return t(lang, "Discord Bot を確認します。", "Check Discord Bot.")
	}
}

func guidedProjectSummary(lang string, diag SetupDiagnostics) string {
	if diag.NotificationVerified {
		if strings.TrimSpace(diag.VerifiedProjectName) != "" {
			return diag.VerifiedProjectName + ": " + t(lang, "Test notification complete.", "Test notification complete.")
		}
		return t(lang, "Project Setup complete", "Project Setup complete")
	}
	if diag.ProjectSetupApplied {
		if strings.TrimSpace(diag.AppliedProjectName) != "" {
			return diag.AppliedProjectName + ": " + t(lang, "Channels and webhooks are ready.", "Channels and webhooks are ready.")
		}
		return t(lang, "Channels and webhooks are ready. Send the test notification next.", "Channels and webhooks are ready. Send the test notification next.")
	}
	return t(lang, "project を選び、preview してから作成します。", "Pick a project, preview it, then create it.")
}

func buildGuidedProjectOptions(db *gorm.DB) (string, string) {
	projects := ListKitsuProjects("")
	configured := make(map[string]struct{})
	for _, p := range model.ListProjects(db) {
		configured[p.KitsuProjectID] = struct{}{}
	}

	var out strings.Builder
	selected := ""
	for _, p := range projects {
		label := p.Name
		if _, ok := configured[p.ID]; ok {
			label += " (configured)"
		}
		if selected == "" {
			if _, ok := configured[p.ID]; !ok {
				selected = p.ID
			}
		}
		option := fmt.Sprintf(`<option value="%s"%s>%s</option>`, html.EscapeString(p.ID), "", html.EscapeString(label))
		if p.ID == selected {
			option = fmt.Sprintf(`<option value="%s" selected>%s</option>`, html.EscapeString(p.ID), html.EscapeString(label))
		}
		out.WriteString(option)
	}
	if selected == "" && len(projects) > 0 {
		selected = projects[0].ID
	}
	if out.Len() == 0 {
		out.WriteString(`<option value="">` + html.EscapeString("No projects found") + `</option>`)
	}
	return out.String(), selected
}

func renderActiveStepBody(lang string, diag SetupDiagnostics, currentStep int, detectedHost, discordServerName, discordServerDetail, defaultProjectID, projectOptions, step3Mode, fallbackGuildID string, r *http.Request) string {
	switch currentStep {
	case 0:
		var envErrors strings.Builder
		for _, c := range diag.Env {
			if c.Status == SetupError {
				envErrors.WriteString(fmt.Sprintf(
					`<div class="guided-banner error" style="margin-bottom:10px"><strong>%s</strong><p style="margin:6px 0 0;font-weight:normal;font-size:.92rem">%s</p></div>`,
					html.EscapeString(c.Label),
					html.EscapeString(firstNonEmpty(c.Fix, c.Summary)),
				))
			}
		}
		return fmt.Sprintf(`
<section class="guided-step active">
  <div class="guided-head"><div><h3>%s</h3><p class="guided-note">%s</p></div><div class="guided-pill">active</div></div>
  <div style="margin-top:14px">%s</div>
  <div class="guided-banner warn" style="margin:12px 0">%s</div>
  <div class="guided-actions" style="margin-top:14px">
    <a class="btn-ghost" href="%s">%s</a>
    <a class="btn-ghost" href="%s">%s</a>
  </div>
</section>`,
			html.EscapeString(t(lang, "Step 0: システム確認", "Step 0: System Check")),
			html.EscapeString(t(lang, "セットアップを始める前に、必要な環境変数が揃っているか確認します。", "Before starting setup, confirm that required environment variables are in place.")),
			envErrors.String(),
			html.EscapeString(t(lang,
				"修正方法: プロジェクトルートの `.env.local` ファイルを編集し、`docker compose restart` で再起動してください。変更はコンテナ再起動後に反映されます。",
				"Fix: Edit `.env.local` in the project root, then run `docker compose restart`. Changes take effect after the container restarts.",
			)),
			withLang("/bot/admin/diagnostics", r),
			html.EscapeString(t(lang, "診断ページを開く", "Open Diagnostics")),
			withLang("/bot/admin/setup", r),
			html.EscapeString(t(lang, "Manual Setup を開く", "Open Manual Setup")),
		)
	case 1:
		envNotice := ""
		for _, c := range diag.Env {
			if c.Status == SetupError {
				envNotice = c.Fix
				break
			}
		}
		if envNotice != "" {
			envNotice = fmt.Sprintf(`<div class="guided-banner warn">%s</div>`, html.EscapeString(envNotice))
		}
		// Build the host note as raw HTML so we can include a link to Bot Settings.
		step1HostNote := t(lang,
			`ホスト名を変更する場合は <a href="/bot/admin/bot" style="color:var(--accent-2)">Bot設定</a> で編集してください。`,
			`To change the hostname, edit it in <a href="/bot/admin/bot" style="color:var(--accent-2)">Bot Settings</a>.`,
		)
		return fmt.Sprintf(`
<section class="guided-step active">
  <div class="guided-head"><div><h3>%s</h3><p class="guided-note">%s</p></div><div class="guided-pill">%s</div></div>
  %s
  <div class="guided-server-box">
    <div class="guided-server-name">%s</div>
    <div class="guided-server-sub">%s</div>
  </div>
  <form id="guidedStep1Form" class="guided-form">
    <input type="hidden" name="hostname" value="%s">
    <div><label>%s</label><input value="%s" readonly></div>
    <div><label>%s</label><input name="email" placeholder="admin@studio.local" autocomplete="email"></div>
    <div><label>%s</label><input name="password" type="password" autocomplete="current-password"></div>
    <div class="guided-actions"><button class="btn" type="submit" id="step1SubmitBtn">%s</button></div>
  </form>
  <div id="guidedStep1Result" class="guided-result hidden"></div>
</section>`,
			html.EscapeString(t(lang, "Step 1: Kitsu Connection", "Step 1: Kitsu Connection")),
			html.EscapeString(t(lang, "Kitsuのタスク変更を検知するために認証が必要です。メールアドレスとパスワードを入力して接続テストしてください。", "Authentication is required so KitsuSync can detect task changes. Enter your email and password, then test the connection.")),
			html.EscapeString(stepBadgeText(lang, diag.Kitsu.Status)),
			envNotice,
			html.EscapeString(t(lang, "Kitsu ホスト（自動検出）", "Kitsu Host (auto-detected)")),
			step1HostNote,
			html.EscapeString(detectedHost),
			html.EscapeString(t(lang, "Kitsu ホスト", "Kitsu Host")),
			html.EscapeString(detectedHost),
			html.EscapeString(t(lang, "メールアドレス（Kitsu 管理者）", "Email (Kitsu admin)")),
			html.EscapeString(t(lang, "パスワード", "Password")),
			html.EscapeString(t(lang, "接続テスト", "Test Connection")),
		)
	case 2:
		// Build hint as raw HTML so we can embed links.
		guildIDHint := t(lang,
			`Discord の <strong>開発者モード</strong> を有効にし（設定 › 詳細設定）、通知先サーバーのアイコンを右クリック → 「ID をコピー」。`,
			`Enable <strong>Developer Mode</strong> in Discord (Settings › Advanced), right-click the server icon → "Copy Server ID".`,
		)
		botTokenHint := t(lang,
			`このフォームは接続テスト専用です。Bot Token を保存するには <a href="/bot/admin/bot" style="color:var(--accent-2)">Bot設定</a> に入力してください。`,
			`This form is for connection testing only. To save the Bot Token permanently, enter it in <a href="/bot/admin/bot" style="color:var(--accent-2)">Bot Settings</a>.`,
		)
		return fmt.Sprintf(`
<section class="guided-step active">
  <div class="guided-head"><div><h3>%s</h3><p class="guided-note">%s</p></div><div class="guided-pill">%s</div></div>
  <div class="guided-server-box">
    <div class="guided-server-name">%s</div>
    <div class="guided-server-sub">%s</div>
  </div>
  <form id="guidedStep2Form" class="guided-form">
    <div>
      <label>%s</label>
      <input name="bot_token" type="password" placeholder="hidden token" autocomplete="off">
      <p class="guided-note" style="font-size:.8rem;margin-top:4px">%s</p>
    </div>
    <div>
      <label>%s</label>
      <input name="guild_id" value="%s" placeholder="123456789012345678" autocomplete="off">
      <p class="guided-note" style="font-size:.8rem;margin-top:4px">%s</p>
    </div>
    <div class="guided-actions"><button class="btn" type="submit" id="step2SubmitBtn">%s</button></div>
  </form>
  <div id="guidedStep2Result" class="guided-result hidden"></div>
</section>`,
			html.EscapeString(t(lang, "Step 2: Discord Bot", "Step 2: Discord Bot")),
			html.EscapeString(t(lang, "BotトークンとDiscordサーバーIDを入力し、接続・権限を確認します。", "Enter your bot token and Discord server ID to verify the connection and permissions.")),
			html.EscapeString(stepBadgeText(lang, diag.Discord.Status)),
			html.EscapeString(discordServerName),
			html.EscapeString(discordServerDetail),
			html.EscapeString(t(lang, "Bot Token", "Bot Token")),
			botTokenHint,
			html.EscapeString(t(lang, "Discord Server ID（Guild ID）", "Discord Server ID (Guild ID)")),
			html.EscapeString(fallbackGuildID),
			guildIDHint,
			html.EscapeString(t(lang, "接続・権限を確認する", "Validate Discord")),
		)
	case 3:
		return renderGuidedStep3(lang, diag, projectOptions, defaultProjectID, step3Mode, fallbackGuildID, r)
	default:
		return fmt.Sprintf(`
<section class="guided-step active">
  <div class="guided-head"><div><h3>%s</h3><p class="guided-note">%s</p></div><div class="guided-pill">%s</div></div>
  <div class="guided-banner ok">%s</div>
  <div class="guided-actions" style="margin-top:12px">
    <a class="btn-ghost" href="%s">%s</a>
    <a class="btn-ghost" href="%s">%s</a>
  </div>
</section>`,
			html.EscapeString(t(lang, "Step 4: Mapping (Optional)", "Step 4: Mapping (Optional)")),
			html.EscapeString(t(lang, "最初の通知成功後に追加すれば十分です。", "Add mapping later after the first successful notification.")),
			html.EscapeString("done"),
			html.EscapeString(t(lang, "Checker / user mapping は初回セットアップの必須項目ではありません。", "Checker / user mapping is optional for initial setup.")),
			html.EscapeString(withLang("/bot/admin/users", r)),
			html.EscapeString(t(lang, "User Mapping", "User Mapping")),
			html.EscapeString(withLang("/bot/admin/checkers", r)),
			html.EscapeString(t(lang, "Checker Mapping", "Checker Mapping")),
		)
	}
}

func renderGuidedStep3(lang string, diag SetupDiagnostics, projectOptions, defaultProjectID, step3Mode, fallbackGuildID string, r *http.Request) string {
	projectLabel := defaultProjectID
	if diag.ProjectSetupApplied && strings.TrimSpace(diag.AppliedProjectName) != "" {
		projectLabel = diag.AppliedProjectName
	}
	if diag.NotificationVerified && strings.TrimSpace(diag.VerifiedProjectName) != "" {
		projectLabel = diag.VerifiedProjectName
	}
	_ = projectLabel
	serverLabel := t(lang, "Discord Server", "Discord Server")
	serverSub := t(lang, "Guild ID は preview で確定します。", "Guild ID will be resolved during preview.")
	if diag.ProjectSetupApplied && !diag.NotificationVerified {
		serverSub = t(lang, "channels / webhooks は作成済みです。テスト通知を送ると Step 3 が完了します。", "Channels and webhooks are ready. Send a test notification to finish Step 3.")
	}
	if diag.NotificationVerified {
		serverSub = t(lang, "テスト通知まで完了しています。", "The test notification has been delivered.")
	}

	planHidden := "hidden"
	testHidden := "hidden"
	doneHidden := "hidden"
	if step3Mode == "plan" {
		planHidden = ""
	} else if step3Mode == "test" {
		testHidden = ""
	} else if step3Mode == "done" {
		doneHidden = ""
	}
	_ = testHidden
	_ = doneHidden

	return fmt.Sprintf(`
<section class="guided-step active">
  <div class="guided-head">
    <div>
      <h3>%s</h3>
      <p class="guided-note">%s</p>
    </div>
    <div class="guided-pill">%s</div>
  </div>
  <div class="guided-server-box">
    <div class="guided-server-name">%s</div>
    <div class="guided-server-sub">%s</div>
  </div>
  <div id="guidedProjectPlanPanel" class="%s">
    <p class="guided-note">%s</p>
    <form id="guidedProjectPreviewForm" class="guided-form">
      <input type="hidden" name="template" value="cg">
      <div>
        <label>%s</label>
        <select name="project_id">%s</select>
      </div>
      <div>
        <label>%s</label>
        <input name="guild_id" value="%s" placeholder="123456789012345678">
      </div>
      <div>
        <label>%s</label>
        <select name="language">
          <option value="ja">%s</option>
          <option value="en">%s</option>
        </select>
      </div>
      <input type="hidden" name="mode" value="create">
      <div class="guided-actions">
        <button class="btn" type="submit">%s</button>
      </div>
      <p class="guided-note" style="font-size:.82rem;margin-top:4px">%s</p>
    </form>
  </div>
  <div id="guidedProjectPreviewResult" class="guided-result hidden"></div>
  <div id="guidedProjectPreviewPanel" class="guided-subpanel hidden"></div>
  <div id="guidedProjectApplyPanel" class="guided-subpanel %s">
    <div class="guided-head"><div><h3>%s</h3></div><div class="guided-pill">%s</div></div>
    <p class="guided-note">%s</p>
    <div id="guidedProjectApplyResult" class="guided-result hidden"></div>
  </div>
  <div id="guidedProjectTestPanel" class="guided-subpanel %s">
    <div class="guided-head"><div><h3>%s</h3></div><div class="guided-pill">%s</div></div>
    <p class="guided-note">%s</p>
    <form id="guidedProjectTestForm" class="guided-form">
      <input type="hidden" name="project_id" value="%s">
      <div><label>%s</label><input name="message" value="KitsuSync test notification"></div>
      <div class="guided-actions"><button class="btn" type="submit">%s</button></div>
    </form>
    <div id="guidedProjectTestResult" class="guided-result hidden"></div>
  </div>
  <div id="guidedProjectCompletePanel" class="guided-subpanel %s">
    <div class="guided-head"><div><h3>%s</h3></div><div class="guided-pill">%s</div></div>
    <p class="guided-note">%s</p>
    <div class="guided-banner ok">%s</div>
    <div class="guided-actions" style="margin-top:12px">
      <a class="btn" href="%s">%s</a>
      <a class="btn-ghost" href="%s">%s</a>
    </div>
  </div>
</section>`,
		html.EscapeString(t(lang, "Step 3: Project Setup", "Step 3: Project Setup")),
		html.EscapeString(t(lang, "project を選び、Discord Server を preview して、確認後にだけ作成します。", "Pick a project, preview the Discord Server, and only create it after confirmation.")),
		html.EscapeString(step3Mode),
		html.EscapeString(serverLabel),
		html.EscapeString(serverSub),
		planHidden,
		html.EscapeString(t(lang, "この段階ではまだ Discord に何も作成しません。", "Nothing is created in Discord yet.")),
		html.EscapeString(t(lang, "Project", "Project")),
		projectOptions,
		html.EscapeString(t(lang, "Discord Server", "Discord Server")),
		html.EscapeString(fallbackGuildID),
		html.EscapeString(t(lang, "Language", "Language")),
		html.EscapeString(t(lang, "日本語", "Japanese")),
		html.EscapeString(t(lang, "英語", "English")),
		html.EscapeString(t(lang, "Preview Setup", "Preview Setup")),
		html.EscapeString(t(lang, "「Confirm and Create」を押すまでは Discord に何も作成されません。Preview はいつでもキャンセルできます。", `Nothing is created in Discord until you click "Confirm and Create". You can cancel the preview at any time.`)),
		func() string {
			if step3Mode == "plan" {
				return "hidden"
			}
			return ""
		}(),
		html.EscapeString(t(lang, "Preview", "Preview")),
		html.EscapeString("ok"),
		html.EscapeString(t(lang, "作成前の preview を確認できます。", "Review the preview before creating anything.")),
		func() string {
			if step3Mode == "test" {
				return ""
			}
			return "hidden"
		}(),
		html.EscapeString(t(lang, "Test Notification", "Test Notification")),
		html.EscapeString("ok"),
		html.EscapeString(t(lang, "Project Setup を apply した後にここで 1 回だけ送ります。", "After applying project setup, send one test message here.")),
		html.EscapeString(firstNonEmpty(diag.AppliedProjectID, defaultProjectID)),
		html.EscapeString(t(lang, "Message", "Message")),
		html.EscapeString(t(lang, "テスト通知を送る", "Send Test Notification")),
		func() string {
			if step3Mode == "done" {
				return ""
			}
			return "hidden"
		}(),
		html.EscapeString(t(lang, "Project Setup complete", "Project Setup complete")),
		html.EscapeString("done"),
		html.EscapeString(t(lang, "テスト通知が完了しました。Step 3 は完了です。", "Test notification delivered — Step 3 is complete.")),
		html.EscapeString(t(lang, "テスト通知成功！セットアップが完了しました。", "Test notification succeeded! Setup is complete.")),
		withLang("/bot/admin", r),
		html.EscapeString(t(lang, "このまま完了する →", "Done — Go to Admin →")),
		withLang("/bot/admin/setup", r),
		html.EscapeString(t(lang, "Manual Setup", "Manual Setup")),
	)
}

func stepBadgeText(lang string, status SetupStatus) string {
	switch status {
	case SetupOK:
		return "done"
	case SetupWarn:
		return "warn"
	case SetupError:
		return "error"
	default:
		return "in progress"
	}
}
