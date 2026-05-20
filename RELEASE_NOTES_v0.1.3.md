# KitsuSync v0.1.3 Release Notes

## Overview

v0.1.3 is a small docs cleanup and docs route consistency release.

It removes an orphaned docs artifact and aligns active in-app docs route aliases with the paths already referenced by the repository and admin UI. It does **not** redesign the docs architecture, and it does **not** change setup/auth/runtime/API/Docker/DB behavior.

## Highlights

### Docs cleanup

- Removed the orphaned `src/docs.html` artifact that was no longer part of the active docs flow
- Kept the active root `docs.html` entry point in place
- Kept `site.jsx`, `build_site.py`, and `diagrams/*.jsx` in place

### Docs route consistency

- Kept the existing `/docs` and `/site.jsx` routes for compatibility
- Added `/bot/docs`, `/bot/docs/`, and `/bot/docs/site.jsx` aliases so the active in-app docs links match Go routing behavior
- Improved consistency between repository/UI references to `/bot/docs/` and the actual docs-serving routes

### Scope control

- No docs architecture redesign is included
- No Markdown-for-HTML docs replacement is included
- No setup/auth/runtime/API/Docker/DB behavior change is intended in this release

## Validation

- #30 validation reported `git diff --check` passing
- #30 validation reported no remaining references to `src/docs.html`
- #30 validation confirmed the active root `docs.html` entry point was kept
- #32 validation reported `/bot/docs` and `docs.html` grep checks were reviewed
- #32 validation reported `docker compose config -q` passing
- #32 validation reported `docker compose build app` passing
- This release note does not claim local `go test` or `go vet` execution unless run separately
- v0.1.3 should be understood as a docs cleanup / docs route consistency release, not a production maturity jump

## Known limitations

- v0.1.3 does not redesign the docs architecture
- Markdown docs still do not replace the HTML docs because they serve different readers
- The active docs page still depends on the existing static `docs.html` plus `site.jsx` structure
- Setup/auth/runtime/API/Docker/DB behavior is unchanged

## Notes for users upgrading from v0.1.2

- No application behavior change is intended beyond the already-merged docs route alias consistency fix
- No data migration is required
- No setup/auth/runtime/API/Docker/DB reconfiguration is introduced by this release itself
- If you use in-app docs links or `/bot/docs/` paths behind your existing app/proxy setup, v0.1.3 aligns those paths with the Go-served docs aliases while keeping `/docs` and `/site.jsx` compatibility in place
