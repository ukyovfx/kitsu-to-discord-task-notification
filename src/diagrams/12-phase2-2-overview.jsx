const Phase22Overview = () => {
    return (
        <section className="doc" id="phase-2-2-overview">
            <div className="kicker">Phase 2-2 Features</div>
            <h1 className="docnum"><span className="n">12</span>DCC Integration & Monitoring</h1>
            <h2 className="docsub">Phase 2-2 で追加された4つの主要機能の概要</h2>
            <p className="doc-intro">
                DCC ツールの直接起動、エラー時のスマートリトライ、Discord通知の監査ログ、そしてステータスベースのロール別メンション機能を追加しました。
            </p>

            <div className="variant-label">
                <span className="v">P2-2-1</span>
                <h3>DCC Direct Link</h3>
                <span className="desc">vfx-launcher Integration</span>
            </div>
            <div className="wf wf-pad box">
                <h3 style={{fontSize: '14px', marginBottom: '12px'}}>✨ What's New</h3>
                    <li>URL Scheme: <code>vfx-launcher://open?software=maya&path=...</code></li>
                    <li>対応 DCC: Maya 🟢 / Nuke 🔴 / Blender 🟠 / Houdini 🟡</li>
                    <li>Discord Action Row ボタン（ファイル自動オープン）</li>
                    <li>Admin UI: /admin/dcc-tools</li>
                </ul>
                <div className="pill solid">ProjectDCCTool</div>
                <div className="pill solid" style={{marginLeft: '4px'}}>FilePathPreset</div>
            </div>

            <div className="variant-label">
                <span className="v">P2-2-2</span>
                <h3>Error Handling</h3>
                <span className="desc">Smart Retry Logic</span>
            </div>
            <div className="wf wf-pad box">
                <h3 style={{fontSize: '14px', marginBottom: '12px'}}>✨ What's New</h3>
                <ul style={{paddingLeft: '20px', marginBottom: '0'}}>
                    <li>429エラー → Retry-After header + 指数バックオフ</li>
                    <li>5xxエラー → 指数バックオフ: 1s → 2s → 4s（最大3回）</li>
                    <li>4xxエラー → リトライなし（即失敗）</li>
                    <li>自動リトライで一時的な通信エラーを自動回復し、送信成功率が向上</li>
                </ul>
            </div>

            <div className="variant-label">
                <span className="v">P2-2-3</span>
                <h3>Audit Trail</h3>
                <span className="desc">Message Audit Log</span>
            </div>
            <div className="wf wf-pad box">
                <h3 style={{fontSize: '14px', marginBottom: '12px'}}>✨ What's New</h3>
                <ul style={{paddingLeft: '20px', marginBottom: '16px'}}>
                    <li>全 Discord 送信を記録（task_id, success/failure, retry_count, error）</li>
                    <li>UI: /admin/audit（最新200件表示、成功/失敗で色分け）</li>
                    <li>30日で自動削除（肥大化防止）</li>
                </ul>
                <div className="pill solid">AuditLog</div>
            </div>

            <div className="variant-label">
                <span className="v">P2-2-4</span>
                <h3>Role Filtering</h3>
                <span className="desc">@Mention Rules</span>
            </div>
            <div className="wf wf-pad box">
                <h3 style={{fontSize: '14px', marginBottom: '12px'}}>✨ What's New</h3>
                <ul style={{paddingLeft: '20px', marginBottom: '16px'}}>
                    <li>TaskType × Status → Discord Role ID マッピング</li>
                    <li>例: "Compositing" + "WFA" → <code>@Reviewers</code></li>
                    <li>空欄 = すべてに適用（ワイルドカード機能）</li>
                </ul>
                <div className="pill solid">MentionRule</div>
            </div>
        </section>
    );
};
