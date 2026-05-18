const WebhookFlow = () => {
    return (
        <section className="doc" id="webhook-flow">
            <div className="kicker">Notification Processing</div>
            <h1 className="docnum"><span className="n">15</span>Delivery Flow</h1>
            <h2 className="docsub">?????? webhook ???????polling ? diff ? delivery ???????</h2>

            <div className="wf wf-pad box">
                <div style={{display: 'flex', flexDirection: 'column', gap: '16px', maxWidth: '700px'}}>
                    {[
                        ['1', 'Kitsu Polling', '1???? Tasks / Comments / Persons / Projects / TaskTypes ???????'],
                        ['2', 'State Compare', 'sqlite.db ???????????????????????????????????'],
                        ['3', 'Route Resolve', 'project_webhooks ?????????? conf.toml ????????????????????'],
                        ['4', 'Discord Payload Build', 'Embed???????????ID?DCC?????????????'],
                        ['5', 'Send / Update', 'Webhook ??????? Message ID ????????????????'],
                        ['6', 'Audit & Persist', '????? task ? Discord message/thread ?????????'],
                    ].map(([step, title, body], idx) => (
                        <div key={idx} style={{display:'flex',gap:'12px'}}>
                            <div style={{width:'24px',height:'24px',borderRadius:'50%',background: idx < 3 ? 'var(--ink)' : 'var(--accent)', color:'var(--paper)',display:'flex',alignItems:'center',justifyContent:'center',fontFamily:'var(--mono)',fontSize:'12px'}}>{step}</div>
                            <div className="box fill" style={{flex:1,borderColor: idx < 3 ? 'var(--line)' : 'var(--accent)'}}>
                                <div className="box-title" style={{color: idx < 3 ? 'var(--ink-3)' : 'var(--accent)'}}>{title}</div>
                                <p style={{margin:0,fontSize:'12px'}}>{body}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="variant-label">
                <span className="v">Reality Check</span>
                <h3>?????????????</h3>
            </div>
            <div className="wf box">
                <table className="adm">
                    <thead><tr><th>??</th><th>??</th><th>??</th></tr></thead>
                    <tbody>
                        <tr><td>????</td><td>Polling</td><td><code>POST /bot/webhook</code> ??????????????????</td></tr>
                        <tr><td>????</td><td>????</td><td>Message ID / Thread ID ????????????????</td></tr>
                        <tr><td>429 / 5xx retry</td><td>?????????</td><td>docs ??????????????????????</td></tr>
                    </tbody>
                </table>
            </div>
        </section>
    );
};
