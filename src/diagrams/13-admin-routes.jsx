const AdminRoutes = () => {
    return (
        <section className="doc" id="admin-routes">
            <div className="kicker">Admin UI</div>
            <h1 className="docnum"><span className="n">13</span>Admin Routes Map</h1>
            <h2 className="docsub">管理画面のルーティング構成と新機能ページ</h2>
            
            <div className="wf wf-pad box">
                <div className="box-title">Routing Tree</div>
                <code>
<pre style={{margin:0, fontFamily: 'var(--mono)', fontSize: '12px', lineHeight: '1.5'}}>
/bot
├── /login        [Kitsu 認証]
├── /setup        [初期セットアップ]
├── /admin        [管理画面トップ]
│   ├── /admin/users           👥 ユーザーマップ
│   ├── /admin/checkers        ✅ チェッカー設定
│   ├── /admin/drive           📁 Google Drive設定
│   ├── /admin/bot             🤖 Bot設定
│   ├── /admin/project-channels 🔗 チャンネル管理
│   ├── <span style={{color: 'var(--accent)', fontWeight: 600}}>/admin/dcc-tools</span>       🎬 DCC ツール設定 (NEW)
│   ├── <span style={{color: 'var(--accent)', fontWeight: 600}}>/admin/path-mapping</span>    📂 パスマッピング (NEW)
│   ├── <span style={{color: 'var(--accent)', fontWeight: 600}}>/admin/mention-rules</span>   🔔 ロールメンション (NEW)
│   ├── <span style={{color: '#c8a13a', fontWeight: 600}}>/admin/webhooks</span>        🔗 Webhook 管理 (UPDATE)
│   └── <span style={{color: 'var(--accent)', fontWeight: 600}}>/admin/audit</span>           📋 監査ログ (NEW)
└── /logout       [ログアウト]
</pre>
                </code>
            </div>

            <div className="variant-label">
                <span className="v">NEW PAGES</span>
                <h3>追加・更新された管理画面</h3>
            </div>
            <div className="wf box">
                <table className="adm">
                    <thead>
                        <tr>
                            <th>パス</th>
                            <th>機能の役割</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td><code>/admin/dcc-tools</code></td>
                            <td>Maya・Nuke 等のファイルパスパターンと対応エンティティ設定</td>
                        </tr>
                        <tr>
                            <td><code>/admin/path-mapping</code></td>
                            <td>サーバーパスからローカルパス（OS別）への自動変換設定</td>
                        </tr>
                        <tr>
                            <td><code>/admin/mention-rules</code></td>
                            <td>タスク型×ステータスから Discord ロール ID へのマッピング設定</td>
                        </tr>
                        <tr>
                            <td><code>/admin/audit</code></td>
                            <td>全 Discord 通知の送受信履歴とエラー内容の表示</td>
                        </tr>
                        <tr>
                            <td><code>/admin/webhooks</code></td>
                            <td>複数 Webhook 管理から、1プロジェクト＝1WebhookURL の単一管理へ移行</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </section>
    );
};
