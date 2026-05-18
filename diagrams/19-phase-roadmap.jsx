window.PhaseRoadmap = () => (
    <section className="doc" id="phase-roadmap">
        <div className="kicker">Project Evolution</div>
        <h1 className="docnum"><span className="n">19</span>Phase Roadmap</h1>
        <h2 className="docsub">?????????????????????????</h2>
        <p className="doc-intro">
            ??????????????????????????????????????????????????
        </p>

        <div className="variant-label">
            <span className="v">Phase 1</span>
            <h3>Phase 1: Basic Notification System (??)</h3>
        </div>
        <div className="wf wf-pad box">
            <ul style={{margin:0,paddingLeft:'20px',fontSize:'12px',lineHeight:'1.9',color:'var(--ink-2)'}}>
                <li>Kitsu ?????????????????</li>
                <li>???????????? Discord ??</li>
                <li>Task ??? message/thread ??</li>
            </ul>
        </div>

        <div className="variant-label">
            <span className="v">Phase 2-1</span>
            <h3>Phase 2-1: Setup & Admin Base (??)</h3>
        </div>
        <div className="wf wf-pad box">
            <ul style={{margin:0,paddingLeft:'20px',fontSize:'12px',lineHeight:'1.9',color:'var(--ink-2)'}}>
                <li><code>/setup</code> ????????????</li>
                <li>SQLite + GORM ???????????</li>
                <li>Kitsu JWT ????????????</li>
            </ul>
        </div>

        <div className="variant-label">
            <span className="v">Phase 2-2</span>
            <h3>Phase 2-2: Reliability & Pipeline UX (??? / ???)</h3>
        </div>
        <div className="wf wf-pad box">
            <ul style={{margin:0,paddingLeft:'20px',fontSize:'12px',lineHeight:'1.9',color:'var(--ink-2)'}}>
                <li><strong>Active:</strong> audit log?DCC tool ???project routing?docs ??</li>
                <li><strong>In Progress:</strong> secret ?????????????healthcheck ??</li>
                <li><strong>Needs Verification:</strong> retry ???mention rule UI?path mapping UI</li>
            </ul>
        </div>

        <div className="variant-label">
            <span className="v">Next</span>
            <h3>???????</h3>
        </div>
        <div className="wf wf-pad box">
            <ol style={{margin:0,paddingLeft:'20px',fontSize:'12px',lineHeight:'1.9',color:'var(--ink-2)'}}>
                <li><strong>Build Recovery:</strong> ???????? UTF-8 ????????????????</li>
                <li><strong>Secret Rotation:</strong> Discord Token / Webhook / Kitsu Password ????????</li>
                <li><strong>Config Unification:</strong> DB ??? live config ??????</li>
                <li><strong>Operational Hardening:</strong> HTTPS ?????session cookie ????????</li>
            </ol>
        </div>
    </section>
);
