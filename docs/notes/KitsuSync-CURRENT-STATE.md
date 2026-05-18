---
title: "KitsuSync Current State"
status: release-candidate
updated: 2026-05-19
release_gate: pending-verification
---

# KitsuSync Current State

## Release Focus

- target release: `v0.1.0`
- current mode: release candidate hardening
- scope rule: do not expand beyond the existing setup, routing, and operator surfaces

## Implemented Surface

- `/bot/login`: `done`
- `/bot/setup-wizard`: `partial`
- `/bot/setup`: `partial`
- `/bot/admin`: `partial`
- `/bot/admin/users`: `partial`
- `/bot/admin/checkers`: `partial`
- `/api/setup/status`: `partial`
- `/bot/admin/audit`: `v0.2.0`

## Main Risks

- release evidence still needs manual infrastructure checks
- clean-clone bring-up cannot be fully proven without real Kitsu and Discord credentials
- screenshot assets are still manual TODO items
- `sanitized_kitsu_schema.sql` and `sanitized_kitsu_sample.json` were not present in this repo-local review copy

## Next Gate

Release passes only after `docs/notes/KitsuSync-v0.1.0-Release-Gate.md` is marked pass for all must-pass items.
