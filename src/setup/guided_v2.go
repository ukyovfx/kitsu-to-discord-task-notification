package setup

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func RenderGuidedSetupPageV2(db *gorm.DB, refreshCreds func() (kitsuHost, botToken, guildID, webhookURL string), r *http.Request) string {
	lang := currentLang(r)
	diag := localizeSetupDiagnostics(lang, BuildSetupDiagnostics(db, refreshCreds))
	kitsuHost, botToken, fallbackGuildID, _ := refreshCreds()
	autoKitsuHost := publicKitsuHostnameFromRequest(r, kitsuHost)
	esc := html.EscapeString

	configuredByID := make(map[string]ProjectSetupStatus, len(diag.Projects))
	for _, p := range diag.Projects {
		configuredByID[p.ProjectID] = p
	}

	var projectOptions strings.Builder
	for _, p := range ListKitsuProjects("") {
		label := p.Name
		if _, ok := configuredByID[p.ID]; ok {
			label += " (" + t(lang, "設定済み", "configured") + ")"
		}
		projectOptions.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, esc(p.ID), esc(label)))
	}
	if projectOptions.Len() == 0 {
		projectOptions.WriteString(`<option value="">` + esc(t(lang, "プロジェクトが見つかりません", "No projects found")) + `</option>`)
	}

	serverLabel := t(lang, "Preview で解決されます", "Resolved in preview")
	if strings.TrimSpace(fallbackGuildID) != "" {
		serverLabel = fallbackGuildID
		if info := checkDiscordStatus(botToken, fallbackGuildID); strings.TrimSpace(info.GuildName) != "" {
			serverLabel = strings.TrimSpace(info.GuildName)
		}
	}

	step3StatusText := t(lang, "まだ project setup は未実行です。", "Project setup has not been run yet.")
	step3StatusClass := "guided-warning"
	if diag.ProjectSetupApplied && !diag.NotificationVerified {
		step3StatusText = t(lang, "channel / webhook の準備は完了しています。テスト通知を送ると Step 3 が完了します。", "Channels and webhooks are ready. Send a test notification to complete Step 3.")
	}
	if diag.NotificationVerified {
		step3StatusClass = "guided-success"
		step3StatusText = t(lang, "テスト通知まで完了しています。", "The test notification has been delivered.")
	}

	body := fmt.Sprintf(`
<style>
.guided-shell{display:grid;grid-template-columns:260px minmax(0,1fr);gap:18px}
.guided-nav,.guided-main{display:flex;flex-direction:column;gap:14px}
.step-card,.guided-step{padding:16px;border-radius:16px;border:1px solid var(--line);background:rgba(255,255,255,.03)}
.guided-step.done{border-color:rgba(142,207,139,.35);background:rgba(142,207,139,.05)}
.guided-step.active{border-color:rgba(255,141,72,.45);background:rgba(255,141,72,.06)}
.step-title{display:flex;justify-content:space-between;gap:10px;align-items:flex-start;margin-bottom:10px}
.step-title h3{margin:0}
.step-list{display:flex;flex-direction:column;gap:8px}
.step-pill{font-size:.75rem;padding:4px 8px;border-radius:999px;background:rgba(255,255,255,.08)}
.guided-form{display:grid;gap:12px;max-width:680px}
.guided-form input,.guided-form select{width:100%%;padding:10px 12px;border-radius:12px;border:1px solid var(--line);background:rgba(255,255,255,.04);color:var(--text)}
.guided-form input[readonly]{opacity:.86}
.guided-form label{font-size:.78rem;color:var(--muted-2)}
.guided-note{color:var(--muted);font-size:.92rem}
.guided-actions{display:flex;gap:10px;flex-wrap:wrap;margin-top:8px}
.guided-success{color:#8ecf8b;font-weight:700}
.guided-warning{color:#ffc850;font-weight:700}
.guided-error{color:#ff6a50;font-weight:700}
.server-box,.preview-box,.summary-box{padding:14px;border-radius:14px;border:1px solid var(--line);background:rgba(255,255,255,.02)}
.server-name{font-weight:700}
.server-sub,.summary-sub{color:var(--muted);font-size:.88rem}
.preview-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:12px;margin-top:12px}
.preview-list{display:grid;gap:8px;margin-top:12px}
.preview-item{padding:12px;border-radius:12px;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06)}
.preview-item strong{display:block;margin-bottom:4px}
.warning-list{margin:10px 0 0;padding-left:18px;color:#ffc850}
.inline-status{margin-top:10px}
.optional-note{font-size:.84rem;color:var(--muted)}
.hidden{display:none!important}
@media (max-width:980px){.guided-shell{grid-template-columns:1fr}}
</style>
<div class="section-card glass">
  <h2 style="margin:0 0 8px">%s</h2>
  <p class="guided-note">%s</p>
</div>
<div class="guided-shell">
  <div class="guided-nav">
    <div class="step-card">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <div class="step-list">
        <div class="guided-step %s"><strong>1. %s</strong><div class="guided-note">%s</div></div>
        <div class="guided-step %s"><strong>2. %s</strong><div class="guided-note">%s</div></div>
        <div class="guided-step %s"><strong>3. %s</strong><div class="guided-note">%s</div></div>
        <div class="guided-step"><strong>4. %s</strong><div class="guided-note">%s</div></div>
      </div>
    </div>
  </div>
  <div class="guided-main">
    <div class="guided-step %s">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <p class="guided-note">%s</p>
      <div class="server-box" style="margin-bottom:12px">
        <div class="server-name">%s</div>
        <div class="server-sub">%s</div>
      </div>
      <form id="guidedKitsuForm" class="guided-form">
        <input type="hidden" name="hostname" value="%s">
        <div><label>%s</label><input value="%s" readonly></div>
        <div><label>%s</label><input name="email" placeholder="admin@studio.local"></div>
        <div><label>%s</label><input name="password" type="password"></div>
        <div class="guided-actions"><button class="btn" type="submit">%s</button></div>
      </form>
      <div id="guidedKitsuResult" class="inline-status"></div>
    </div>

    <div class="guided-step %s">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <p class="guided-note">%s</p>
      <form id="guidedDiscordForm" class="guided-form">
        <div><label>%s</label><input name="bot_token" type="password" placeholder="hidden token"></div>
        <div><label>%s</label><input name="guild_id" value="%s" placeholder="123456789012345678"></div>
        <div class="guided-actions"><button class="btn" type="submit">%s</button></div>
      </form>
      <div id="guidedDiscordResult" class="inline-status"></div>
    </div>

    <div class="guided-step %s">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <p class="guided-note">%s</p>
      <p class="guided-note">%s</p>
      <div class="%s">%s</div>
      <form id="guidedPreviewForm" class="guided-form">
        <div><label>%s</label><select name="project_id">%s</select></div>
        <div>
          <label>%s</label>
          <div class="server-box">
            <div id="resolvedServerName" class="server-name">%s</div>
            <div id="resolvedServerMeta" class="server-sub">%s</div>
          </div>
          <input type="hidden" name="guild_id" value="">
        </div>
        <div><label>%s</label><select name="language"><option value="ja">%s</option><option value="en">%s</option></select></div>
        <div><label>%s</label><select name="mode"><option value="create">%s</option><option value="reuse" disabled>%s</option></select></div>
        <div class="guided-actions"><button class="btn" type="submit">%s</button></div>
        <div class="optional-note">%s</div>
      </form>
      <div id="guidedPreviewPanel" class="preview-box hidden"></div>
      <div id="guidedExecutePanel" class="summary-box hidden"></div>
      <div id="guidedTestPanel" class="%s">
        <div class="step-title" style="margin-top:14px"><h3>%s</h3><span class="step-pill">%s</span></div>
        <p class="guided-note">%s</p>
        <form id="guidedTestForm" class="guided-form">
          <input type="hidden" name="project_id" value="%s">
          <div><label>%s</label><input name="message" value="KitsuSync test notification"></div>
          <div class="guided-actions"><button class="btn" type="submit">%s</button></div>
        </form>
        <div id="guidedTestResult" class="inline-status"></div>
      </div>
    </div>

    <div class="guided-step %s">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <p class="guided-note">%s</p>
      <div class="guided-actions">
        <a class="btn-ghost" href="%s">%s</a>
        <a class="btn-ghost" href="%s">%s</a>
      </div>
    </div>

    <div class="guided-step %s">
      <div class="step-title"><h3>%s</h3><span class="step-pill">%s</span></div>
      <p class="guided-note">%s</p>
      <div id="guidedCompleteText" class="%s">%s</div>
      <div class="guided-actions">
        <a class="btn" href="%s">%s</a>
        <a class="btn-ghost" href="%s">%s</a>
      </div>
    </div>
  </div>
</div>
<script>
(function(){
  var endpoints = {
    testKitsu: %q,
    testDiscord: %q,
    previewProject: %q,
    applyProject: %q,
    testNotification: %q
  };

  function setResult(id, kind, html){
    var el = document.getElementById(id);
    if(!el) return;
    el.className = 'inline-status ' + kind;
    el.innerHTML = html;
  }

  function escapeHtml(value){
    return String(value || '').replace(/[&<>"']/g, function(ch){
      return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch];
    });
  }

  function renderWarnings(list){
    if(!list || !list.length) return '';
    return '<ul class="warning-list">' + list.map(function(item){
      return '<li>' + escapeHtml(item) + '</li>';
    }).join('') + '</ul>';
  }

  function projectNameFromSelect(){
    var sel = document.querySelector('#guidedPreviewForm select[name="project_id"]');
    if(!sel) return '';
    var opt = sel.options[sel.selectedIndex];
    return opt ? opt.textContent : '';
  }

  document.getElementById('guidedKitsuForm').addEventListener('submit', async function(ev){
    ev.preventDefault();
    var payload = {};
    new FormData(this).forEach(function(value, key){ payload[key] = value; });
    var res = await fetch(endpoints.testKitsu, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    var data = await res.json();
    if(data.authenticated){
      setResult('guidedKitsuResult', 'guided-success', %q);
    } else {
      setResult('guidedKitsuResult', 'guided-error', escapeHtml(data.error || %q));
    }
  });

  document.getElementById('guidedDiscordForm').addEventListener('submit', async function(ev){
    ev.preventDefault();
    var payload = {};
    new FormData(this).forEach(function(value, key){ payload[key] = value; });
    var res = await fetch(endpoints.testDiscord, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    var data = await res.json();
    if(data.bot_valid && data.guild_valid){
      setResult('guidedDiscordResult', 'guided-success', escapeHtml((data.bot_name || 'Bot') + ' / ' + (data.guild_name || payload.guild_id)));
    } else {
      setResult('guidedDiscordResult', 'guided-error', escapeHtml(data.error || %q));
    }
  });

  document.getElementById('guidedPreviewForm').addEventListener('submit', async function(ev){
    ev.preventDefault();
    var payload = {template:'cg'};
    new FormData(this).forEach(function(value, key){ payload[key] = value; });
    var res = await fetch(endpoints.previewProject, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    var data = await res.json();
    var executePanel = document.getElementById('guidedExecutePanel');
    if(!data.ok){
      setResult('guidedExecutePanel', 'guided-error', escapeHtml(data.error || %q));
      executePanel.classList.remove('hidden');
      return;
    }

    document.getElementById('resolvedServerName').textContent = data.discord_server_name || %q;
    document.getElementById('resolvedServerMeta').textContent = 'Guild ID: ' + (data.guild_id || '-');
    document.querySelector('#guidedPreviewForm input[name="guild_id"]').value = data.guild_id || '';

    var html = '<div class="server-name">' + escapeHtml(data.project_name || projectNameFromSelect()) + '</div>';
    html += '<div class="server-sub">' + escapeHtml(data.discord_server_name || '') + ' / Guild ID: ' + escapeHtml(data.guild_id || '') + '</div>';
    html += '<div class="preview-grid">';
    html += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(data.category_name || '') + '</div></div>';
    html += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(String((data.channels_to_create || []).length)) + '</div></div>';
    html += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(String(data.webhooks_to_create || 0)) + '</div></div>';
    html += '</div>';
    html += '<div class="preview-list">' + (data.channels_to_create || []).map(function(item){
      return '<div class="preview-item"><strong>#' + escapeHtml(item.name) + '</strong><div class="summary-sub">' + escapeHtml((item.task_types || []).join(', ')) + '</div></div>';
    }).join('') + '</div>';
    html += renderWarnings(data.warnings);
    html += '<div class="guided-actions"><button class="btn" type="button" id="confirmCreateBtn">%s</button><button class="btn-ghost" type="button" id="backPreviewBtn">%s</button></div>';

    var panel = document.getElementById('guidedPreviewPanel');
    panel.innerHTML = html;
    panel.classList.remove('hidden');

    document.getElementById('confirmCreateBtn').addEventListener('click', async function(){
      var applyPayload = {
        project_id: payload.project_id,
        guild_id: data.guild_id || '',
        template: 'cg',
        language: payload.language || 'ja'
      };
      var applyRes = await fetch(endpoints.applyProject, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(applyPayload)});
      var applyData = await applyRes.json();
      var summary = '<div class="server-name">' + escapeHtml(applyData.project_name || '') + '</div>';
      summary += '<div class="server-sub">' + escapeHtml(applyData.discord_server_name || data.discord_server_name || '') + ' / Guild ID: ' + escapeHtml(applyData.guild_id || data.guild_id || '') + '</div>';
      if(applyData.ok){
        summary += '<div class="preview-grid">';
        summary += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(String(applyData.channels_created || 0)) + '</div></div>';
        summary += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(String(applyData.webhooks_created || 0)) + '</div></div>';
        summary += '<div><strong>%s</strong><div class="summary-sub">' + escapeHtml(applyData.safe_to_retry ? 'yes' : 'no') + '</div></div>';
        summary += '</div>';
        if((applyData.channels_created || 0) !== (data.channels_to_create || []).length || (applyData.webhooks_created || 0) !== (data.webhooks_to_create || 0)){
          summary += '<div class="guided-warning" style="margin-top:10px">%s</div>';
        }
        summary += renderWarnings(applyData.warnings);
        document.querySelector('#guidedTestForm input[name="project_id"]').value = applyData.project_id || payload.project_id;
        document.getElementById('guidedTestPanel').classList.remove('hidden');
      } else {
        summary += '<div class="guided-error" style="margin-top:10px">' + escapeHtml(applyData.error || %q) + '</div>';
      }
      summary += '<details style="margin-top:10px"><summary>%s</summary><pre>' + escapeHtml((applyData.lines || []).join('\n')) + '</pre></details>';
      executePanel.innerHTML = summary;
      executePanel.classList.remove('hidden');
    });

    document.getElementById('backPreviewBtn').addEventListener('click', function(){
      panel.classList.add('hidden');
      panel.innerHTML = '';
    });
  });

  document.getElementById('guidedTestForm').addEventListener('submit', async function(ev){
    ev.preventDefault();
    var payload = {};
    new FormData(this).forEach(function(value, key){ payload[key] = value; });
    var res = await fetch(endpoints.testNotification, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    var data = await res.json();
    if(data.ok){
      setResult('guidedTestResult', 'guided-success', %q + ' ' + escapeHtml(data.project_name || payload.project_id));
      document.getElementById('guidedCompleteText').className = 'guided-success';
      document.getElementById('guidedCompleteText').textContent = %q;
    } else {
      setResult('guidedTestResult', 'guided-error', escapeHtml(data.error || %q));
    }
  });
})();
</script>`,
		esc(t(lang, "Guided Setup", "Guided Setup")),
		esc(t(lang, "最初の通知が届くところまでを順番に進めます。詳しい診断や編集は Manual Setup に残します。", "This flow walks you to the first successful notification. Detailed diagnostics and editing stay in Manual Setup.")),
		esc(t(lang, "Setup Flow", "Setup Flow")),
		esc(func() string { if diag.SetupComplete { return "complete" } ; return "in progress" }()),
		stepClassFromStatus(diag.Kitsu.Status), esc(t(lang, "Kitsu 接続", "Kitsu Connection")), esc(t(lang, "検出済み host で接続確認", "Check the detected host and test authentication.")),
		stepClassFromStatus(diag.Discord.Status), esc(t(lang, "Discord Bot", "Discord Bot")), esc(t(lang, "Bot と Discord Server の確認", "Validate the bot and Discord Server access.")),
		stepClassFromBool(diag.NotificationVerified), esc(t(lang, "Project Setup", "Project Setup")), esc(t(lang, "preview -> confirm -> execute -> test notification", "preview -> confirm -> execute -> test notification")),
		esc(t(lang, "Mapping (Optional)", "Mapping (Optional)")), esc(t(lang, "最初の通知後に追加できます", "Add mentions later after the first notification.")),
		stepClassFromStatus(diag.Kitsu.Status), esc(t(lang, "Step 1: Kitsu 接続確認", "Step 1: Kitsu Connection")), esc(t(lang, "Detected Host", "Detected Host")), esc(t(lang, "Kitsu Hostname は自動解決されています。変更が必要な場合は Manual Setup で編集できます。", "Kitsu Hostname is auto-detected. If it needs to change, edit it in Manual Setup.")), esc(t(lang, "Detected Kitsu Host", "Detected Kitsu Host")), esc(t(lang, "変更が必要な場合は Manual Setup で編集できます。", "If this needs to change, edit it in Manual Setup.")), esc(autoKitsuHost), esc(t(lang, "Detected Kitsu Host", "Detected Kitsu Host")), esc(autoKitsuHost), esc(t(lang, "メールアドレス", "Email")), esc(t(lang, "パスワード", "Password")), esc(t(lang, "接続テスト", "Test Kitsu")),
		stepClassFromStatus(diag.Discord.Status), esc(t(lang, "Step 2: Discord Bot 確認", "Step 2: Discord Bot")), esc(stepBadgeText(lang, diag.Discord.Status)), esc(t(lang, "Bot token と Discord Server への参加状態を確認します。", "Confirm the bot token and Discord Server access.")), esc(t(lang, "Bot Token", "Bot Token")), esc(t(lang, "Guild ID", "Guild ID")), esc(fallbackGuildID), esc(t(lang, "確認する", "Validate Discord")),
		func() string { if diag.NotificationVerified || diag.ProjectSetupApplied { return "guided-step done" } ; return "guided-step active" }(), esc(t(lang, "Step 3: Project Setup", "Step 3: Project Setup")), esc(t(lang, "Project Setup", "Project Setup")), esc(t(lang, "project を選び、通知先の Discord Server を preview し、確認後にだけ作成します。", "Pick a project, preview the Discord Server changes, and only create them after confirmation.")), esc(t(lang, "project = Kitsu側の制作単位 / Discord Server = 通知先サーバー", "project = Kitsu unit of work / Discord Server = notification target server")), step3StatusClass, esc(step3StatusText), esc(t(lang, "Project", "Project")), projectOptions.String(), esc(t(lang, "Discord Server", "Discord Server")), esc(serverLabel), esc(t(lang, "Preview で確定した Guild ID がここに表示されます。", "The resolved Guild ID will appear here after preview.")), esc(t(lang, "Language", "Language")), esc(t(lang, "日本語", "Japanese")), esc(t(lang, "英語", "English")), esc(t(lang, "Setup mode", "Setup mode")), esc(t(lang, "新規作成", "Create new channels")), esc(t(lang, "既存チャンネルを使う (v0.1.x 予定)", "Reuse existing channels (coming in v0.1.x)")), esc(t(lang, "Preview Setup", "Preview Setup")), esc(t(lang, "この段階ではまだ Discord に何も作成しません。", "This step is read-only. Nothing is created in Discord yet.")), func() string { if diag.ProjectSetupApplied { return "" } ; return "hidden" }(), esc(t(lang, "Test Notification", "Test Notification")), esc(func() string { if diag.NotificationVerified { return "ok" } ; return "pending" }()), esc(t(lang, "Send one message only after project setup has been applied.", "Send one message only after project setup has been applied.")), esc(firstNonEmpty(diag.AppliedProjectID, diag.VerifiedProjectID)), esc(t(lang, "Message", "Message")), esc(t(lang, "テスト通知を送る", "Send Test Notification")),
		stepClassFromBool(diag.NotificationVerified), esc(t(lang, "Step 4: Mapping は Optional", "Step 4: Mapping is Optional")), esc(t(lang, "Optional", "Optional")), esc(t(lang, "Checker / user mapping は最初の通知成功後に追加すれば十分です。", "Checker and user mapping can wait until after the first successful notification.")), withLang("/bot/admin/users", r), esc(t(lang, "User Mapping", "User Mapping")), withLang("/bot/admin/checkers", r), esc(t(lang, "Checker Mapping", "Checker Mapping")),
		stepClassFromBool(diag.NotificationVerified), esc(t(lang, "Complete", "Complete")), esc(func() string { if diag.NotificationVerified { return "done" } ; return "pending" }()), esc(t(lang, "Project Setup complete は test notification 成功後のみ表示されます。", "Project Setup complete only appears after a successful test notification.")), func() string { if diag.NotificationVerified { return "guided-success" } ; return "guided-warning" }(), esc(func() string { if diag.NotificationVerified { return t(lang, "Project Setup complete", "Project Setup complete") } ; return t(lang, "まだ Step 3 は完了していません", "Step 3 is not complete yet") }()), withLang("/bot/admin/setup", r), esc(t(lang, "Open Manual Setup", "Open Manual Setup")), withLang("/bot/admin", r), esc(t(lang, "Open Admin", "Open Admin")),
		withLang("/api/setup/test-kitsu", r), withLang("/api/setup/test-discord", r), withLang("/api/setup/preview-project", r), withLang("/api/setup/apply-project", r), withLang("/api/setup/test-notification", r),
		esc(t(lang, "Kitsu に接続できました。", "Kitsu connection succeeded.")), esc(t(lang, "Kitsu 認証に失敗しました。", "Kitsu authentication failed.")), esc(t(lang, "Discord Bot の確認に失敗しました。", "Discord validation failed.")), esc(t(lang, "preview の取得に失敗しました。", "Preview could not be created.")), esc(t(lang, "Resolved in preview", "Resolved in preview")), esc(t(lang, "Category", "Category")), esc(t(lang, "Channels", "Channels")), esc(t(lang, "Webhooks", "Webhooks")), esc(t(lang, "Confirm and Create", "Confirm and Create")), esc(t(lang, "Back", "Back")), esc(t(lang, "Channels created", "Channels created")), esc(t(lang, "Webhooks created", "Webhooks created")), esc(t(lang, "Safe to retry", "Safe to retry")), esc(t(lang, "preview と execute の結果に差分がありました。実行結果を正として表示しています。", "Preview and execute differed. Showing the execute result as the source of truth.")), esc(t(lang, "project setup の実行に失敗しました。", "Project setup execution failed.")), esc(t(lang, "Raw details", "Raw details")), esc(t(lang, "テスト通知に成功しました:", "Test notification delivered for:")), esc(t(lang, "Project Setup complete", "Project Setup complete")), esc(t(lang, "テスト通知の送信に失敗しました。", "Test notification failed.")),
	)

	return adminPage(lang, t(lang, "Guided Setup", "Guided Setup"), r, body)
}
