# Historical Tag Policy for v0.1.1

## Context
KitsuSync now has a public release line starting at v0.1.0. The repository history also contains older non-semantic tags inherited from earlier upstream or fork activity.

## Existing historical tags
- 2021-10-19
- 2021-10-29
- 2024-01-21
- 2024-02-08
- windows

## Decision for v0.1.1
- Keep the historical tags for now.
- Do not delete or rename them during v0.1.1.
- Treat v0.1.0 as the first public KitsuSync release gate.
- Treat the older non-v tags as historical upstream/fork artifacts, not as official KitsuSync releases.
- Use semantic version tags such as v0.1.1, v0.2.0, and later for future KitsuSync releases.

## Rationale
- Deleting historical tags is not required for v0.1.1 hardening.
- Leaving them in place avoids unnecessary history cleanup during a stabilization cycle.
- A semantic version line makes the public release story clearer without rewriting inherited history.
- Public-facing docs and release messaging should point to v0.1.0+ as the KitsuSync release line.

## Future cleanup conditions
Consider historical tag cleanup only after a separate decision or issue that confirms:
- which tags are safe to remove
- whether any automation, links, or users still rely on them
- how the cleanup will be communicated in public repo history

## What not to do now
- Do not delete old tags in v0.1.1.
- Do not rewrite git history.
- Do not present the old non-v tags as KitsuSync release milestones.
