const WebhookFlow = () => {
    return (
        <section className="doc" id="webhook-flow">
            <div className="kicker">Webhook Processing</div>
            <h1 className="docnum"><span className="n">15</span>Webhook Flow — Phase 2-2 Enhanced</h1>
            <h2 className="docsub">KitsuのイベントからDiscordへの通知、そしてログ記録までの流れ</h2>

            <div className="wf wf-pad box">
                <div style={{display: 'flex', flexDirection: 'column', gap: '16px', maxWidth: '600px'}}>
                    
                    {/* Step 1 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--ink)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px'}}>1</div>
                        <div className="box fill" style={{flex: 1}}>
                            <div className="box-title">Kitsu Webhook 受信</div>
                            <p style={{margin: 0, fontSize: '12px'}}><code>POST /bot/webhook</code> — Task 更新通知受信</p>
                        </div>
                    </div>

                    <div style={{paddingLeft: '11px', margin: '-10px 0'}}>
                        <div style={{borderLeft: '2px dashed var(--line-soft)', height: '16px'}}></div>
                    </div>

                    {/* Step 2 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--ink)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px'}}>2</div>
                        <div className="box fill" style={{flex: 1}}>
                            <div className="box-title">Task データ抽出 + 認証確認</div>
                            <p style={{margin: 0, fontSize: '12px'}}>JWT Token 取得（セッションストア優先 → 環境変数フォールバック）</p>
                        </div>
                    </div>

                    <div style={{paddingLeft: '11px', margin: '-10px 0'}}>
                        <div style={{borderLeft: '2px dashed var(--line-soft)', height: '16px'}}></div>
                    </div>

                    {/* Step 3 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--ink)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px'}}>3</div>
                        <div className="box fill" style={{flex: 1}}>
                            <div className="box-title">Status フィルタリング</div>
                            <p style={{margin: 0, fontSize: '12px'}}>WFA / READY / DONE 等の重要ステータスのみ処理。それ以外はスキップ。</p>
                        </div>
                    </div>

                    <div style={{paddingLeft: '11px', margin: '-10px 0'}}>
                        <div style={{borderLeft: '2px dashed var(--line-soft)', height: '16px'}}></div>
                    </div>

                    {/* Step 4 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--accent)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px', boxShadow: '0 0 0 2px var(--paper), 0 0 0 4px var(--accent)'}}>4</div>
                        <div className="box fill" style={{flex: 1, borderColor: 'var(--accent)'}}>
                            <div className="box-title" style={{color: 'var(--accent)'}}>Discord メッセージ組立 (Phase 2-2 新機能)</div>
                            <ul style={{margin: 0, paddingLeft: '20px', fontSize: '12px'}}>
                                <li><strong>Base Embed:</strong> プロジェクト名・Entity・Task型・Status・色分け</li>
                                <li><strong>DCC Action Row:</strong> <code>[🔴 Open in Nuke] [🟢 Open in Maya]</code> などの vfx-launcher ボタン</li>
                                <li><strong>@Mention 組立:</strong> 既存のチェッカー・アーティストメンションに加え、<code>MentionRule</code> に基づくロールメンションを追加</li>
                            </ul>
                        </div>
                    </div>

                    <div style={{paddingLeft: '11px', margin: '-10px 0'}}>
                        <div style={{borderLeft: '2px dashed var(--line-soft)', height: '16px'}}></div>
                    </div>

                    {/* Step 5 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--accent)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px', boxShadow: '0 0 0 2px var(--paper), 0 0 0 4px var(--accent)'}}>5</div>
                        <div className="box fill" style={{flex: 1, borderColor: 'var(--accent)'}}>
                            <div className="box-title" style={{color: 'var(--accent)'}}>Webhook 送信 & エラーリトライ (Phase 2-2 新機能)</div>
                            <ul style={{margin: 0, paddingLeft: '20px', fontSize: '12px'}}>
                                <li><span style={{color: '#3e8f5b', fontWeight: 'bold'}}>2xx成功:</span> 即座に次へ</li>
                                <li><span style={{color: '#c8a13a', fontWeight: 'bold'}}>429エラー:</span> Retry-After 待機 + バックオフ</li>
                                <li><span style={{color: '#c44a2c', fontWeight: 'bold'}}>5xxエラー:</span> 指数バックオフ (1s→2s→4s, 最大3回)</li>
                            </ul>
                        </div>
                    </div>

                    <div style={{paddingLeft: '11px', margin: '-10px 0'}}>
                        <div style={{borderLeft: '2px dashed var(--line-soft)', height: '16px'}}></div>
                    </div>

                    {/* Step 6 */}
                    <div style={{display: 'flex', gap: '12px'}}>
                        <div style={{width: '24px', height: '24px', borderRadius: '50%', background: 'var(--accent)', color: 'var(--paper)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--mono)', fontSize: '12px', boxShadow: '0 0 0 2px var(--paper), 0 0 0 4px var(--accent)'}}>6</div>
                        <div className="box fill" style={{flex: 1, borderColor: 'var(--accent)'}}>
                            <div className="box-title" style={{color: 'var(--accent)'}}>監査ログ記録 (Phase 2-2 新機能)</div>
                            <p style={{margin: 0, fontSize: '12px'}}><code>AuditLog</code> に <code>task_id</code>, <code>success</code>, <code>retry_count</code>, <code>error_message</code> を保存し、管理画面から参照可能にする。</p>
                        </div>
                    </div>

                </div>
            </div>
        </section>
    );
};
