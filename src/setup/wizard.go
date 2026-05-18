package setup

import (
	"fmt"
	"net/http"
)

// WizardHandler handles GET /bot/setup-wizard.
// Renders a 4-step setup wizard backed by /api/setup/* endpoints.
func WizardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		lang := currentLang(r)
		fmt.Fprint(w, renderWizard(lang, r))
	}
}

func renderWizard(lang string, r *http.Request) string {
	apiBase := "/api/setup"
	adminURL := withLang("/bot/admin", r)
	setupURL := withLang("/bot/setup", r)

	body := fmt.Sprintf(`
<style>
.wiz-steps{display:flex;align-items:center;gap:0;margin-bottom:28px;overflow:hidden;border-radius:var(--radius-md);border:1px solid var(--line)}
.wiz-step-tab{flex:1;padding:12px 8px;text-align:center;font-size:.8rem;font-weight:600;color:var(--muted-2);background:var(--panel-soft);transition:all .2s;white-space:nowrap}
.wiz-step-tab.active{color:var(--accent-2);background:rgba(255,141,72,.1);border-bottom:2px solid var(--accent-2)}
.wiz-step-tab.done{color:var(--success);background:rgba(142,207,139,.07);border-bottom:2px solid var(--success)}
.wiz-step-tab.done::before{content:"✓ "}
.wiz-panel{display:none;animation:riseIn .2s ease}
.wiz-panel.active{display:block}
.status-badge{display:inline-flex;align-items:center;gap:8px;padding:8px 14px;border-radius:var(--radius-sm);font-size:.88rem;font-weight:600;margin-bottom:16px}
.badge-ok{background:rgba(142,207,139,.14);border:1px solid rgba(142,207,139,.35);color:var(--success)}
.badge-err{background:rgba(255,106,80,.13);border:1px solid rgba(255,106,80,.35);color:var(--danger)}
.badge-warn{background:rgba(255,200,80,.1);border:1px solid rgba(255,200,80,.3);color:#ffc850}
.badge-info{background:rgba(255,255,255,.05);border:1px solid var(--line);color:var(--muted)}
.wiz-form{display:flex;flex-direction:column;gap:14px;max-width:520px}
.wiz-form label{font-size:.78rem;color:var(--muted-2);margin-bottom:3px;display:block}
.wiz-form input,.wiz-form select{width:100%%;padding:10px 14px;background:rgba(255,255,255,.04);border:1px solid var(--line);border-radius:var(--radius-sm);color:var(--text);font-size:.9rem;transition:border-color .15s}
.wiz-form input:focus,.wiz-form select:focus{outline:none;border-color:var(--accent-2)}
.wiz-actions{display:flex;gap:10px;align-items:center;flex-wrap:wrap;margin-top:6px}
.project-item{display:flex;align-items:center;gap:8px;padding:10px 0;border-bottom:1px solid var(--line)}
.project-item:last-child{border-bottom:none}
.project-badge{font-size:.72rem;padding:2px 8px;border-radius:20px;background:rgba(142,207,139,.14);border:1px solid rgba(142,207,139,.35);color:var(--success)}
.result-box{margin-top:16px;padding:14px 18px;border-radius:var(--radius-sm);background:rgba(255,255,255,.03);border:1px solid var(--line);font-size:.84rem}
.result-box ul{margin:6px 0;padding-left:18px;color:var(--muted)}
.result-box li.ok{color:var(--success)}
.result-box li.fail{color:var(--danger);font-weight:600}
.result-box li.warn{color:#ffc850}
.spinner{display:inline-block;width:14px;height:14px;border:2px solid rgba(255,255,255,.2);border-top-color:var(--accent-2);border-radius:50%%;animation:spin .6s linear infinite;vertical-align:middle;margin-left:6px}
@keyframes spin{to{transform:rotate(360deg)}}
#step3ProjectList{margin:14px 0}
.map-row{display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid rgba(255,255,255,.05)}
.map-row:last-child{border-bottom:none}
.map-person{flex:1;min-width:0}
.map-name{font-size:.88rem;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.map-email{font-size:.74rem;color:var(--muted-2);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.map-input{width:200px;padding:6px 10px;background:rgba(255,255,255,.04);border:1px solid var(--line);border-radius:var(--radius-sm);color:var(--text);font-size:.84rem}
.map-input:focus{outline:none;border-color:var(--accent-2)}
.map-section{margin-bottom:24px}
.map-section-title{font-size:.76rem;font-weight:700;color:var(--muted-2);text-transform:uppercase;letter-spacing:.05em;margin-bottom:10px;padding-bottom:6px;border-bottom:1px solid var(--line)}
.map-task{flex:1;font-size:.88rem;font-weight:600}
.optional-notice{display:flex;align-items:flex-start;gap:10px;padding:12px 16px;border-radius:var(--radius-sm);background:rgba(255,200,80,.07);border:1px solid rgba(255,200,80,.25);color:var(--muted);font-size:.84rem;margin-bottom:20px;line-height:1.5}
.optional-icon{font-size:1.1rem;flex-shrink:0;margin-top:1px}
</style>

<!-- Step tabs -->
<div class="wiz-steps">
  <div class="wiz-step-tab" id="tab1">1. Kitsu</div>
  <div class="wiz-step-tab" id="tab2">2. Discord</div>
  <div class="wiz-step-tab" id="tab3">3. %s</div>
  <div class="wiz-step-tab" id="tab4">4. %s</div>
</div>

<!-- Step 1: Kitsu -->
<div class="wiz-panel" id="panel1">
  <div class="section-card glass" style="max-width:600px">
    <h3>%s</h3>
    <div id="step1Status"></div>
    <div id="step1Form" class="wiz-form">
      <div>
        <label>Kitsu Hostname</label>
        <input id="kitsuHostname" type="url" placeholder="http://kitsu.studio.local/">
      </div>
      <div>
        <label>%s</label>
        <input id="kitsuEmail" type="email" placeholder="bot@studio.com" autocomplete="username">
      </div>
      <div>
        <label>%s</label>
        <input id="kitsuPassword" type="password" autocomplete="current-password">
      </div>
    </div>
    <div class="wiz-actions" style="margin-top:16px">
      <button class="btn" id="testKitsuBtn" onclick="testKitsu()">%s</button>
      <span id="step1Spinner" style="display:none"><span class="spinner"></span></span>
    </div>
    <div id="step1Error" style="display:none;margin-top:10px"></div>
  </div>
</div>

<!-- Step 2: Discord -->
<div class="wiz-panel" id="panel2">
  <div class="section-card glass" style="max-width:600px">
    <h3>%s</h3>
    <div id="step2Status"></div>
    <div id="step2Form" class="wiz-form">
      <div>
        <label>Bot Token</label>
        <input id="discordToken" type="password" placeholder="MTUwMTUx..." autocomplete="off">
      </div>
      <div>
        <label>Guild ID</label>
        <input id="discordGuildID" type="text" placeholder="1234567890123456789" autocomplete="off">
      </div>
    </div>
    <div class="wiz-actions" style="margin-top:16px">
      <button class="btn" id="testDiscordBtn" onclick="testDiscord()">%s</button>
      <span id="step2Spinner" style="display:none"><span class="spinner"></span></span>
    </div>
    <div id="step2Error" style="display:none;margin-top:10px"></div>
  </div>
</div>

<!-- Step 3: Project -->
<div class="wiz-panel" id="panel3">
  <div class="section-card glass" style="max-width:600px">
    <h3>%s</h3>
    <p class="hint">%s</p>
    <div id="step3ProjectList"><div class="badge-info status-badge">%s<span class="spinner"></span></div></div>
    <div id="step3Form" style="display:none" class="wiz-form">
      <div>
        <label>%s</label>
        <select id="projectSelect"></select>
      </div>
      <div>
        <label>%s</label>
        <select id="langSelect">
          <option value="ja">日本語</option>
          <option value="en">English</option>
        </select>
      </div>
    </div>
    <div class="wiz-actions" style="margin-top:16px">
      <button class="btn" id="applyBtn" onclick="applyProject()" style="display:none">%s</button>
      <span id="step3Spinner" style="display:none"><span class="spinner"></span></span>
    </div>
    <div id="step3Result" style="display:none;margin-top:16px"></div>
  </div>
</div>

<!-- Step 4: User & Checker Mapping (Optional) -->
<div class="wiz-panel" id="panel4">
  <div class="section-card glass" style="max-width:680px">
    <div class="optional-notice">
      <span class="optional-icon">📌</span>
      <span>%s</span>
    </div>
    <h3>%s</h3>
    <div id="step4Loading"><div class="badge-info status-badge">%s<span class="spinner"></span></div></div>
    <div id="step4Content" style="display:none">
      <div class="map-section">
        <div class="map-section-title">%s</div>
        <div id="userMappingList"></div>
      </div>
      <div class="map-section">
        <div class="map-section-title">%s</div>
        <div id="checkerMappingList"></div>
      </div>
    </div>
    <div class="wiz-actions" style="margin-top:20px">
      <button class="btn" id="saveMappingBtn" onclick="saveMapping()" style="display:none">%s</button>
      <button class="btn-ghost" id="skipMappingBtn" onclick="skipMapping()" style="display:none">%s</button>
      <span id="step4Spinner" style="display:none"><span class="spinner"></span></span>
    </div>
    <div id="step4Error" style="display:none;margin-top:10px"></div>
  </div>
</div>

<!-- Done -->
<div id="wizardDone" style="display:none;margin-top:20px">
  <div class="section-card glass" style="text-align:center;padding:32px 20px">
    <div style="font-size:2.4rem;margin-bottom:8px">✅</div>
    <h2 style="color:var(--success);margin:0 0 10px">%s</h2>
    <p class="hint">%s</p>
    <div style="display:flex;gap:12px;justify-content:center;flex-wrap:wrap;margin-top:20px">
      <a class="btn" href="%s">%s</a>
      <a class="btn-ghost" href="%s">%s</a>
    </div>
  </div>
</div>

<script>
var API_BASE = %q;
var stepDone = {1: false, 2: false, 3: false, 4: false};
var currentProjectID = '';
var mappingData = null;

function el(id){ return document.getElementById(id); }

function showTab(n){
  for(var i=1;i<=4;i++){
    var tab = el('tab'+i), panel = el('panel'+i);
    if(!tab||!panel) continue;
    panel.classList.toggle('active', i===n);
    tab.classList.toggle('active', i===n && !stepDone[i]);
  }
}

function markTabDone(n){
  stepDone[n] = true;
  var tab = el('tab'+n);
  if(tab){ tab.classList.remove('active'); tab.classList.add('done'); }
}

function badge(cls, msg){ return '<div class="status-badge '+cls+'">'+msg+'</div>'; }

function setBtn(id, loading){
  var b = el(id);
  if(!b) return;
  b.disabled = loading;
}

function showSpinner(id, show){
  var s = el(id);
  if(s) s.style.display = show ? '' : 'none';
}

function resultLines(lines){
  if(!lines||!lines.length) return '';
  var html = '<ul>';
  lines.forEach(function(l){
    var cls = l.startsWith('OK:') ? 'ok' : l.startsWith('FAIL:') ? 'fail' : l.startsWith('WARN:') ? 'warn' : '';
    html += '<li class="'+cls+'">'+escHtml(l)+'</li>';
  });
  return html+'</ul>';
}

function escHtml(s){
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

async function loadStatus(){
  try{
    var res = await fetch(API_BASE+'/status');
    var d = await res.json();
    applyStatus(d);
  }catch(e){
    el('step1Status').innerHTML = badge('badge-warn', '%s');
  }
}

function applyStatus(d){
  var kitsuOK = d.kitsu && d.kitsu.authenticated;
  var discordOK = d.discord && d.discord.bot_valid && d.discord.guild_valid &&
    d.discord.permissions && d.discord.permissions.manage_channels && d.discord.permissions.manage_webhooks;
  var projectOK = d.project && d.project.selected;

  if(kitsuOK){
    el('step1Status').innerHTML = badge('badge-ok', '✓ Kitsu: %s');
    el('step1Form').style.display = 'none';
    markTabDone(1);
  } else if(d.kitsu && d.kitsu.configured && d.kitsu.error){
    el('step1Status').innerHTML = badge('badge-err', '✗ '+escHtml(d.kitsu.error));
  }

  if(discordOK){
    var name = (d.discord.bot_name || 'Bot')+' / '+(d.discord.guild_name || 'Guild');
    el('step2Status').innerHTML = badge('badge-ok', '✓ Discord: '+escHtml(name));
    el('step2Form').style.display = 'none';
    markTabDone(2);
  } else if(d.discord && d.discord.configured && d.discord.error){
    el('step2Status').innerHTML = badge('badge-err', '✗ '+escHtml(d.discord.error));
  }

  if(projectOK){
    currentProjectID = d.project.project_id || '';
    markTabDone(3);
  }

  if(!kitsuOK){ showTab(1); }
  else if(!discordOK){ showTab(2); }
  else if(!projectOK){ showTab(3); loadProjects(); }
  else { showTab(4); loadMapping(); }
}

// --- Step 1 ---

async function testKitsu(){
  var hostname = el('kitsuHostname').value.trim();
  var email = el('kitsuEmail').value.trim();
  var password = el('kitsuPassword').value;
  if(!hostname||!email||!password){
    el('step1Error').innerHTML = badge('badge-err', '%s');
    el('step1Error').style.display='';
    return;
  }
  showSpinner('step1Spinner', true);
  setBtn('testKitsuBtn', true);
  el('step1Error').style.display='none';
  try{
    var res = await fetch(API_BASE+'/test-kitsu', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({hostname, email, password})
    });
    var d = await res.json();
    if(d.authenticated){
      el('step1Status').innerHTML = badge('badge-ok', '✓ Kitsu: %s');
      el('step1Form').style.display='none';
      markTabDone(1);
      showTab(2);
    } else {
      el('step1Error').innerHTML = badge('badge-err', '✗ '+(d.error||'%s'));
      el('step1Error').style.display='';
    }
  }catch(e){
    el('step1Error').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
    el('step1Error').style.display='';
  }
  showSpinner('step1Spinner', false);
  setBtn('testKitsuBtn', false);
}

// --- Step 2 ---

async function testDiscord(){
  var botToken = el('discordToken').value.trim();
  var guildID = el('discordGuildID').value.trim();
  if(!botToken||!guildID){
    el('step2Error').innerHTML = badge('badge-err', '%s');
    el('step2Error').style.display='';
    return;
  }
  showSpinner('step2Spinner', true);
  setBtn('testDiscordBtn', true);
  el('step2Error').style.display='none';
  try{
    var res = await fetch(API_BASE+'/test-discord', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({bot_token: botToken, guild_id: guildID})
    });
    var d = await res.json();
    if(d.bot_valid && d.guild_valid && d.permissions && d.permissions.manage_channels && d.permissions.manage_webhooks){
      var name = (d.bot_name||'Bot')+' / '+(d.guild_name||guildID);
      el('step2Status').innerHTML = badge('badge-ok', '✓ Discord: '+escHtml(name));
      el('step2Form').style.display='none';
      markTabDone(2);
      showTab(3);
      loadProjects();
    } else {
      var errMsg = d.error || '';
      if(d.bot_valid && d.guild_valid && d.permissions){
        var missing = [];
        if(!d.permissions.manage_channels) missing.push('MANAGE_CHANNELS');
        if(!d.permissions.manage_webhooks) missing.push('MANAGE_WEBHOOKS');
        if(missing.length) errMsg = '%s: '+missing.join(', ');
      }
      el('step2Error').innerHTML = badge('badge-err', '✗ '+(errMsg||'%s'));
      el('step2Error').style.display='';
    }
  }catch(e){
    el('step2Error').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
    el('step2Error').style.display='';
  }
  showSpinner('step2Spinner', false);
  setBtn('testDiscordBtn', false);
}

// --- Step 3 ---

async function loadProjects(){
  el('step3ProjectList').innerHTML = badge('badge-info', '%s<span class="spinner"></span>');
  try{
    var res = await fetch(API_BASE+'/projects');
    var projects = await res.json();
    renderProjects(projects);
  }catch(e){
    el('step3ProjectList').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
  }
}

function renderProjects(projects){
  if(!projects||!projects.length){
    el('step3ProjectList').innerHTML = badge('badge-warn', '%s');
    return;
  }

  var configured = projects.filter(function(p){ return p.is_configured; });
  var available = projects.filter(function(p){ return !p.is_configured; });

  var html = '';

  if(configured.length){
    html += '<div style="margin-bottom:12px"><div style="font-size:.76rem;color:var(--muted-2);margin-bottom:6px">%s</div>';
    configured.forEach(function(p){
      html += '<div class="project-item"><span>'+escHtml(p.project_name)+'</span><span class="project-badge">%s</span></div>';
    });
    html += '</div>';
  }

  if(!available.length){
    html += badge('badge-ok', '✓ %s');
    el('step3ProjectList').innerHTML = html;
    markTabDone(3);
    showTab(4);
    loadMapping();
    return;
  }

  html += '<div style="font-size:.76rem;color:var(--muted-2);margin-bottom:6px">%s</div>';
  el('step3ProjectList').innerHTML = html;

  var sel = el('projectSelect');
  sel.innerHTML = '<option value="">-- %s --</option>';
  available.forEach(function(p){
    var opt = document.createElement('option');
    opt.value = p.project_id;
    opt.textContent = p.project_name+' ('+p.task_type_count+' task types)';
    sel.appendChild(opt);
  });

  el('step3Form').style.display='';
  el('applyBtn').style.display='';
}

async function applyProject(){
  var projectID = el('projectSelect').value;
  var lang = el('langSelect').value;
  if(!projectID){
    el('step3Result').innerHTML = badge('badge-err', '%s');
    el('step3Result').style.display='';
    return;
  }
  el('step3Result').style.display='none';
  showSpinner('step3Spinner', true);
  setBtn('applyBtn', true);
  try{
    var res = await fetch(API_BASE+'/apply-project', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({project_id: projectID, template: 'cg', language: lang})
    });
    var d = await res.json();
    if(d.ok){
      currentProjectID = projectID;
      el('step3Result').innerHTML =
        badge('badge-ok', '✓ '+escHtml(d.project_name)+': '+d.channels_created+' %s, '+d.webhooks_created+' webhooks') +
        resultLines(d.lines);
      el('step3Result').style.display='';
      markTabDone(3);
      setTimeout(function(){ showTab(4); loadMapping(); }, 600);
    } else {
      el('step3Result').innerHTML =
        badge('badge-err', '✗ '+(d.error||'%s')) +
        resultLines(d.lines);
      el('step3Result').style.display='';
      if(d.safe_to_retry){
        el('step3Result').innerHTML += badge('badge-warn', '%s');
      }
    }
  }catch(e){
    el('step3Result').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
    el('step3Result').style.display='';
  }
  showSpinner('step3Spinner', false);
  setBtn('applyBtn', false);
}

// --- Step 4 ---

async function loadMapping(){
  el('step4Loading').style.display='';
  el('step4Content').style.display='none';
  el('saveMappingBtn').style.display='none';
  el('skipMappingBtn').style.display='none';
  try{
    var res = await fetch(API_BASE+'/mapping');
    var d = await res.json();
    if(d.error){
      el('step4Loading').innerHTML = badge('badge-warn', d.error);
      el('skipMappingBtn').style.display='';
      return;
    }
    mappingData = d;
    renderMapping(d);
  }catch(e){
    el('step4Loading').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
    el('skipMappingBtn').style.display='';
  }
}

function renderMapping(d){
  // Build existing lookup maps
  var userLookup = {};
  (d.user_maps||[]).forEach(function(m){ userLookup[m.kitsu_name] = m.discord_user_id; });
  var checkerLookup = {};
  (d.checker_maps||[]).forEach(function(m){ checkerLookup[m.task_type] = m.discord_user_id; });

  // User mapping rows
  var userHTML = '';
  if(!d.persons||!d.persons.length){
    userHTML = '<p class="hint">%s</p>';
  } else {
    d.persons.forEach(function(p){
      var existing = userLookup[p.kitsu_name] || '';
      userHTML += '<div class="map-row">' +
        '<div class="map-person"><div class="map-name">'+escHtml(p.kitsu_name)+'</div>' +
        '<div class="map-email">'+escHtml(p.kitsu_email)+'</div></div>' +
        '<input class="map-input user-discord-input" data-name="'+escHtml(p.kitsu_name)+'" data-email="'+escHtml(p.kitsu_email)+'" ' +
        'type="text" placeholder="Discord ID" value="'+escHtml(existing)+'">' +
        '</div>';
    });
  }
  el('userMappingList').innerHTML = userHTML;

  // Checker mapping rows
  var checkerHTML = '';
  if(!d.task_types||!d.task_types.length){
    checkerHTML = '<p class="hint">%s</p>';
  } else {
    d.task_types.forEach(function(tt){
      var existing = checkerLookup[tt] || '';
      checkerHTML += '<div class="map-row">' +
        '<div class="map-task">'+escHtml(tt)+'</div>' +
        '<input class="map-input checker-discord-input" data-tasktype="'+escHtml(tt)+'" ' +
        'type="text" placeholder="Discord ID" value="'+escHtml(existing)+'">' +
        '</div>';
    });
  }
  el('checkerMappingList').innerHTML = checkerHTML;

  el('step4Loading').style.display='none';
  el('step4Content').style.display='';
  el('saveMappingBtn').style.display='';
  el('skipMappingBtn').style.display='';
}

async function saveMapping(){
  if(!mappingData){ showDone(); return; }
  showSpinner('step4Spinner', true);
  setBtn('saveMappingBtn', true);
  el('step4Error').style.display='none';

  var userMappings = [];
  document.querySelectorAll('.user-discord-input').forEach(function(inp){
    userMappings.push({
      kitsu_name: inp.dataset.name,
      kitsu_email: inp.dataset.email,
      discord_user_id: inp.value.trim()
    });
  });

  var checkerMappings = [];
  document.querySelectorAll('.checker-discord-input').forEach(function(inp){
    checkerMappings.push({
      task_type: inp.dataset.tasktype,
      discord_user_id: inp.value.trim()
    });
  });

  var projectID = (mappingData && mappingData.project_id) || currentProjectID;
  var ok = true;

  try{
    if(userMappings.length){
      var r1 = await fetch(API_BASE+'/mapping/users', {
        method:'POST', headers:{'Content-Type':'application/json'},
        body: JSON.stringify({project_id: projectID, mappings: userMappings})
      });
      if(!r1.ok) ok = false;
    }
    if(checkerMappings.length){
      var r2 = await fetch(API_BASE+'/mapping/checkers', {
        method:'POST', headers:{'Content-Type':'application/json'},
        body: JSON.stringify({project_id: projectID, mappings: checkerMappings})
      });
      if(!r2.ok) ok = false;
    }
  }catch(e){
    el('step4Error').innerHTML = badge('badge-err', '✗ '+escHtml(e.message));
    el('step4Error').style.display='';
    showSpinner('step4Spinner', false);
    setBtn('saveMappingBtn', false);
    return;
  }

  showSpinner('step4Spinner', false);
  setBtn('saveMappingBtn', false);
  if(ok){ markTabDone(4); showDone(); }
  else {
    el('step4Error').innerHTML = badge('badge-err', '%s');
    el('step4Error').style.display='';
  }
}

function skipMapping(){
  markTabDone(4);
  showDone();
}

function showDone(){
  el('wizardDone').style.display='';
  el('panel4').classList.remove('active');
  window.scrollTo({top:0, behavior:'smooth'});
}

document.addEventListener('DOMContentLoaded', loadStatus);
</script>`,
		// step tab labels
		t(lang, "プロジェクト", "Project"),
		t(lang, "マッピング", "Mapping"),
		// step 1 heading
		t(lang, "Step 1: Kitsu 接続確認", "Step 1: Verify Kitsu Connection"),
		// field labels
		t(lang, "メールアドレス", "Email"),
		t(lang, "パスワード", "Password"),
		// button
		t(lang, "接続テスト", "Test Connection"),
		// step 2 heading
		t(lang, "Step 2: Discord Bot 確認", "Step 2: Verify Discord Bot"),
		// button
		t(lang, "Bot を確認", "Verify Bot"),
		// step 3 heading
		t(lang, "Step 3: プロジェクトセットアップ", "Step 3: Project Setup"),
		// step 3 hint
		t(lang, "プロジェクトを選択して Discord チャンネルを作成します。", "Select a project to create Discord channels."),
		// loading label
		t(lang, "プロジェクトを読み込み中...", "Loading projects..."),
		// field labels
		t(lang, "プロジェクト", "Project"),
		t(lang, "通知言語", "Notification language"),
		// apply button
		t(lang, "チャンネルを作成", "Create Channels"),
		// step 4 optional notice
		t(lang, "この設定は Optional です。未設定でも通知は届きます（@mention がつきません）。プロジェクトメンバーの Discord ID を入力するとメンションが届きます。", "This step is optional. Notifications will still be sent without @mentions. Enter Discord IDs to enable @mention for each member."),
		// step 4 heading
		t(lang, "Step 4: ユーザー・チェッカーマッピング", "Step 4: User & Checker Mapping"),
		// loading
		t(lang, "マッピング情報を読み込み中...", "Loading mapping data..."),
		// section titles
		t(lang, "Kitsu ユーザー → Discord", "Kitsu Users → Discord"),
		t(lang, "タスクタイプ別チェッカー", "Task Type Checkers"),
		// buttons
		t(lang, "保存して完了", "Save & Finish"),
		t(lang, "スキップ", "Skip"),
		// done section
		t(lang, "セットアップ完了", "Setup Complete"),
		t(lang, "Discord チャンネルの作成が完了しました。通知の到着を確認してください。", "Discord channels have been created. Verify that notifications arrive."),
		adminURL,
		t(lang, "管理画面へ", "Go to Admin"),
		setupURL,
		t(lang, "Setup に戻る", "Back to Setup"),
		// JS string: API_BASE
		apiBase,
		// inline JS strings
		t(lang, "ステータス取得に失敗しました", "Failed to load status"),
		t(lang, "認証済み", "Authenticated"),
		t(lang, "すべての項目を入力してください", "All fields are required"),
		t(lang, "認証成功", "Authenticated"),
		t(lang, "認証失敗", "Authentication failed"),
		t(lang, "すべての項目を入力してください", "All fields are required"),
		t(lang, "権限不足", "Missing permissions"),
		t(lang, "確認失敗", "Verification failed"),
		t(lang, "プロジェクトを読み込み中...", "Loading projects..."),
		t(lang, "Kitsuにプロジェクトが見つかりません", "No projects found in Kitsu"),
		t(lang, "設定済みプロジェクト", "Configured projects"),
		t(lang, "設定済み", "configured"),
		t(lang, "全プロジェクトの設定が完了しています", "All projects are configured"),
		t(lang, "未設定のプロジェクト", "Available projects"),
		t(lang, "プロジェクトを選択", "Select a project"),
		t(lang, "プロジェクトを選択してください", "Please select a project"),
		t(lang, "チャンネル", "channels"),
		t(lang, "セットアップ失敗", "Setup failed"),
		t(lang, "✅ 再試行可能です", "✅ Safe to retry"),
		t(lang, "Kitsu ユーザーが見つかりません", "No Kitsu users found"),
		t(lang, "チャンネルが作成されていません", "No channels created yet"),
		t(lang, "保存に失敗しました", "Failed to save"),
	)

	return adminPage(lang, t(lang, "セットアップウィザード", "Setup Wizard"), r, body)
}
