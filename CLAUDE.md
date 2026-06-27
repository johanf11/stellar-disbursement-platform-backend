# CLAUDE.md — Stellar Disbursement Platform Backend

## What this repo is

Fork of [stellar/stellar-disbursement-platform-backend](https://github.com/stellar/stellar-disbursement-platform-backend) (v6.6.1).
Go monolith that handles bulk Stellar payments: API server, Transaction Submission Service (TSS), multi-tenant DB migrations, SEP-10/24/45 auth.

## Deployment (production)

| Thing | Value |
|---|---|
| Droplet IP | `64.225.19.99` |
| SSH key | `~/.ssh/id_do_sdp` |
| Binary on server | `/opt/sdp/stellar-disbursement-platform` |
| Env file on server | `/opt/sdp/.env` |
| API | `http://64.225.19.99:8000` |
| Admin API | `http://64.225.19.99:8003` |
| Frontend | `http://64.225.19.99:3000` (nginx → `/opt/sdp-frontend`) |
| Database | Supabase `pwbjulfhqryrinkwijrx` (PostgreSQL 17, IPv6 direct) |
| Network | Stellar testnet |

**Deploy backend:** `./dev/deploy.sh` — cross-compiles `linux/amd64`, SCPs binary + systemd units, skips `.env` if already present.

**Deploy frontend:** `./dev/deploy-frontend.sh` — builds the frontend repo (expected at `../stellar-disbursement-platform-frontend`), SCPs `build/` to droplet, reloads nginx.

## Systemd services

Both are enabled and start on boot.

```
sdp-api   — runs 5 DB migrations then `serve` on :8000/:8002/:8003
sdp-tss   — runs after sdp-api, manages channel accounts + submits txs
```

Check status / logs:
```bash
ssh -i ~/.ssh/id_do_sdp root@64.225.19.99
systemctl status sdp-api sdp-tss --no-pager -l
journalctl -u sdp-api -f
```

## Database

- **Provider:** Supabase (project `pwbjulfhqryrinkwijrx`)
- **Connection:** direct IPv6 — `db.pwbjulfhqryrinkwijrx.supabase.co:5432` (pooler doesn't work with this setup)
- **Schemas:** `admin`, `tss`, `auth` — NOT `public` (public schema is empty by design)
- **`uuid-ossp`** lives in the `extensions` schema; a wrapper `public.uuid_generate_v4()` was created manually

Local dev uses the same Supabase DB. Set `DATABASE_URL` in `dev/.env.default` (gitignored).

## Build & run locally

```bash
# Build
go build -o stellar-disbursement-platform .

# Run (uses dev/.env.default via docker-compose or export manually)
./stellar-disbursement-platform serve

# Cross-compile for the droplet
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stellar-disbursement-platform .
```

## Common make targets

```bash
make go-build    # compile check
make go-test     # unit tests (requires gotestsum)
make go-lint     # golangci-lint
make setup       # interactive setup wizard
```

## Local dev with Docker Compose

Docker Compose (`dev/docker-compose.yml`) runs `sdp-api`, `sdp-tss`, and `sdp-frontend` — all connect to Supabase (no local DB container).

```bash
cp dev/.env.default dev/.env   # then fill in secrets
docker compose -f dev/docker-compose.yml up
```

**Note:** Docker is disabled on the production droplet to free RAM. Do not re-enable it there.

## Environment files

| File | Purpose |
|---|---|
| `dev/.env.default` | Local dev defaults (gitignored) |
| `dev/.env.production` | Production values, uploaded once by deploy.sh (gitignored) |
| `dev/env.production.example` | Template committed to git — no secrets |

Key env vars to know:
- `DATABASE_URL` — must be the direct Supabase URL (not pooler)
- `INSTANCE_NAME` — no spaces (systemd EnvironmentFile limitation)
- `SINGLE_TENANT_MODE=true` — one tenant, owner email is `DEFAULT_TENANT_OWNER_EMAIL`
- `EMAIL_SENDER_TYPE=DRY_RUN` — emails logged, not sent; check `journalctl -u sdp-api` for reset links
- `DISABLE_MFA=true` — MFA off for demo
- `DISABLE_RECAPTCHA=true` — reCAPTCHA off for demo

## Known gotchas

- **Port 8000 conflict on restart:** Docker proxy used to hold port 8000. Docker is now disabled (`systemctl disable docker`). If port conflicts recur: `ss -tlnp | grep :8000` then `pkill docker-proxy`.
- **`uuid_generate_v4` missing:** Supabase puts `uuid-ossp` in the `extensions` schema. A wrapper function exists in `public`. If migrations fail with this error, check the wrapper is still there via Supabase SQL editor.
- **Supabase IPv6:** The direct connection host resolves IPv6-only. The droplet has IPv6 enabled (`2604:a880:800:14:0:3:2b75:4000/64`). The pooler (`aws-0-us-west-2.pooler.supabase.com`) does not work.
- **systemd EnvironmentFile:** Quoted values with spaces break parsing. Use `KEY=value` without quotes for values containing spaces (e.g. `INSTANCE_NAME=SDP-Demo`).
- **TSS channel accounts:** `channel-accounts ensure 3` (matches `NUM_CHANNEL_ACCOUNTS=3`). Using `ensure 1` causes TSS to exit immediately.

## Frontend repo

Companion repo: `johanf11/stellar-disbursement-platform-frontend`
Expected locally at: `../stellar-disbursement-platform-frontend` (relative to this repo)
Runtime config (gitignored on frontend): `public/settings/env-config.js` → on server at `/opt/sdp-frontend/settings/env-config.js`
