# KitsuSync v0.2.0 Release Notes

## Overview

v0.2.0 is a Setup Wizard structure and operator clarity release.

It improves setup entry chooser clarity, renames the status overview surface from `Quick Setup` to `Setup Status`, reframes `System Check` as intentional Guided Setup preflight, makes Mapping read as optional post-setup follow-up work, and clarifies the role boundaries across the setup wizard and admin/setup surfaces. It does **not** change setup/auth/runtime/API/Docker/DB behavior. It does **not** change routes or remove manual setup paths. It does **not** add automation.

## Highlights

### Setup Wizard clarity

- Renamed the visible `Quick Setup` label to `Setup Status` while keeping the compatible `?mode=quick` route in place
- Clarified that `Setup Status` is for readiness/status review and does not replace Guided Setup
- Emphasized `Guided Setup` as the recommended first-time path in the wizard entry chooser

### Guided flow clarity

- Reframed `System Check` so it reads as Guided Setup preflight/status rather than a broken or hidden Step 0 state
- Clarified that Step 4 Mapping is optional post-setup/admin follow-up work
- Reinforced that the initial setup can already feel complete before optional mapping work begins

### Setup surface role clarity

- Clarified `/bot/setup-wizard` as the recommended first-time setup path
- Clarified `/bot/admin/setup` as manual setup / diagnostics / repair support
- Clarified `/bot/setup` as project/channel management rather than a competing first-time wizard
- Clarified `/bot/admin/bot` as shared bot/runtime credential configuration
- Clarified `/bot/admin/projects` as project-to-Discord server/guild assignment

### Scope control

- No setup/auth/runtime/API/Docker/DB behavior change is included
- No route renames or route removals are included
- No manual setup path removal is included
- No automation is included

## Validation

- Implementation PRs that changed app/UI copy reported `docker compose config -q` passing where run: #40, #44, #46, #48, #50
- Implementation PRs that changed app/UI copy reported `docker compose build app` passing where run: #40, #44, #46, #48, #50
- Docs-only terminology alignment work did not claim build/test execution: #42
- This release note does not claim local `go test` or `go vet` execution unless run separately
- v0.2.0 should be understood as a Setup Wizard structure / operator clarity release, not as a production maturity jump

## Known limitations

- v0.2.0 does not redesign setup architecture beyond wording, guidance, and surface-role clarity
- Manual setup paths remain necessary for some repair, diagnostics, and direct management cases
- Mapping remains optional in presentation, but it can still matter for richer `@mention` routing workflows
- `/bot/setup` remains an active management surface rather than being replaced or removed
- Setup/auth/runtime/API/Docker/DB behavior is unchanged

## Notes for users upgrading from v0.1.3

- No data migration is required
- No route migration is required
- Existing manual setup/admin paths remain available
- Existing `?mode=quick` compatibility remains in place even though the visible UI wording now says `Setup Status`
- The release improves how the setup surfaces explain themselves, but it does not introduce new automation or new runtime behavior
