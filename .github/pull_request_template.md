## Summary

- What changed?
- Why was it needed?

## Verification

- [ ] `docker compose build app` succeeded
- [ ] `docker compose ps` shows `app` healthy
- [ ] polling still works (`Connected to Kitsu`, `Got tasks`, `Got taskStatuses`)
- [ ] relevant routes were checked
- [ ] docs or onboarding text was updated if behavior changed

Routes checked:

- [ ] `/`
- [ ] `/bot/login`
- [ ] `/bot/docs/`
- [ ] `/bot/admin`
- [ ] `/bot/setup`

## Runtime Impact

- [ ] no runtime behavior change
- [ ] setup flow touched
- [ ] auth or cookie handling touched
- [ ] notification routing touched
- [ ] SQLite or persistence touched

## Security Review Notes

- [ ] this change does not expose secrets
- [ ] reverse proxy assumptions remain valid
- [ ] FileBrowser/debug profile safety remains valid

## Extra Notes

- rollout concerns:
- follow-up checks:
