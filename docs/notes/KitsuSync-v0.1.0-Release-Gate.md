---
title: "KitsuSync v0.1.0 Release Gate"
status: release-candidate
updated: 2026-05-19
release_gate: pending-verification
---

# KitsuSync v0.1.0 Release Gate

## Must Pass

- tracked secrets / runtime data / deploy artifacts excluded from release diff
- `.gitignore` blocks common local runtime outputs
- README includes quick start, limitations, and roadmap
- CURRENT-STATE / Technical Spec / Release Readiness use aligned status vocabulary
- `/api/setup/status` smoke test passes
- auth session smoke test passes
- notification routing regression test passes
- log redaction test passes
- clean-clone startup path is documented

## Manual Checks

- real Kitsu connectivity
- real Discord bot/guild validation
- actual browser walkthrough for `/bot/login`, `/bot/setup-wizard`, `/bot/setup`, `/bot/admin`, `/bot/docs/`
- screenshot capture and scrubbing
- presence of `sanitized_kitsu_schema.sql` and `sanitized_kitsu_sample.json` if they are part of the release packet

## Current Verdict

- automated gate: pending
- manual gate: pending
- ship decision: hold until both sections are reviewed
