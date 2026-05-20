# KitsuSync v0.2.1 Release Notes

## Overview

v0.2.1 is a repository rename and public URL alignment follow-up to v0.2.0.

This release aligns the GitHub repository slug with the public product name KitsuSync, updates current public clone URLs and operator-facing repository references, and adds a rename decision note for historical clarity. It does **not** change runtime behavior. It does **not** change setup/auth/runtime/API/Docker/DB behavior. It should **not** be understood as a production maturity jump.

## Highlights

### Repository rename

- Renamed the GitHub repository from `ukyovfx/kitsu-to-discord-task-notification` to `ukyovfx/kitsusync`
- Aligned the canonical repository slug with the public product name KitsuSync
- Preserved the historical context around the older slug instead of rewriting prior notes and release-era records

### Public docs and onboarding URLs

- Updated current public clone URLs to use:
  - `https://github.com/ukyovfx/kitsusync.git`
- Updated current `cd` targets in onboarding docs to use:
  - `cd kitsusync`
- Kept current public docs aligned with the new canonical repository path instead of relying on redirects

### Rename decision note

- Added `docs/notes/repo-rename-v0.2.1.md`
- Documented the rename reason, canonical repository URLs, redirect guidance, and local clone follow-up
- Clarified that older historical notes remain historical records and should not be silently rewritten

## Validation

- Rename follow-up validation in #56 checked `git remote -v`, `git status --short --untracked-files=all`, `git diff --stat`, `git diff`, and old-slug grep results after the docs/remote update
- #56 confirmed that remaining old-slug references were historical only after the current public clone URLs were updated
- This release note does not claim local `go test` or `go vet` execution unless run separately
- This release note does not claim Docker build/test execution
- v0.2.1 should be understood as a repository rename / public URL alignment release, not as a runtime or production maturity jump

## Known limitations

- v0.2.1 does not change setup, login, admin, mapping, or notification behavior
- GitHub redirects may help older URLs continue working for a time, but redirects should not be treated as the long-term canonical path in docs or operator instructions
- Historical notes and earlier release documents still mention the prior repository slug where that wording accurately describes earlier project state
- Setup/auth/runtime/API/Docker/DB behavior remains unchanged

## Notes for users upgrading from v0.2.0

- No data migration is required
- No setup/auth/runtime/API/Docker/DB reconfiguration is introduced by this release itself
- Existing local clones should update `origin` to:
  - `git@github.com:ukyovfx/kitsusync.git`
- The most visible differences are repository slug alignment, current public clone URL updates, and the added rename decision/release-readiness documentation
