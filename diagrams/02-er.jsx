window.ERDiagram = () => (
    <section className="doc" id="er-diagram">
        <div className="kicker">Database</div>
        <h1 className="docnum"><span className="n">02</span>DB Schema (ER Diagram)</h1>
        <h2 className="docsub">SQLite ????????????????????????</h2>
        <p className="doc-intro">
            ?????? GORM ??????????????????<code>sqlite.db</code> ????????????????????????????????????
            Secret ? DB ??????????????
        </p>

        <div style={{display:'flex',flexWrap:'wrap',gap:'20px',marginBottom:'32px'}}>
            <div className="er-table">
                <div className="h"><span>tasks</span><span>Core</span></div>
                <div className="row pk"><span>id</span><span>uint PK</span></div>
                <div className="row"><span>task_id</span><span>string idx</span></div>
                <div className="row"><span>task_updated_at</span><span>string</span></div>
                <div className="row"><span>task_status</span><span>string idx</span></div>
                <div className="row"><span>comment_id</span><span>string</span></div>
                <div className="row"><span>comment_updated_at</span><span>string</span></div>
                <div className="row"><span>discord_message_id</span><span>string</span></div>
                <div className="row"><span>webhook_url</span><span>legacy / ???</span></div>
                <div className="row"><span>discord_thread_id</span><span>string</span></div>
            </div>

            <div className="er-table">
                <div className="h"><span>projects</span><span>Setup</span></div>
                <div className="row pk"><span>id</span><span>uint PK</span></div>
                <div className="row"><span>kitsu_project_id</span><span>string unique</span></div>
                <div className="row"><span>name</span><span>string</span></div>
                <div className="row"><span>project_type</span><span>cg/jissha/anime</span></div>
                <div className="row"><span>discord_category_id</span><span>string</span></div>
                <div className="row"><span>language</span><span>ja / en</span></div>
                <div className="row"><span>storage_url</span><span>string</span></div>
            </div>

            <div className="er-table">
                <div className="h"><span>project_webhooks</span><span>Routing</span></div>
                <div className="row pk"><span>id</span><span>uint PK</span></div>
                <div className="row fk"><span>kitsu_project_id</span><span>idx ? projects</span></div>
                <div className="row"><span>channel_name</span><span>string</span></div>
                <div className="row"><span>task_type</span><span>"*" = catch-all</span></div>
                <div className="row"><span>webhook_url</span><span>runtime routing</span></div>
                <div className="row"><span>discord_channel_id</span><span>string</span></div>
            </div>

            <div className="er-table">
                <div className="h"><span>audit_logs</span><span>Observability</span></div>
                <div className="row pk"><span>id</span><span>uint PK</span></div>
                <div className="row"><span>created_at</span><span>idx</span></div>
                <div className="row"><span>task_id</span><span>idx</span></div>
                <div className="row"><span>project_name</span><span>string</span></div>
                <div className="row"><span>entity_name</span><span>string</span></div>
                <div className="row"><span>old_status / new_status</span><span>string</span></div>
                <div className="row"><span>discord_msg_id</span><span>string</span></div>
                <div className="row"><span>webhook_url</span><span>masked</span></div>
                <div className="row"><span>success / retry_count</span><span>bool / int</span></div>
            </div>
        </div>

        <div className="variant-label">
            <span className="v">State</span>
            <h3>??????????????</h3>
        </div>
        <div className="wf box">
            <table className="adm">
                <thead><tr><th>??</th><th>????</th><th>??</th></tr></thead>
                <tbody>
                    <tr><td>Kitsu hostname / email</td><td>????</td><td>?? UI ????????? seed ????</td></tr>
                    <tr><td>Kitsu password</td><td>?????</td><td>Secret ? DB ???????????????????</td></tr>
                    <tr><td>Discord bot token</td><td>?????</td><td>??????????????????????????</td></tr>
                    <tr><td>Main webhook URL</td><td>?????</td><td>??????? secret ??? DB ????????</td></tr>
                    <tr><td>Project webhook URL</td><td>????</td><td>????????????????? docs ???? URL ???????</td></tr>
                    <tr><td>Task ?????? URL</td><td>?????</td><td>?????? Message ID ?????????legacy ????????</td></tr>
                </tbody>
            </table>
        </div>
    </section>
);
