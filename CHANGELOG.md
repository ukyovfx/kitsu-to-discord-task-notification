# Changelog

## [v1.0.0] — 2026-05-06

Initial release of the heavily customized fork. Based on [keshon/kitsu-to-discord-task-notification](https://github.com/keshon/kitsu-to-discord-task-notification).

### Security & Stability (Phase 1)
- **Fix**: Remove JWT token leak to Docker logs (was printed every polling cycle)
- **Fix**: Reverse message replace order to send-first-then-delete (prevents history loss on send failure)
- **Fix**: Replace `slog.Fatal` with `slog.Error` + continue in polling cycle (no crash on Kitsu API errors)
- **Fix**: Ensure HTTP response bodies are always closed (connection leak)
- **Fix**: Move Kitsu password and Discord Webhook URL to `.env` (out of `conf.toml`)

### Config & Reliability (Phase 2)
- **Fix**: `checkerStatuses` / `artistStatuses` in conf.toml are now actually read (were previously ignored)
- **Add**: `slog.Warn` when a Kitsu user or task type is not found in the mention maps
- **Add**: DB `SELECT + UPDATE` wrapped in GORM transactions (race condition fix for 12-goroutine pool)
- **Add**: 3-attempt retry with 2s→6s exponential backoff in `request.Do` (transient 5xx/429 tolerance)
- **Add**: Startup config validation — detects empty required fields and Japanese placeholder values
- **Fix**: `go.mod` version aligned to Go 1.21

### Features (Phase 3)
- **Add**: Kitsu preview thumbnail displayed in Discord embed image (`embed.image`)
- **Add**: `@here` support via `mention.hereStatuses` in conf.toml (urgent status broadcast)
- **Add**: Comment-only changes now trigger notifications (previously ignored)
- **Add**: Status transition messages (e.g. `RETAKE→DONE` → "再修正版がアップされました")
- **Fix**: `fields.tpl` was missing the opening `[` in the JSON array

### Advanced Features (Phase 4)
- **Add**: Task-type based channel routing via `[[discord.taskTypeWebhooks]]`
- **Add**: Per-task Discord threads via `discord.useThreads = true`
- **Add**: Daily digest at 09:00 JST — posts status-count summary to main webhook
- **Add**: Health check HTTP endpoint at `:8090/health`
- **Add**: `DiscordThreadID` column in SQLite for thread state persistence

### Publication Prep (Phase 5)
- **Add**: `conf.toml.example` — safe template with placeholder values and full comments
- **Add**: `.env.example` — template for secret values
- **Add**: `README.md` — Japanese (primary) + English (secondary) setup guide
- **Add**: `LICENSE` — Apache 2.0, original + fork copyright notice
- **Add**: GitHub Actions CI — `go vet` + `go build` on push/PR
