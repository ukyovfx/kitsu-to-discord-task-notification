# Setup Wizard Reference

The Setup Wizard (`/bot/setup-wizard`) is the **recommended first-run entry point** for v0.1.0. It guides you through the complete first-time configuration in 4 steps. This document covers each step in detail, including common errors and how to fix them.

If you are starting from a clean clone, the intended operator path is:

1. Start the app with `.env.local` and `conf.toml`
2. Open `/bot/login`
3. Sign in with a personal Kitsu manager/admin account
4. Continue into `/bot/setup-wizard`
5. Use `/bot/setup` and `/bot/admin/*` only after the wizard finishes

---

## Overview

| Step | Name | Required | What it configures |
|------|------|----------|--------------------|
| 1 | Kitsu Connection | ✅ Yes | Kitsu hostname, email, password |
| 2 | Discord Bot | ✅ Yes | Bot token, Guild ID, permissions |
| 3 | Project Setup | ✅ Yes | Discord channels and webhooks for one Kitsu project |
| 4 | User & Checker Mapping | Optional | @mention routing for users and reviewers |

The wizard pre-populates steps from the current system state. If Kitsu and Discord are already configured, the wizard opens directly at Step 3 or 4.

---

## Step 1: Kitsu Connection

**What it does:** Verifies that KitsuSync can reach your Kitsu server and authenticate with the runtime account credentials.

**Inputs:**

| Field | Example | Notes |
|-------|---------|-------|
| Kitsu Hostname | `http://kitsu.studio.local/` | Must include `http://` or `https://` and a trailing slash |
| Email | `bot@studio.com` | The dedicated runtime Kitsu account, not your personal login |
| Password | — | The runtime account password |

**Checking connection** calls `POST /api/setup/test-kitsu`. The API:
1. Pings `{hostname}/api/` — if unreachable, reports "server not reachable"
2. Authenticates at `{hostname}/api/auth/login` with the credentials
3. Returns success only if both succeed

**On success:** The step collapses and shows a green `✓ Kitsu: Authenticated` badge. The wizard advances to Step 2.

### Step 1 errors

| Error | Cause | Fix |
|-------|-------|-----|
| `server not reachable` | Kitsu is not accessible from the KitsuSync container | Check `KITSU_HOSTNAME`; test from the server with `curl http://YOUR_KITSU_HOST/api/` |
| `authentication failed — check email and password` | Credentials are wrong | Verify the runtime account exists in Kitsu and is active |
| `hostname must start with http:// or https://` | URL format error | Add the scheme to the hostname |
| `All fields are required` | A field is empty | Fill in all three fields |

---

## Step 2: Discord Bot

**What it does:** Verifies the bot token, confirms the bot is in your server, and checks that it has the required permissions.

**Inputs:**

| Field | Example | Notes |
|-------|---------|-------|
| Bot Token | `MTUwMTUx...` | From Discord Developer Portal → Bot tab → Reset Token |
| Guild ID | `1234567890123456789` | Right-click server name in Discord (Developer Mode must be on) |

**Checking the bot** calls `POST /api/setup/test-discord`. The API:
1. Validates the bot token against Discord's `/users/@me`
2. Confirms the guild is accessible via `/guilds/{guild_id}`
3. Checks the bot's membership and permission bits (Manage Channels, Manage Webhooks)

**On success:** Shows `✓ Discord: BotName / ServerName` and advances to Step 3.

### Step 2 errors

| Error | Cause | Fix |
|-------|-------|-----|
| `bot token invalid (HTTP 401)` | Token is wrong or revoked | Generate a new token from Discord Developer Portal |
| `guild not accessible (HTTP 404)` | Guild ID is wrong or the bot is not in the server | Check the Guild ID; invite the bot using the OAuth2 URL |
| `Missing permissions: MANAGE_CHANNELS` | Bot lacks the permission | Edit the bot role in Discord Server Settings → Roles |
| `Missing permissions: MANAGE_WEBHOOKS` | Bot lacks the permission | Same as above |

> **Tip:** To invite the bot with the correct permissions, go to Discord Developer Portal → OAuth2 → URL Generator. Select scope `bot` and permissions `Manage Channels` + `Manage Webhooks`. Open the generated URL to add the bot to your server.

---

## Step 3: Project Setup

**What it does:** Creates a Discord category, text channels, and webhooks for one Kitsu project using the `cg` template.

**Inputs:**

