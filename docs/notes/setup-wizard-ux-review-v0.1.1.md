# Setup Wizard UX Review for v0.1.1

## Executive summary
The current setup surfaces are functional, but the first-time operator experience still has a trust gap. The biggest problems are not deep behavior bugs; they are inconsistencies between the public docs, the wizard entry flow, and the older /bot/setup surface. A beginner can reach the right outcome, but they are asked to reconcile multiple setup paths and multiple meanings of  Bot Setup on their own.

The safest next step is not a flow rewrite. The safest next step is to align the public docs and on-screen copy around one first-time path, then tighten the wording around what is only a test, what is saved permanently, and when rollback is guaranteed.

## High-confidence UX issues

### 1. First-time setup path is inconsistent across public docs and UI
- Risk: high
- Evidence:
  - README.md still presents /bot/setup as a core first-time setup surface in multiple places.
  - docs/QUICK_START.md sends the user to /bot/admin/bot and /bot/admin/projects before the wizard.
  - docs/SETUP_WIZARD.md describes /bot/admin/setup as the intended diagnostics-first entry.
  - docs/SETUP_FOR_STUDIOS.md labels Step 2 as /bot/admin/bot but then tells the user to open /bot/setup and click Bot Setup.
  - src/setup/wizard.go shows that /bot/setup-wizard first opens an entry chooser, not the old direct 4-step flow.
- Why this matters:
  - A beginner cannot tell which page is the real first stop.
  - Trust drops when the docs and UI appear to disagree before any setup work is done.
- Recommended fix type: docs-only first, then small copy/text change

### 2. The public description still frames the wizard like a direct 4-step flow, but the actual product is now entry -> guided/quick/manual
- Risk: high
- Evidence:
  - README.md says /bot/setup-wizard is a guided 4-step first-time setup.
  - docs/QUICK_START.md says The wizard walks you through 4 steps.
  - At the time of this review, `src/setup/quick_setup.go` rendered an entry page with two modes: Guided Setup and Quick Setup.
  - src/setup/guided_stepwise.go can start from System Check as step 0 when env prerequisites are missing.
- Why this matters:
  - A first-time operator expects one linear wizard, but the actual UX starts with a mode choice and can insert a prerequisite gate before Kitsu.
  - This makes the first screen feel wrong even when the product is behaving correctly.
- Recommended fix type: docs-only first, then small copy/text change

### 3. Rollback and retry guarantees are described too absolutely in some places and more cautiously in others
- Risk: high
- Evidence:
  - docs/SETUP_WIZARD.md says the entire setup is rolled back automatically if channel creation fails.
  - docs/SETUP_FOR_STUDIOS.md also says partial changes are rolled back automatically.
  - docs/TROUBLESHOOTING.md later says rollback is best-effort and manual cleanup may still be required.
  - src/setup/handler.go and src/setup/setupapi.go show both cases exist:
    - normal failure paths that roll back and can become Safe to retry
    - database-transaction failure after Discord creation, where manual cleanup is explicitly required
- Why this matters:
  - This is the trust-critical part of the setup flow. If the user is told automatic rollback and then has to clean up Discord manually, they will stop trusting the tool.
- Recommended fix type: docs-only first, behavior change only if a separate issue is opened later

### 4. The product does not clearly separate connection testing from persistent configuration at the public onboarding level
- Risk: medium
- Evidence:
  - src/setup/guided_stepwise.go explains in Step 2 that the form is for connection testing only, and permanent save belongs in /bot/admin/bot.
  - That distinction is much less visible in README.md and docs/QUICK_START.md.
  - docs/SETUP_FOR_STUDIOS.md mixes Bot Setup terminology with both admin credential storage and /bot/setup actions.
- Why this matters:
  - Beginners want to know whether pressing a button will only test a credential, or actually save and change runtime state.
  - Hidden persistence boundaries make bot token handling feel riskier than it needs to.
- Recommended fix type: docs-only first, then small copy/text change

### 5. Discord terminology is technically accurate but not always operator-friendly at the moment of action
- Risk: medium
- Evidence:
  - Public docs mix Guild ID, Discord Server, and server membership explanations across README.md, docs/QUICK_START.md, and docs/SETUP_WIZARD.md.
  - src/setup/guided_stepwise.go tries to prefer Discord Server in Step 3, but Step 2 still centers Discord Server ID (Guild ID).
  - src/setup/admin.go still exposes Discord Guild ID directly in project assignment.
- Why this matters:
  - Beginners usually understand Discord server faster than guild.
  - The current wording is not wrong, but it still asks the operator to translate platform terminology mentally.
- Recommended fix type: copy/text change

### 6. Recovery guidance exists, but it is spread across too many places to feel safe during failure
- Risk: medium
- Evidence:
  - docs/TROUBLESHOOTING.md has the right failure guidance.
  - src/setup/guided_stepwise.go can show Safe to retry and raw details.
  - src/setup/quick_setup.go and /bot/admin/setup present diagnostics separately.
  - The public onboarding path does not strongly tell the user where to go next when setup partially works but verification fails.
- Why this matters:
  - A first-time operator needs one obvious if this fails go here next answer.
  - Right now that answer exists, but it is distributed between wizard, manual setup, and troubleshooting docs.
- Recommended fix type: docs-only first, then UI layout change if still confusing after copy cleanup

## What not to change yet
- Do not rewrite setup provisioning behavior.
- Do not merge /bot/setup-wizard, /bot/setup, and /bot/admin/* into a new architecture during v0.1.1.
- Do not change auth/session flow as part of UX cleanup.
- Do not change Discord/Kitsu provisioning logic just to make the UI narrative simpler.
- Do not add more automation until the current manual / semi-automatic / automatic boundaries are explained more clearly.

## Suggested next PRs in safe order
1. Docs-only alignment PR
   - Pick one first-time operator path and state it consistently in README.md, docs/QUICK_START.md, docs/SETUP_WIZARD.md, docs/SETUP_FOR_STUDIOS.md, and docs/TROUBLESHOOTING.md.
2. Copy-only PR for setup surfaces
   - Clarify test only vs saved permanently for Kitsu runtime credentials and Discord bot token handling.
3. Copy-only PR for rollback language
   - Explicitly separate automatic rollback complete from manual cleanup required cases.
4. Small layout PR if still needed
   - Add one stronger handoff from Guided Setup failure states to Manual Setup / Troubleshooting.
5. Separate behavior issue, not PR-by-default
   - Only if desired later, open a separate issue for any real setup behavior simplification.

## Review scope used for this pass
- Routes reviewed: /bot/setup-wizard, /bot/setup, /bot/admin/setup, /bot/admin/bot, /bot/admin/projects
- Public docs reviewed: README.md, docs/QUICK_START.md, docs/SETUP_WIZARD.md, docs/SETUP_FOR_STUDIOS.md, docs/TROUBLESHOOTING.md
- Relevant source reviewed for UX behavior and error states: src/setup/wizard.go, src/setup/quick_setup.go, src/setup/guided_stepwise.go, src/setup/setupapi.go, src/setup/handler.go, src/setup/admin.go, src/setup/ux_setup.go
