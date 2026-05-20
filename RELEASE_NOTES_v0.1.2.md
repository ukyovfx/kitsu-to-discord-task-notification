# KitsuSync v0.1.2 Release Notes

## Overview

v0.1.2 is a small setup wizard trust/copy release.

It improves first-time operator clarity and confidence during setup, especially around Kitsu credentials, Discord bot permissions, provisioning confirmation, and retry guidance. It does **not** change runtime behavior or setup automation behavior.

## Highlights

### Setup wizard trust/copy improvements

- Added safety copy before Step 3 Discord provisioning so operators can see more clearly when Discord resources will actually be created
- Replaced raw Safe to retry: yes/no wording with human-readable retry guidance
- Clarified the purpose of the Kitsu credentials shown in Step 1
- Clarified what the Discord bot needs permission to do in Step 2

### Scope control

- No setup/auth/runtime/API/Docker/DB behavior change is intended in this release
- No setup flow redesign is included
- No new automation behavior is introduced

## Validation

- Individual v0.1.2 implementation PRs reported the following checks passing after each copy change:
  - docker compose config -q
  - docker compose build app
- This release note does not claim local go test / go vet execution unless run separately
- v0.1.2 should be understood as a trust/copy improvement release, not a runtime maturity jump

## Known limitations

- v0.1.2 improves setup clarity, but does not redesign the setup wizard
- Rollback remains best-effort; this release only makes that easier to understand
- Kitsu credential persistence behavior is unchanged
- Discord permission requirements are unchanged
- Setup/auth/runtime/API/Docker/DB behavior is unchanged

## Notes for users upgrading from v0.1.1

- No application behavior change is intended
- No data migration is required
- No setup/auth/runtime/API/Docker/DB reconfiguration is introduced by this release itself
- The main visible changes are clearer setup wording and more explicit operator guidance in the wizard
