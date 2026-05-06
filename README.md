# Kitsu × Discord リッチ通知システム

> Kitsu のタスクステータス変更を Discord に自動通知する Bot。  
> VFX・アニメ制作チーム向けに日本語対応・メンション・プレビュー画像表示に対応。

[English](#english) | [日本語](#日本語)

---

## 日本語

### 概要

[Kitsu](https://www.cg-wire.com/) のポーリング（デフォルト1分間隔）でタスクの状態変化を検知し、Discord Webhook でリッチ通知を送信します。

**主な機能:**
- 🔔 ステータス変更・コメント追加を自動検知
- 👤 担当者・チェッカーへの @メンション（conf.toml で設定）
- 📸 Kitsu プレビュー画像を Discord に表示（nginx 設定が必要、後述）
- 🧵 タスクごとのスレッド作成（`useThreads = true`）
- 📊 工程別・プロジェクト別チャンネルへの振り分け
- 🚨 緊急ステータスで `@here` 通知
- 📅 毎朝9時の日次サマリー
- 🏥 ヘルスチェックエンドポイント（`:8090/health`）
- 🔁 Kitsu API 障害時の自動リトライ（指数バックオフ）

### 必要環境

- Docker + Docker Compose
- Kitsu サーバー（self-hosted または Zou/CGWire cloud）
- Discord サーバーの Webhook URL

### セットアップ

```bash
# 1. リポジトリをクローン
git clone https://github.com/YOUR_USERNAME/kitsu-discord-rich-notification.git
cd kitsu-discord-rich-notification

# 2. 設定ファイルを作成
cp conf.toml.example conf.toml
cp .env.example .env

# 3. .env を編集（パスワード・Webhook URL）
nano .env

# 4. conf.toml を編集（Kitsu ホスト・チームメンバー設定）
nano conf.toml

# 5. 起動
docker-compose up -d

# 6. ログ確認
docker-compose logs -f app
```

### conf.toml 主要設定

| キー | 説明 |
|------|------|
| `kitsu.hostname` | Kitsu サーバーの URL（末尾スラッシュ必須） |
| `discord.webhookURL` | 通知先 Discord Webhook URL |
| `discord.useThreads` | `true` でタスクごとにスレッドを作成 |
| `mention.checkerStatuses` | チェッカーをメンションするステータス（例: `["WFA"]`） |
| `mention.artistStatuses` | アーティストをメンションするステータス |
| `mention.hereStatuses` | `@here` を送るステータス（緊急用、空でも可） |
| `[[mention.userMap]]` | Kitsu フルネーム → Discord ユーザーID の紐付け |
| `[[mention.checkers]]` | タスクタイプ → チェッカー Discord ID |
| `[[discord.taskTypeWebhooks]]` | 工程別チャンネル振り分け |
| `[[discord.productions]]` | プロジェクト別チャンネル振り分け |

詳細は `conf.toml.example` のコメントを参照してください。

### プレビュー画像表示（オプション）

Discord の embed 画像には **認証なしでアクセス可能な URL** が必要です。Kitsu のプレビュー API はデフォルトで JWT 認証が必要なため、nginx にバイパス設定を追加します。

**nginx 設定（Kitsu サーバー上で実施）:**

```nginx
# サムネイルのみ認証なしで公開（ファイル名が UUID なので推測困難）
location ~ ^/api/pictures/thumbnails/preview-files/ {
    proxy_pass http://localhost:5000;   # Kitsu バックエンド（zou）のポートに合わせる
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

```bash
sudo nginx -t && sudo systemctl reload nginx
```

設定後、以下で動作確認できます：
```bash
# 認証なしで 200 が返れば OK（{id} は任意の preview_file_id に置換）
curl -o /dev/null -w "%{http_code}" http://YOUR_KITSU_HOST/api/pictures/thumbnails/preview-files/{id}.png
```

> ⚠️ サムネイル URL は UUID で推測困難ですが、Discord チャンネル自体のアクセス制御を必ず設定してください。社内素材が外部に漏洩しないよう注意。

### ディレクトリ構成

```
├── conf.toml.example      # 設定テンプレート
├── .env.example           # 環境変数テンプレート
├── docker-compose.yml
├── Dockerfile
├── tpl/
│   └── rich/              # Discord メッセージテンプレート（Go text/template）
│       ├── author.tpl
│       ├── title.tpl
│       ├── description.tpl
│       ├── fields.tpl
│       └── footer.tpl
└── src/
    ├── main.go
    ├── api/
    │   ├── kitsu/         # Kitsu API クライアント
    │   └── discord/       # Discord Webhook 送信
    ├── model/             # SQLite モデル（GORM）
    └── utils/
        ├── config/        # TOML 設定・${VAR} 環境変数展開
        ├── request/       # HTTP クライアント（リトライ+バックオフ）
        └── truncate/
```

### 運用コマンド

```bash
# ログ確認
docker-compose logs -f app

# 設定変更後の再起動（conf.toml, .env を変えたとき）
docker-compose restart app

# コードを含む完全再ビルド
docker-compose down && docker-compose up -d --build

# ヘルスチェック確認
curl http://localhost:8090/health
# → {"status":"ok"}

# コンテナの状態確認
docker-compose ps
```

### トラブルシューティング

| 症状 | 確認箇所 |
|------|----------|
| 通知が来ない | `docker-compose logs app` でエラーを確認 |
| メンションされない | `[[mention.userMap]]` に全員分の Discord ID があるか確認 |
| 起動時に `[WARN] placeholder` が出る | `conf.toml` のプレースホルダを実際の値に書き換える |
| 画像が表示されない | nginx の preview-files location を設定済みか確認 |
| `Initial Kitsu authentication failed` | hostname / email / password を確認 |

---

## English

### Overview

Polls [Kitsu](https://www.cg-wire.com/) every minute, detects task status changes, and posts rich Discord notifications. Built for Japanese VFX/animation studios but works with any Kitsu setup.

**Features:**
- 🔔 Auto-detect status changes and new comments
- 👤 @mention assignees and checkers, configurable per status
- 📸 Show Kitsu preview thumbnails in Discord (requires nginx setup)
- 🧵 Per-task Discord threads (`useThreads = true`)
- 📊 Route notifications by project or task type to separate channels
- 🚨 `@here` mention for urgent statuses
- 📅 Daily digest at 9 AM (JST)
- 🏥 Health check at `:8090/health`
- 🔁 Automatic retry with exponential backoff on Kitsu API failures

### Quick Start

```bash
git clone https://github.com/YOUR_USERNAME/kitsu-discord-rich-notification.git
cd kitsu-discord-rich-notification

cp conf.toml.example conf.toml
cp .env.example .env
# Fill in .env (Kitsu password + Discord webhook URL)
# Edit conf.toml (Kitsu hostname, team members, mention settings)

docker-compose up -d
docker-compose logs -f app
```

Secrets go in `.env` only — never in `conf.toml`. See `conf.toml.example` for all available options with inline comments.

### Architecture

```
Kitsu API ──(poll every N min)──► Go app ──► Discord Webhook
                                     │
                                   SQLite (change detection)
```

The app stores task state in SQLite. On each poll it compares the stored state to the current Kitsu state and only sends a Discord message when something changed (status, timestamp, or latest comment).

---

## Credits

Forked from [keshon/kitsu-to-discord-task-notification](https://github.com/keshon/kitsu-to-discord-task-notification) and heavily extended for Japanese VFX/animation production workflows.

## License

MIT — see [LICENSE](LICENSE).
