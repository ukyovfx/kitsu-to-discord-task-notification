# KitsuSync v0.3.1

## Summary
v0.3.1 is a small Discord notification UX refinement release. It preserves the compact rich Discord embed style while making WFA, RETAKE, and DONE messages clearer for production workflows.

## Highlights
- Refined compact Discord notification wording.
- Clarified WFA / RETAKE / DONE intent.
- Preserved transition-aware Japanese status messages.
- Kept Drive links and comments conditional.
- Added tests for rendered notification output and rich fields JSON validity.

## Changes

### Discord Notification Message UX
- PR #66
- Preserves compact rich embed style.
- Improves WFA / RETAKE / DONE wording.
- Keeps transition-aware status messages.
- Keeps Drive links visible only when available.
- Keeps comments visible only when present.

### Release Readiness Documentation
- PR #67
- Adds v0.3.1 release readiness note.

## Validation
- `go test ./src/api/discord/... -count=1`
- `go test ./src/...`
- `docker compose build app`
- docs-only validation for readiness note

## Known Limitations
- No routing behavior changes.
- No status filtering changes.
- No mention targeting changes.
- No setup/admin behavior changes.
- No broader notification workflow redesign.

## Upgrade Notes
- No DB migration required.
- No configuration change required.
- Existing Discord webhook routing remains unchanged.
- Notification wording is slightly changed for clarity.