| Field | Options | Notes |
|-------|---------|-------|
| Project | List loaded from Kitsu | Projects already configured show as "configured" |
| Notification language | Japanese / English | Controls channel names; can be changed per-project |

**Creating channels** calls `POST /api/setup/apply-project`. The API:
1. Resolves the project name from Kitsu
2. Creates a Discord category named after the project
3. Creates text channels for each task type defined in the template
4. Creates a Discord webhook in each channel
5. Records all channels and webhook URLs in the database

Each step is atomic — if any channel creation fails, the entire setup is rolled back automatically. You can re-run setup after fixing the issue.

**On success:** Shows channel count and webhook count, then advances to Step 4.

### Step 3 errors

| Error | Cause | Fix |
|-------|-------|-----|
| `project not found in Kitsu` | Project was deleted from Kitsu after the list was loaded | Reload the wizard and try again |
| `unsupported template: ...` | Invalid template name | Use `cg` (the only supported template currently) |
| `FAIL: failed to create Discord category` | Bot lacks Manage Channels | Fix bot permissions (see Step 2) |
| `FAIL: failed to create webhook` | Bot lacks Manage Webhooks | Fix bot permissions (see Step 2) |
| `✅ Safe to retry` badge | Setup failed but no partial state remains | Fix the reported error and click Create Channels again |

> **Tip:** If you see `WARN:` lines but the setup completes (OK), the warnings are non-fatal. Common warning: a channel name already exists and was reused instead of created.

### Already-configured projects

Projects that already have channels show in a "Configured" list at the top of Step 3. If all your projects are already configured, the wizard shows a completion message and offers a link to Step 4 for mapping.

To add a second project, come back to the wizard after completing the first.

---

## Step 4: User & Checker Mapping (Optional)

> **This step is optional.** If you skip it, notifications are still delivered. Only @mentions will be missing.

**What it does:** Maps Kitsu user accounts to Discord user IDs (for @mention on task assignment), and maps task types to reviewer Discord IDs (for @mention when status changes to WFA).

**Loading mapping data** calls `GET /api/setup/mapping`, which returns:
- All active Kitsu persons
- Task types that have Discord channels (from Step 3)
- Any existing mappings already saved

### User mapping

For each Kitsu person, enter their Discord User ID in the input field.

**How to find a Discord User ID:**
1. In Discord, enable Developer Mode (User Settings → Advanced).
2. Right-click the user's name → Copy User ID.

Leave the field empty to skip that person. Entering an empty ID where one was previously set will remove the mapping.

### Checker mapping

For each task type that has a Discord channel, enter the Discord User ID of the person who reviews that task type.

When a task transitions to **WFA**, the checker for that task type is @mentioned in the notification.

### Saving

Click **Save & Finish** to write all mappings at once. The wizard calls:
- `POST /api/setup/mapping/users` — saves user → Discord ID entries
- `POST /api/setup/mapping/checkers` — saves task type → Discord ID entries

Entries with an empty ID field are deleted from the database.

Click **Skip** to go to the Done screen without saving.

---

## After the wizard

Once the wizard completes, the **Dashboard** (`/bot/admin`) shows:

- **System status** — whether the poller is running and the last sync time
- **Active Projects** — channel and webhook counts per project
- **Warnings** — broken webhooks or poller errors that need attention

### Ongoing management

| Task | Where |
|------|-------|
| Add or remove channels for a project | `/bot/setup` — manage individual channels |
| Edit user/checker mappings | `/bot/admin/users`, `/bot/admin/checkers` |
| Reconnect a broken webhook | `/bot/admin/health` → Reconnect button |
| Add a second project | Run the wizard again at `/bot/setup-wizard` |
| Check detailed poller stats | `/bot/admin/health` |

---

## Verifying the API directly

If you need to diagnose setup state programmatically:

```bash
# Full system status (requires session cookie or Basic Auth)
curl http://localhost:8090/api/setup/status

# Test Kitsu credentials without saving
curl -X POST http://localhost:8090/api/setup/test-kitsu \
  -H "Content-Type: application/json" \
  -d '{"hostname":"http://kitsu.local/","email":"bot@studio.com","password":"xxx"}'

# Test Discord bot
curl -X POST http://localhost:8090/api/setup/test-discord \
  -H "Content-Type: application/json" \
  -d '{"bot_token":"xxx","guild_id":"yyy"}'
```

The `/api/setup/status` response is the authoritative source of truth for the current setup state — all wizard UI logic reads from it.
