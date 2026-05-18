window.AdminRoutes = () => {
    return (
        <section className="doc" id="admin-routes">
            <div className="kicker">Admin UI</div>
            <h1 className="docnum"><span className="n">13</span>Admin Routes Map</h1>
            <h2 className="docsub">????????????????????</h2>
            
            <div className="wf wf-pad box">
                <div className="box-title">Routing Tree</div>
                <code>
<pre style={{margin:0, fontFamily: 'var(--mono)', fontSize: '12px', lineHeight: '1.5'}}>
/bot
??? /login        [Kitsu ??]
??? /setup        [????????????]
??? /admin        [???????]
?   ??? /admin/users           ?? ?????????
?   ??? /admin/checkers        ? ???????
?   ??? /admin/drive           ?? Google Drive??
?   ??? /admin/bot             ?? Bot??
?   ??? /admin/dcc-tools       ?? DCC ?????
?   ??? /admin/audit           ?? ????
??? /logout       [?????]
</pre>
                </code>
            </div>

            <div className="variant-label">
                <span className="v">Routes</span>
                <h3>????????</h3>
            </div>
            <div className="wf box">
                <table className="adm">
                    <thead>
                        <tr>
                            <th>??</th>
                            <th>?????</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td><code>/admin/dcc-tools</code></td>
                            <td>Maya?Nuke ???????????????</td>
                        </tr>
                        <tr>
                            <td><code>/admin/audit</code></td>
                            <td>? Discord ?????????????????</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </section>
    );
};
