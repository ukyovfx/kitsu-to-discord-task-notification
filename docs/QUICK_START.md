# Quick Start

Get KitsuSync running in about 5 minutes.

---

## Prerequisites checklist

Before you begin, confirm you have:

- [ ] Docker and Docker Compose installed on the server
- [ ] Kitsu running and reachable from the server (test: `curl http://YOUR_KITSU_HOST/api/`)
- [ ] Discord server where you are an administrator
- [ ] Discord **Bot Token** — create one at [discord.com/developers/applications](https://discord.com/developers/applications) (Bot tab → Reset Token)
- [ ] One or more Discord **Guild IDs** — enable Developer Mode in Discord settings, then right-click each target server → Copy Server ID
- [ ] A dedicated Kitsu **runtime account** (any role ≥ CG Artist) for the bot to poll with — do not reuse a personal account

> **Bot permissions required:** Manage Channels, Manage Webhooks.
> Generate the invite URL from OAuth2 → URL Generator with scope `bot` and those two permissions.

---

## Step 1 — Clone and configure

```bash
git clone https://github.com/ukyovfx/kitsu-to-discord-task-notification.git
cd kitsu-to-discord-task-notification
cp .env.example .env.local
```

Open `.env.local` and fill in the minimum required values:

```env
KITSU_HOSTNAME=http://YOUR_KITSU_HOST/        # include http:// and trailing slash
KITSU_RUNTIME_EMAIL=bot@yourstudio.com
KITSU_RUNTIME_PASSWORD=your-runtime-password

DISCORD_BOT_TOKEN=your-bot-token
DISCORD_GUILD_ID=optional-fallback-server-id
```

Leave `DISCORD_WEBHOOK_URL` empty for now — the wizard sets up routing in the browser.

---

## Step 2 — Start the app

```bash
docker compose up -d --build
```

Wait for the app to be ready:

```bash
docker compose logs -f app
```

You should see within a few seconds:

```
Connected to Kitsu in ...
HTTP server listening on :8090
```

Verify health:

```bash
curl http://localhost:8090/health
# Expected: {"status":"ok"}
```

If you get a 502 or connection refused, the app is still starting — wait a moment and retry.

---

## Step 3 — Open the Setup Wizard

1. Open `http://YOUR_SERVER:8090/bot/login` in a browser.
2. Sign in with your **personal** Kitsu manager or admin account.
3. After login, open `/bot/admin/bot` then `/bot/admin/projects` first.

The wizard walks you through 4 steps:

| Step | What it does |
|------|-------------|
| 1. Kitsu | Verifies your Kitsu connection and authentication |
| 2. Discord | Confirms the bot token, server membership, and permissions |
| 3. Project | Creates Discord channels and webhooks for one Kitsu project in its assigned guild |
| 4. Mapping | (Optional) Maps Kitsu users to Discord IDs for @mentions |

Each step tests the connection live and shows errors inline before you proceed.

---

## Step 4 — Verify notifications

1. In Kitsu, change a task status to **WFA**, **Retake**, or **Done**.
2. Watch the logs: `docker compose logs -f app`
3. Look for `Notification route dispatch` followed by a send result.
4. Check the Discord channel — the notification should appear within one poll cycle.

---

## What's next

| Task | Where |
|------|-------|
| Monitor system health | `/bot/admin` (dashboard) |
| Assign guild per project | `/bot/admin/projects` |
| Edit channel routing | `/bot/setup` |
| Add more user/checker mappings | `/bot/admin/users`, `/bot/admin/checkers` |
| Set per-project storage links | `/bot/admin/drive` |
| Detailed setup reference | `docs/SETUP_WIZARD.md` |
| Troubleshooting | `docs/TROUBLESHOOTING.md` |

---

## Production deployment (Traefik + HTTPS)

For a production server using the Traefik stack in `deploy/`:

```bash
cp .env.local .env.production
# Edit .env.production: set PUBLIC_HOST and ALIAS
cd deploy
docker compose up -d
```

HTTPS and rate limiting are handled automatically by the Traefik labels in `deploy/docker-compose.yml`.

See `docs/ENVIRONMENTS.md` for full environment configuration details.
