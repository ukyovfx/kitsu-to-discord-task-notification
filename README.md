# KitsuSync

[![Build](https://img.shields.io/badge/build-GitHub%20Actions-blue)](#)
[![Docker](https://img.shields.io/badge/runtime-Docker-2496ED)](#quick-start)
[![License](https://img.shields.io/badge/license-MIT-green)](#license)
[![Security Policy](https://img.shields.io/badge/security-policy-important)](#contributing-and-security)
[![Release](https://img.shields.io/badge/release-v0.1.0-orange)](#v010-release-focus)

KitsuSync is a Kitsu x Discord pipeline bridge for VFX and animation teams. It polls Kitsu, detects task changes, and posts structured Discord notifications with setup and admin tools in the browser.

**v0.1.0 focus:** secure first release, setup integrity, notification visibility, SQLite durability, and OSS onboarding clarity.

## Why This Repo Exists

KitsuSync helps teams that already track production in Kitsu but need Discord to act like an operator-friendly notification surface instead of a manual status relay.

- Poll Kitsu and post only meaningful task changes
- Create and manage Discord routing from the browser
- Keep setup, login, docs, and admin tooling in one service
- Stay deployable with Docker Compose and SQLite

## What It Solves

- Sends Kitsu task updates to Discord without manual copy/paste.
- Routes notifications by project, task type, or the default webhook.
- Helps operators create Discord channels and webhooks from `/bot/setup`.
- Keeps lightweight change state in SQLite so only new changes are posted.

## Architecture

```text
Kitsu API
  -> KitsuSync poller
  -> SQLite change tracking
  -> Discord webhook delivery

Browser admin UI
  -> /bot/login
  -> /bot/setup
  -> /bot/admin
```

## v0.1.0 Release Focus

- Security hardening for admin auth and runtime credential separation
- FileBrowser restricted to explicit debug usage
- Notification routing now surfaces silent-drop conditions in logs
- Setup flow no longer reports partial Discord provisioning as success
- SQLite now runs with WAL, busy timeout, and graceful shutdown logging

See `RELEASE_NOTES_v0.1.0.md` for the full release summary.

## Limitations (v0.1.0)

This release focuses on **single-project setups** with **notification routing and task tracking**. The following are **not** supported yet:

- **Multiple concurrent projects** — admin dashboard shows active projects but doesn't support project-scoped user/checker mappings (planned for v0.2.0)
- **Full Discord channel management** — KitsuSync creates channels but cannot delete or restructure them; manual cleanup required
- **All Kitsu task statuses** — filtering configured for `wfa`, `retake`, `done` and `assign` notifications only; other status transitions are silently dropped
- **Complete task detail retrieval** — reviews all task fields but does not embed full task description, comments history, or rich metadata
- **Large-scale queue handling** — no batching optimization for studios with 500+ task updates per minute (may hit Discord rate limits)
- **Enterprise authentication** — supports Kitsu Basic Auth only (no LDAP, OAuth, SSO)
- **Kitsu individual task URL links** — embeds generic Kitsu project URL, not direct task permalink (workaround: `/bot/docs/` link points to studio dashboard)

Most of these are roadmap items for v0.2.0. For production use, see `docs/SETUP_FOR_STUDIOS.md` and verify routing against your expected notification volume.

## Roadmap

The following items are intentionally **out of v0.1.0 scope** and are tracked for **v0.2.0+** rather than being partially shipped now:

- Project-scoped multi-project admin management
- Admin audit surface in the browser (`/bot/admin/audit`)
- Safer delete-and-recreate handling for existing Discord channel layouts
- Direct Kitsu task deep links in notifications
- UI controls for `@here` routing and broader mention policies
- More explicit setup dry-run and confirmation workflows
- Screenshot asset completion for public marketing/release pages

If you need these capabilities, keep them in roadmap planning instead of extending the v0.1.0 release branch.

## UI Surfaces

- `/bot/login` — admin sign-in
- `/bot/setup-wizard` — guided 4-step first-time setup (recommended entry point)
- `/bot/setup` — direct project and channel management
- `/bot/admin` — operational dashboard: system health, active projects, warnings
- `/bot/admin/users` — map Kitsu users to Discord IDs for @mentions
- `/bot/admin/checkers` — map task types to reviewer Discord IDs
- `/bot/docs/` — operator-facing pipeline documentation

Screenshot placeholders and capture guidance live in `screenshots/README.md`.

## Getting Started

See `docs/QUICK_START.md` for a 5-minute startup guide.
See `docs/SETUP_WIZARD.md` for a detailed walkthrough of each wizard step.

## Repository Layout

```text
docker-compose.yml        App + optional FileBrowser(debug profile only)
.env.example              Environment variable template
conf.toml.example         App configuration template
docs.html / site.jsx      Browser docs page content
src/                      Go application source
tpl/                      Discord message templates
diagrams/                 Supporting docs assets
```

## Requirements

- Docker
- Docker Compose
- A reachable Kitsu server
- A Discord server where you can create webhooks
- A reverse proxy in production if you want `/bot/*` under the same public host

## Quick Start

### 1. Clone and prepare config

```bash
git clone https://github.com/YOUR_ORG/kitsusync.git
cd kitsusync
cp .env.example .env.local
cp conf.toml.example conf.toml
mkdir -p data
```

### 2. Fill in `.env.local`

At minimum, set:

- `DISCORD_BOT_TOKEN`
- `DISCORD_GUILD_ID`

You have two ways to provide Kitsu runtime credentials:

1. Recommended for reproducible local bring-up
   - Fill `KITSU_HOSTNAME`
   - Fill `KITSU_RUNTIME_EMAIL`
   - Fill `KITSU_RUNTIME_PASSWORD`
2. Guided browser setup
   - Fill only `KITSU_HOSTNAME`
   - Start the app
   - Open `/bot/login` and create/apply the runtime bot account from `/bot/setup`

If you want a default catch-all Discord route before project setup, also set:

- `DISCORD_WEBHOOK_URL`

### 3. Review `conf.toml`

`conf.toml.example` already points secret values at environment variables. Common values to review:

- `kitsu.hostname`
- `discord.useThreads`
- `mention.checkerStatuses`
- `mention.artistStatuses`
- `mention.hereStatuses`

### 4. Start the app

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f app
```

Health check:

```bash
curl http://localhost:8090/health
```

Expected response:

```json
{"status":"ok"}
```

If you want a quick release-readiness sanity check after boot:

```bash
docker compose ps
docker compose logs --tail=50 app
curl http://localhost:8090/health
```

## First-Time Setup Flow

### Direct local access

- Login page: `http://localhost:8090/bot/login`
- Docs: `http://localhost:8090/bot/docs/`
- Setup: `http://localhost:8090/bot/setup`
- Admin: `http://localhost:8090/bot/admin`

### Behind a reverse proxy

Use the public `/bot/*` paths exposed by your proxy, for example:

- `/bot/login`
- `/bot/docs/`
- `/bot/setup`
- `/bot/admin`

### Operator checklist

1. Open `/bot/login` and sign in with your Kitsu manager or admin account.
2. If runtime credentials are not already in `.env.local`, run Bot Setup from `/bot/setup`.
3. Create or choose a Kitsu project from Project Setup.
4. Confirm the Discord category, channels, and webhooks are created.
5. Review routing and user mappings in `/bot/admin`.
6. Wait for polling logs such as `Connected to Kitsu`, `Got tasks`, and `Done FilterTasks`.

## Environment Variables

See `.env.example` for the full template. Copy it to `.env.local` (development) or `.env.production` (production) — never commit these files to git.

### Required for most installs

- `KITSU_HOSTNAME`
- `DISCORD_BOT_TOKEN`
- `DISCORD_GUILD_ID`

### Required if you want polling to work immediately on first boot

- `KITSU_RUNTIME_EMAIL`
- `KITSU_RUNTIME_PASSWORD`

### Optional

- `DISCORD_WEBHOOK_URL`
- `FB_USERNAME` (debug profile only)
- `FB_PASSWORD` (debug profile only — generate with: `openssl rand -base64 20`)

## Debug vs Production

### Production defaults

- `app` starts by default.
- `editor` does not start by default.
- Runtime secrets must stay in `.env.production` or your deployment secret store. Never commit secret files to git.
- The app is expected to run behind a trusted reverse proxy if you expose `/bot/*` publicly.

### Debug profile

Start FileBrowser only when you explicitly need it:

```bash
docker compose --profile debug up -d editor
```

FileBrowser is for local/debug inspection only.

- It is not recommended for production.
- It does not mount `.env`.
- It does not mount `conf.toml`.
- It does not mount `sqlite.db`.
- It does not mount runtime credential storage.

## Routing Behavior

Notification routing priority is:

1. Project/task-type webhook records created by `/bot/setup`
2. `[[discord.productions]]`
3. `[[discord.taskTypeWebhooks]]`
4. `discord.webhookURL` / `DISCORD_WEBHOOK_URL` as the main fallback

If no fallback webhook is configured for unmatched tasks, the app now logs the drop explicitly instead of silently swallowing it.

## Production Notes

- Cookies are hardened for HTTPS and trusted reverse proxy operation.
- `SameSite=Lax` is intentionally used to preserve login redirect flows.
- `X-Forwarded-Proto` is trusted only when your reverse proxy overwrites it.
- Do not expose the app directly to the internet without a properly configured reverse proxy.
- Large first-time syncs can temporarily hit Discord rate limits. For an existing Kitsu with many tasks, start with `ignoreMessagesDaysOld = 7` and widen it after setup is stable.
- If nginx proxies `/bot/`, set `proxy_read_timeout 300;` so long Discord setup calls do not appear as false browser-side 504 failures.

Example `/bot/` nginx snippet:

```nginx
location /bot/ {
    proxy_pass http://127.0.0.1:8090;
    proxy_read_timeout 300;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

## Preview Images

Discord preview thumbnails require Kitsu preview files to be reachable without interactive browser auth. If your Kitsu deployment protects preview files differently, add a reverse-proxy rule for the preview endpoint.

Example nginx snippet:

```nginx
location ~ ^/api/pictures/thumbnails/preview-files/ {
    proxy_pass http://localhost:5000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Troubleshooting

### `Initial Kitsu authentication failed`

- Check `KITSU_HOSTNAME`
- Check `KITSU_RUNTIME_EMAIL`
- Check `KITSU_RUNTIME_PASSWORD`
- Or complete Bot Setup from `/bot/setup`

### Notifications are not arriving

- Check `docker compose logs -f app`
- Check `DISCORD_WEBHOOK_URL`
- Check `/bot/admin` routing
- Check the new notification routing warnings in the logs

### Setup failed partway through

- The setup screen now reports failure explicitly.
- Partial Discord resources are rolled back on a best-effort basis.
- Re-run setup only after reading the returned `FAIL:` / `WARN:` lines.

### FileBrowser is not reachable

- This is expected unless you started the debug profile explicitly.

## Documentation

- Template variables and custom preset guide: `docs/TEMPLATES.md`
- Dev vs production environment setup: `docs/ENVIRONMENTS.md`
- First-time studio setup walkthrough: `docs/SETUP_FOR_STUDIOS.md`
- Error messages and diagnostics: `docs/TROUBLESHOOTING.md`
- Repo-local release notes and scope sync: `docs/notes/`

## Contributing and Security

- Contributor guide: `CONTRIBUTING.md`
- Security reporting: `SECURITY.md`
- Changelog: `CHANGELOG.md`
- Release notes: `RELEASE_NOTES_v0.1.0.md`
- Screenshot guidance: `screenshots/README.md`

## License

MIT. Keep the upstream notices when redistributing.
