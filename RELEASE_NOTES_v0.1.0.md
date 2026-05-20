# KitsuSync v0.1.0 Release Notes

## Summary

v0.1.0 is the first public-ready release of KitsuSync.

This release focuses on making the project safer to operate, easier to onboard, and easier to trust as an OSS repository without changing the core Kitsu -> Discord runtime model.

## What Is Included

### Security

- Hardened admin cookie handling for HTTPS and trusted reverse proxy operation
- Separated admin login credentials from runtime polling credentials
- Restricted FileBrowser to an explicit debug-only profile with secrets excluded from mounts

### Reliability

- Notification routing now logs route dispatch, send results, and unmatched drop conditions
- Setup flow no longer reports partial Discord provisioning as success
- Best-effort rollback added for setup failures during channel/webhook creation

### Setup and DX

- Added onboarding-focused `README.md`
- Added `.env.example` and `conf.toml.example`
- Added `CONTRIBUTING.md` and `SECURITY.md`
- Added GitHub issue templates, PR template, CODEOWNERS, and release checklist

### SQLite Durability

- SQLite now configures WAL mode at startup
- `busy_timeout` is applied to reduce transient contention failures
- Graceful shutdown now closes cron, HTTP server, and SQLite connections cleanly

## Architecture Snapshot

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

## Production Assumptions

- Run behind a trusted reverse proxy if exposing `/bot/*`
- Keep secrets out of git and out of FileBrowser mounts
- Use `SameSite=Lax` cookies intentionally to preserve login redirect flows
- Treat FileBrowser as debug-only, not as a production operations UI

## Known Limitations

- SQLite remains the only supported runtime database in v0.1.0
- Setup rollback is best-effort, not a full distributed transaction
- Discord routing verification still depends on live webhook reachability
- Preview image behavior depends on your Kitsu/reverse-proxy exposure rules

## Suggested Release Summary

KitsuSync v0.1.0 delivers a safer first public release with auth hardening, runtime credential separation, setup rollback, notification observability, SQLite durability improvements, and improved OSS onboarding.
