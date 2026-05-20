# KitsuSync v0.1.1 Release Notes

## Overview

v0.1.1 is a post-release hardening follow-up to v0.1.0.

This release improves repository safety, CI coverage, public onboarding clarity, and project policy documentation. It does **not** expand runtime behavior or change the core Kitsu -> Discord application flow.

## Highlights

### CI and validation

- Added a GitHub Actions CI baseline for:
  - go vet ./src/...
  - go test ./src/... -count=1 -timeout=120s
  - docker compose config -q
  - docker compose build app
- Added CI-safe .env.local preparation so Docker Compose validation can run in GitHub Actions without relying on a host machine setup

### Public docs and onboarding

- Corrected public clone/setup docs to use the actual repository path:
  - https://github.com/ukyovfx/kitsu-to-discord-task-notification.git
- Reviewed first-time setup UX from a beginner operator perspective
- Aligned setup wizard docs with the current product behavior:
  - /bot/setup-wizard is the recommended first-time entry point
  - the first screen may show a setup mode chooser
  - System Check may appear before Kitsu/Discord/project setup
  - connection testing, persistent configuration, and Discord resource creation are now described more clearly in the docs
  - rollback language is more cautious and matches the current best-effort behavior

### Repository policy notes

- Added a historical tag policy for inherited non-v tags
- Added a branding policy for the relationship between the public name KitsuSync and the historical repository slug kitsu-to-discord-task-notification
- Added a v0.1.1 release readiness note documenting why this release is considered safe to ship as a hardening/docs/policy release

## Validation

- GitHub Actions CI is available after the v0.1.1 hardening work
- This release note does not claim local go test / go vet execution on the host unless explicitly run elsewhere
- v0.1.1 is intended as a docs/CI/policy hardening release, not a runtime maturity jump

## Known limitations

- v0.1.1 does not change setup/auth/runtime/API/Docker/DB behavior
- Setup wizard trust gaps were reduced in docs, but not fully solved in-product
- Branch protection and ruleset behavior should still be monitored over time
- Historical non-v tags remain in place for now by policy
- Repository slug and public product name remain intentionally different for now

## Notes for users upgrading from v0.1.0

- No application behavior change is intended
- No data migration is required
- No setup/auth/runtime/API/Docker/DB reconfiguration is introduced by this release itself
- The most visible differences are CI coverage, docs clarity, and repository governance notes
