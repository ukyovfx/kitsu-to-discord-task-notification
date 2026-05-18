---
title: "KitsuSync v0.1.0 Release Readiness"
status: release-candidate
updated: 2026-05-19
release_gate: pending-verification
---

# KitsuSync v0.1.0 Release Readiness

## Summary

Implementation is close enough to evaluate as a release candidate, but final release depends on release-gate evidence rather than a percentage score.

## Readiness by Area

- tracked secrets / runtime data cleanup: `partial`
- README / docs alignment: `partial`
- setup wizard operator guidance: `partial`
- `/api/setup/status` coverage: `partial`
- auth session coverage: `partial`
- notification routing regression coverage: `partial`
- log redaction coverage: `partial`
- screenshots for public release: `human-check`
- clean-clone boot validation with real infra: `human-check`

## Deferred to v0.2.0

- admin audit UI
- delete and recreate safety flows
- direct task deep-link hardening
- `@here` management UI

## Decision Rule

Do not ship based on readiness percentage alone. Ship only when the must-pass section in `KitsuSync-v0.1.0-Release-Gate.md` is green or explicitly waived by a human release owner.
