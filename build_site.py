#!/usr/bin/env python3
import os, re

from pathlib import Path

base = Path(__file__).resolve().parent
diagrams_dir = base / 'diagrams'
out_path = base / 'site.jsx'

files = sorted([f for f in os.listdir(diagrams_dir) if f.endswith('.jsx')])
parts = []
for f in files:
    try:
        content = (diagrams_dir / f).read_text(encoding='utf-8')
        content = re.sub(r'^window\.', 'const ', content, count=1)
        parts.append(content)
    except Exception as e:
        print(f"Error reading {f}: {e}")

shell = r"""
const { useState } = React;

const navSections = [
    { group: 'Core System', items: [
        { id: 'architecture',  label: 'System Architecture',  num: '01' },
        { id: 'er-diagram',    label: 'DB Schema (ER)',        num: '02' },
        { id: 'admin-ui',      label: 'Admin UI',              num: '03' },
        { id: 'status-map',    label: 'Status Map',            num: '04' },
        { id: 'routing',       label: 'Message Routing',       num: '05' },
    ]},
    { group: 'Discord Output', items: [
        { id: 'discord-messages', label: 'Message Format',    num: '06' },
        { id: 'mention-system',   label: 'Mention System',    num: '07' },
    ]},
    { group: 'Pipeline System', items: [
        { id: 'vfx-pipeline',  label: 'VFX Production Pipeline',num: '17' },
        { id: 'dataflow',      label: 'Data Flow & State',     num: '18' },
        { id: 'pipeline',      label: 'Full Pipeline Flow',    num: '08' },
        { id: 'templates',     label: 'Config Reference',      num: '09' },
        { id: 'status-all',    label: 'Template System',       num: '10' },
        { id: 'threads',       label: 'Thread & Lifecycle',    num: '11' },
    ]},
    { group: 'Phase Updates', items: [
        { id: 'phase-roadmap',      label: 'Phase Roadmap',        num: '19' },
        { id: 'phase-2-2-overview', label: 'Phase 2-2 Features',   num: '12' },
        { id: 'admin-routes',       label: 'Admin Routes',         num: '13' },
        { id: 'db-schema',          label: 'DB Schema Updates',    num: '14' },
        { id: 'webhook-flow',       label: 'Webhook Flow',         num: '15' },
        { id: 'dcc-integration',    label: 'DCC Integration',      num: '16' },
    ]},
];

const componentMap = {
    'architecture':        () => typeof Architecture        !== 'undefined' ? <Architecture />        : null,
    'er-diagram':          () => typeof ERDiagram           !== 'undefined' ? <ERDiagram />           : null,
    'admin-ui':            () => typeof AdminUI             !== 'undefined' ? <AdminUI />             : null,
    'status-map':          () => typeof StatusMap           !== 'undefined' ? <StatusMap />           : null,
    'routing':             () => typeof Routing             !== 'undefined' ? <Routing />             : null,
    'discord-messages':    () => typeof DiscordMessages     !== 'undefined' ? <DiscordMessages />     : null,
    'mention-system':      () => typeof MentionSystem       !== 'undefined' ? <MentionSystem />       : null,
    'pipeline':            () => typeof Pipeline            !== 'undefined' ? <Pipeline />            : null,
    'templates':           () => typeof Templates           !== 'undefined' ? <Templates />           : null,
    'status-all':          () => typeof StatusAll           !== 'undefined' ? <StatusAll />           : null,
    'threads':             () => typeof Threads             !== 'undefined' ? <Threads />             : null,
    'phase-2-2-overview':  () => typeof Phase22Overview     !== 'undefined' ? <Phase22Overview />     : null,
    'admin-routes':        () => typeof AdminRoutes         !== 'undefined' ? <AdminRoutes />         : null,
    'db-schema':           () => typeof DBSchema            !== 'undefined' ? <DBSchema />            : null,
    'webhook-flow':        () => typeof WebhookFlow         !== 'undefined' ? <WebhookFlow />         : null,
    'dcc-integration':     () => typeof DCCIntegration      !== 'undefined' ? <DCCIntegration />      : null,
    'vfx-pipeline':        () => typeof VFXPipeline         !== 'undefined' ? <VFXPipeline />         : null,
    'dataflow':            () => typeof DataFlow            !== 'undefined' ? <DataFlow />            : null,
    'phase-roadmap':       () => typeof PhaseRoadmap        !== 'undefined' ? <PhaseRoadmap />        : null,
};

function Site() {
    const [activeTab, setActiveTab] = useState('architecture');

    const content = componentMap[activeTab] ? componentMap[activeTab]() : <div style={{padding:'2rem'}}>Loading...</div>;

    return (
        <div className="site">
            <aside className="nav">
                <div className="brand">
                    Kitsu × Discord
                    <small>Pipeline Documentation</small>
                </div>

                {navSections.map(sec => (
                    <div className="navgroup" key={sec.group}>
                        <div className="navtitle">{sec.group}</div>
                        {sec.items.map(item => (
                            <a
                                key={item.id}
                                className={`navlink${activeTab === item.id ? ' active' : ''}`}
                                onClick={() => setActiveTab(item.id)}
                            >
                                <span className="num">{item.num}</span>
                                {item.label}
                            </a>
                        ))}
                    </div>
                ))}

                <div className="meta">
                    v2.2.0<br/>
                    Go + React + GORM
                </div>
            </aside>

            <main className="body">
                <div className="hero">
                    <div className="container">
                        <h1>Kitsu <em>×</em> Discord<br/>Pipeline</h1>
                        <p className="lead">
                            Kitsu のタスク更新を 1 分ポーリングで取得し、差分だけ Discord に送る運用パイプラインのドキュメントです。
                            現在動いている実装、運用上の制約、これから直すべき点を分けて整理しています。
                        </p>
                        <div className="stats">
                            <div className="stat"><div className="n">19</div><div className="l">Sections</div></div>
                            <div className="stat"><div className="n">7</div><div className="l">DB Tables</div></div>
                            <div className="stat"><div className="n">11</div><div className="l">Active Admin Routes</div></div>
                            <div className="stat"><div className="n">Now</div><div className="l">Operational Phase</div></div>
                        </div>
                        <div className="toc">
                            {navSections.map(sec =>
                                sec.items.map(item => (
                                    <a key={item.id} onClick={() => setActiveTab(item.id)}>
                                        {item.num} {item.label}
                                    </a>
                                ))
                            )}
                        </div>
                    </div>
                </div>

                <div className="container">
                    {content}
                </div>

                <footer className="foot">
                    <div className="left">Kitsu × Discord Pipeline v2.2</div>
                    <div>Go · React 18 · GORM · SQLite · Docker</div>
                </footer>
            </main>
        </div>
    );
}

const rootElement = document.getElementById('root');
if (rootElement) {
    ReactDOM.createRoot(rootElement).render(<Site />);
}
"""

full = '\n'.join(parts) + '\n' + shell
with open(out_path, 'w', encoding='utf-8') as f:
    f.write(full)
print(f'OK: site.jsx written, {len(full)} chars')
