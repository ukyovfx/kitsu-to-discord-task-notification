# Security Policy

## Supported Use

KitsuSync is intended to run:

- with secrets provided outside git
- with `/bot/*` protected by admin login
- behind a trusted reverse proxy in production

## Report a Vulnerability

Please do not open a public issue for sensitive problems.

Share the following privately with the maintainer:

- affected version or commit
- deployment style
  - direct `:8090`
  - reverse proxy
- reproduction steps
- impact summary
- relevant log snippets with secrets removed

## Scope to Review Carefully

- admin login and cookie handling
- runtime credential storage and migration
- Discord webhook routing
- setup rollback and partial failure handling
- FileBrowser debug profile exposure
- reverse proxy assumptions such as `X-Forwarded-Proto`

## Deployment Safety Notes

- Keep `.env`, `conf.toml`, and SQLite data out of public file browsers.
- Do not trust forwarded headers unless your reverse proxy overwrites them.
- Keep FileBrowser disabled unless you intentionally enable the debug profile.
