# CLAUDE.md — Stellar Disbursement Platform Backend

## HTGC Account Setup & Minting Flow

### Key accounts

| Account | Public Key | Role |
|---|---|---|
| **Issuer** | `GDSRYZWTLQLBECKCL4TV7ZRGBZGBMSPD4V47B7Y7JSQVDJRSEXQTFCQT` | Mints HTGC, authorizes trustlines. Keep cold. |
| **Distributor** | `GCP6VMZS3SJ4CSOT3ZVMMJIOXOHTMJK47YQ4RTUJN7P2KYKDVRCUBS2X` | Holds pre-minted float, sends day-to-day payments. Used by SDP/TSS. |

The issuer has `AUTH_REVOCABLE` and `CLAWBACK_ENABLED` flags set. `AUTH_REQUIRED` was removed — trustlines auto-authorize, so payments flow immediately after the recipient creates a trustline. The issuer retains freeze and clawback rights.

### One-time setup: fund the distributor

1. **Distributor creates trustline** — sign with distributor key:
   - Operation: `Change Trust`
   - Asset: `HTGC` / issuer `GDSRYZWTLQLBECKCL4TV7ZRGBZGBMSPD4V47B7Y7JSQVDJRSEXQTFCQT`
2. **Issuer authorizes distributor trustline** — sign with issuer key:
   - Operation: `Set Trust Line Flags`
   - Trustor: `GCP6V...UBS2X`, Asset: `HTGC`, Set flags: `Authorized`
3. **Issuer mints to distributor** — sign with issuer key:
   - Operation: `Payment`
   - Destination: `GCP6V...UBS2X`, Asset: `HTGC`, Amount: large float (e.g. 10,000,000)

Repeat step 3 whenever the distributor float needs topping up.

### Per-recipient setup (every new user)

These steps must happen in order before a user can receive HTGC:

**Step 1 — Recipient creates trustline** (signed by recipient)
- Operation: `Change Trust`
- Asset: `HTGC` / issuer `GDSRYZWTLQLBECKCL4TV7ZRGBZGBMSPD4V47B7Y7JSQVDJRSEXQTFCQT`
- The recipient needs XLM for the base reserve (0.5 XLM per trustline)

**Step 2 — Distributor sends payment** (signed by distributor key)
- Operation: `Payment`
- Source: `GCP6V...UBS2X`
- Destination: `<recipient address>`
- Asset: `HTGC` / `GDSRY...FCQT`
- Amount: disbursement amount

### Correct flow diagram

```
Issuer mints → Distributor        (one-time or periodic top-up)
Issuer authorizes → Recipient     (per new user, before first payment)
Distributor pays → Recipient      (every disbursement)
```

### Common errors

| Error | Cause | Fix |
|---|---|---|
| `op_malformed` on change_trust | Issuer account not funded on testnet | Friendbot both accounts |
| `op_not_authorized` on payment | Recipient trustline not authorized by issuer | Run Step 2 first |
| `op_no_trust` on payment | Recipient has no trustline at all | Run Step 1 first |
| `tx_bad_auth` | Signed with wrong secret key | Check which account is `source_account`, sign with its key |
| `tx_failed` + 0 balance despite sent | Trustline is for a different HTGC issuer | Recipient must trust the correct issuer address |

---

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

---

## Theo — Business Model

This SDP fork is the disbursement backbone of **Theo**, a stablebond-backed payment and savings protocol for the US–Haiti corridor.

### What Theo is

Theo lets Haitian importers pay foreign suppliers instantly via an Odoo plugin ("Pay with Theo"), and lets international NGOs disburse to Haitian beneficiaries via the SDP dashboard. Both use **HTG-C**, a Haitian Gourde-denominated token backed entirely by a USD reserve held offshore.

### The two tokens

| Token | Peg | Purpose |
|---|---|---|
| **HTG-C** | 1 HTG-C = 1 HTG (display); backed by USD reserve | Local spending, NGO disbursements, importer settlement |
| **THEO-USD** | 1:1 USD | Transmission token and dollar savings account (later phase) |

### Legal structure

- **Theo managed service** — operates the SDP (this repo), the Odoo plugin, and customer relationships.
- **BVI SPV** — separate legal entity that issues HTG-C on Stellar and holds the reserve. This separation keeps the token issuance outside Haiti's BRH (central bank) controls.

### The dollar sourcing problem and how it's solved

The BRH tightly rations USD inside Haiti — Haitian banks cannot sell dollars freely at scale. Theo solves this by sourcing dollars **entirely outside Haiti**:

- **International NGOs** wire USD to the BVI SPV (outside Haiti) to fund disbursements → SPV mints HTG-C → NGOs disburse via SDP to beneficiaries
- **Haitian importers** pre-fund in HTG via local agents → credited with HTG-C → use Odoo plugin to pay foreign suppliers in USD drawn from the reserve

NGOs and importers are opposite sides of a natural dollar netting book. The reserve never needs to touch the Haitian banking system.

### Reserve composition (target)

```
50%  USTRY  — US Treasury bills, ~3.2% APY (Etherfuse, native Stellar)
40%  CETES  — Mexican Treasury bills, ~9% APY (Etherfuse, native Stellar)
10%  USDC   — Instant redemption buffer (Circle, native Stellar)

Over-collateralization: 105% enforced at every state change
```

The 105% OC buffer covers HTG appreciation of up to ~4.8% — larger than Haiti's historical annual max.

### Revenue streams

| Stream | Mechanism |
|---|---|
| FX fees | 1.5% on every HTG-C conversion (in and out) |
| Bond yield | USTRY (~3.2%) + CETES (~9%) on the float |
| HTG depreciation gains | When HTG weakens, HTG-C liability falls in USD terms → surplus accrues |
| AMM fees | CETES/USDC pool on Stellar DEX earns 30 bps per trade (later phase) |

### Go-to-market sequencing

1. **Phase 1 — Enterprise importers** (Odoo plugin, `THEO_ODOO_PLUGIN_SPEC.md`): Haitian importers pre-fund in HTG, pay foreign suppliers at 1.5% vs 4–7% bank rates. Anchor customer: NABATCO.
2. **Phase 2 — International NGOs** (SDP dashboard): NGOs bring USD in, disburse HTG-C to beneficiaries. Dollar inflows fund the reserve.
3. **Phase 3 — Consumer remittance** (deferred): Diaspora USD → HTG-C to family. Slots into an already-proven netting book.

### How this repo fits

The SDP (`sdp-api` + `sdp-tss`) is the **disbursement layer** — it handles bulk sends of HTG-C to recipient wallets (NGO beneficiaries, importer payment confirmations). The BVI SPV and Soroban vault contract handle minting/burning and reserve management separately.

### Key specs

- `THEO_SPEC.md` — Soroban vault contract spec (mint/burn, 105% OC invariant, accounting model)
- `THEO_ODOO_PLUGIN_SPEC.md` — Odoo 17 plugin spec ("Pay with Theo" button, FX wizard, mock API)
- `THEO_PRODUCT_SUMMARY.md` — full product brief with unit economics
- `THEO_ODOO_PR_FAQ.md` — press release and stakeholder FAQ
