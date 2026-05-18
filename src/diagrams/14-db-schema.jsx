const DBSchema = () => {
    return (
        <section className="doc" id="db-schema">
            <div className="kicker">Database</div>
            <h1 className="docnum"><span className="n">14</span>Phase 2-2 Database Schema</h1>
            <h2 className="docsub">追加および更新されたテーブル構成</h2>

            <div style={{display: 'flex', gap: '24px', flexWrap: 'wrap'}}>
                {/* ProjectDCCTool Table */}
                <div className="er-table">
                    <div className="h">
                        <span>ProjectDCCTool</span>
                        <span>Phase 2-2</span>
                    </div>
                    <div className="row pk">
                        <span>id</span>
                        <span>uint</span>
                    </div>
                    <div className="row fk">
                        <span>kitsu_project_id</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>dcc_type</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>file_path_pattern</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>entity_type_filter</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>enabled</span>
                        <span>bool</span>
                    </div>
                </div>



                {/* AuditLog Table */}
                <div className="er-table">
                    <div className="h">
                        <span>AuditLog</span>
                        <span>Phase 2-2</span>
                    </div>
                    <div className="row pk">
                        <span>id</span>
                        <span>uint</span>
                    </div>
                    <div className="row fk">
                        <span>kitsu_project_id</span>
                        <span>string</span>
                    </div>
                    <div className="row fk">
                        <span>task_id</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>success</span>
                        <span>bool</span>
                    </div>
                    <div className="row">
                        <span>retry_count</span>
                        <span>int</span>
                    </div>
                    <div className="row">
                        <span>error_message</span>
                        <span>string</span>
                    </div>
                </div>

                {/* MentionRule Table */}
                <div className="er-table">
                    <div className="h">
                        <span>MentionRule</span>
                        <span>Phase 2-2</span>
                    </div>
                    <div className="row pk">
                        <span>id</span>
                        <span>uint</span>
                    </div>
                    <div className="row fk">
                        <span>kitsu_project_id</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>task_type</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>task_status</span>
                        <span>string</span>
                    </div>
                    <div className="row">
                        <span>discord_role_id</span>
                        <span>string</span>
                    </div>
                </div>
            </div>
            
            <div className="wf wf-pad box" style={{marginTop: '32px'}}>
                <div className="box-title">Schema Notes</div>
                <ul style={{margin: 0, paddingLeft: '20px'}}>
                    <li><code>ProjectDCCTool</code> は <code>FilePathPreset</code> を利用してテンプレートから素早く作成可能です。</li>
                    <li><code>AuditLog</code> は容量節約のため、作成から30日経過したレコードを自動パージします。</li>
                    <li><code>MentionRule</code> は <code>task_type</code> または <code>task_status</code> が空欄の場合、「すべてに適用（ワイルドカード）」として機能します。</li>
                    <li><code>UserMap</code> にはフェーズ2-2の一環で、リネーム保護のための <code>kitsu_email</code> フィールドが追加されました。</li>
                </ul>
            </div>
        </section>
    );
};
