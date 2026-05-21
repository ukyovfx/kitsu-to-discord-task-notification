# KitsuSync v0.3.0

## Summary
v0.3.0 is an operational hardening release focused on logging safety, safe setup testing behavior, and cleanup for partial failures during Discord provisioning.

## Highlights
- Secret redaction is now wired into stdout and file logs.
- Discord connection testing no longer mutates runtime or persisted setup state.
- Discord resources created during setup are cleaned up if the DB step fails afterward.
- v0.3.0 readiness note documents validation, known gaps, and non-goals.

## Changes

### Logging and Secret Safety
Implemented in PR #59. Redaction was wired into slog output paths so both console output and file logs are protected, including runtime env-aware secret redaction and hardened writer behavior.

### Setup Test Safety
Implemented in PR #60. `/api/setup/test-discord` is now validation-only and no longer updates runtime env values or persists setup state as a side effect.

### Partial Failure Cleanup
Implemented in PR #62. When Discord resources are created but the subsequent DB transaction fails, setup now performs request-scoped best-effort cleanup of created channels/category without changing success-path behavior.

### Release Readiness Documentation
Implemented in PR #63. Added `docs/notes/v0.3.0-release-readiness.md` to record validation, known gaps, and non-goals for this hardening phase.

## Validation
- `go test ./src/setup/... -count=1`
- `go test ./src/...`
- `docker compose build app`
- docs-only validation for the readiness note content

## Known Limitations
- no full rollback orchestration
- no persistent leftover-resource tracking
- limited repair UI
- token persistence design unchanged
- no release/tag behavior changes yet

## Upgrade Notes
- No DB migration required.
- No configuration change required.
- Operators should still treat `.env.local` / environment as the source of truth for durable bot token rotation.
- Setup test actions are now validation-only.
