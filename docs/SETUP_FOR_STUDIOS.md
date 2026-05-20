# Studio Setup Guide

This guide walks through the complete setup of KitsuSync for a studio that already uses Kitsu and Discord. No deep engineering knowledge required.

## What You Need Before Starting

- **Kitsu** running and reachable from the server
- **Discord server** where you are an administrator
- A **server or VM** with Docker and Docker Compose installed
- A **Kitsu account** with `manager` or `admin` role (for logging into KitsuSync)
- A **dedicated Kitsu runtime account** (a bot/service account used for polling — see below)

---

## Part 1: Create a Runtime Account in Kitsu

KitsuSync needs a dedicated Kitsu account to poll for changes. This should be a service account, not a personal account.

1. Log into Kitsu as admin.
2. Go to **Admin → People** → Add a new person.
3. Set the role to **CG Artist** or **Manager** (Manager lets it see all projects).
4. Note the email and password — you'll need them for `.env.local`.

> Using a dedicated account means you can rotate the password without affecting real users, and you can see exactly which API calls come from the bot in Kitsu audit logs.

---

## Part 2: Add the Bot to Your Discord Server

1. Go to [Discord Developer Portal](https://discord.com/developers/applications).
2. Create a new application, then go to the **Bot** tab.
3. Click **Reset Token** and copy the bot token — save it securely.
4. Under **Privileged Gateway Intents**, no special intents are needed.
5. Go to **OAuth2 → URL Generator**.
6. Select scopes: `bot`
7. Select permissions: `Manage Channels`, `Manage Webhooks`
8. Copy the generated URL and open it in a browser to add the bot to your server.

Note your **Discord Guild IDs** (Server IDs) for each production server:

- In Discord, go to Settings → Advanced → enable Developer Mode.
- Right-click your server name → Copy Server ID.

---

## Part 3: Deploy KitsuSync

### On the server

```bash
git clone https://github.com/ukyovfx/kitsu-to-discord-task-notification.git
cd kitsu-to-discord-task-notification

# Copy and fill in the config files
cp .env.example .env.local
cp conf.toml.example conf.toml
mkdir -p data
```

### Fill in `.env.local`

Open `.env.local` in a text editor and fill in:

```env
KITSU_HOSTNAME=http://YOUR_KITSU_HOST/      # Include http:// and trailing slash
KITSU_RUNTIME_EMAIL=bot@yourstudio.com      # The runtime account you created in Part 1
KITSU_RUNTIME_PASSWORD=yourpassword

DISCORD_BOT_TOKEN=your-bot-token            # From Discord Developer Portal (Part 2)
DISCORD_GUILD_ID=optional-fallback-server-id  # Optional fallback only

DISCORD_WEBHOOK_URL=                        # Leave blank for now; set up via /bot/setup
```

### Review `conf.toml`

The most important settings at first:

```toml
tplPreset = "rich"            # "rich" (Japanese), "eng" (English), or "rus" (Russian)
ignoreMessagesDaysOld = 5     # How many days back to pick up missed changes on first boot

[discord]
  useThreads = false          # Set true to post updates as thread replies instead of new messages
```

Leave other settings at their defaults for now.

### Start the app

```bash
docker compose up -d --build
docker compose logs -f app
```

Wait until you see:

```
Connected to Kitsu in ...
HTTP server listening on :8090
```

Verify it is healthy:

```bash
curl http://localhost:8090/health
```

Expected: `{"status":"ok"}`

---

## Part 4: First-Time Setup in the Browser

Open a browser and go to `http://YOUR_SERVER:8090/bot/login`.

If you are using a reverse proxy (nginx, Traefik), use your public URL instead.

### Step 1: Log in

Sign in with your **personal** Kitsu manager or admin account. This is not the runtime account — it's your own login used to operate the admin UI.

Sessions expire after 15 minutes.

### Step 2: Bot Setup (`/bot/admin/bot`)

This step verifies that the bot token and guild ID are working.

1. Open `/bot/setup`.
2. Click **Bot Setup** (or the equivalent button in your language).
3. If successful, you'll see the bot username and server name confirmed.

### Step 3: Assign Guild per Project (`/bot/admin/projects`)

For each Kitsu project, set the Discord Guild ID where channels should be created.

### Step 4: Project Setup (`/bot/setup` → Project Setup)

This creates Discord channels and webhooks for one Kitsu project.

1. Select a Kitsu project from the dropdown.
2. Choose a project type (`cg` is currently the only supported type).
3. Choose the language for channel names (Japanese or English).
4. Select how to handle Discord channel routing:
   - **Create new category** — KitsuSync creates a new category and all channels under it.
   - **Use existing category** — choose an existing Discord category to place channels in.
5. Click **Run Setup**.

If setup succeeds, you'll see channel names confirmed in the output.
If setup fails partway through, read the `FAIL:` / `WARN:` lines in the output — partial changes are rolled back automatically.

### Step 5: Review Routing (`/bot/admin`)

After project setup, open `/bot/admin` to verify:

- **Routing** tab: Shows which Kitsu task types send to which Discord channels.
- **Users** tab: Map Kitsu usernames to Discord user IDs so the bot can `@mention` people.
- **Checkers** tab: Map task types to checker Discord IDs so the bot mentions reviewers when status changes to WFA.

---

## Part 5: Verify Notifications Are Working

1. In Kitsu, change a task status on a project you set up.
2. Watch the Docker logs: `docker compose logs -f app`
3. You should see `Notification route dispatch` followed by `SendMessage`.
4. Check the Discord channel — the notification should appear within one poll cycle (a few seconds).

If notifications are not appearing, see `docs/TROUBLESHOOTING.md`.

---

## Part 6: Production Deployment (with Traefik)

For a production server using the Traefik setup in `deploy/docker-compose.yml`:

1. Copy `.env.local` to `.env.production` and update values as needed.
2. Set `PUBLIC_HOST=kitsusync.example.com` and `ALIAS=kitsusync` in `.env.production`.
3. Ensure Traefik is running and the `proxy` Docker network exists.
4. Deploy:

```bash
cd deploy
docker compose up -d
```

HTTPS and rate limiting are handled automatically by Traefik using the labels in `deploy/docker-compose.yml`.

---

## Updating KitsuSync

```bash
git pull
docker compose up -d --build
```

Config files (`conf.toml`, `.env.local`) are not touched by git pull — your settings are preserved.

---

## Common First-Setup Issues

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| Login fails immediately | Kitsu credentials or URL wrong | Check `KITSU_HOSTNAME`, `KITSU_RUNTIME_EMAIL`, `KITSU_RUNTIME_PASSWORD` |
| Bot Setup fails | Bot token or Guild ID wrong | Double-check both in `.env.local` |
| Project Setup fails | Bot lacks Discord permissions | Verify bot has Manage Channels and Manage Webhooks |
| Notifications not arriving | No fallback webhook | Set `DISCORD_WEBHOOK_URL` or ensure project routing is active |
| `/health` returns 502 | App not started | Run `docker compose up -d --build` |

For more detailed troubleshooting, see `docs/TROUBLESHOOTING.md`.
