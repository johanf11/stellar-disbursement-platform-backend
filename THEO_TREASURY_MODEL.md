# THEO_TREASURY_MODEL.md — Treasury, Float & Matched-Book Mechanics

Companion to `THEO_SPEC.md` (vault contract), `THEO_ODOO_PLUGIN_SPEC.md` (importer UX), and `THEO_PRODUCT_SUMMARY.md` (business model). This doc covers **how the money actually moves**: the netting book, float sizing, inventory math, yield models, and a concrete plan to launch and validate the accounting with **$100–200K of seed capital**.

---

## 1. Purpose

Validate, with real but small capital, that the matched-book model clears correctly:

- NGOs bring **USD wanting HTG** (beneficiary disbursements)
- Importers bring **HTG wanting USD** (supplier payments)
- **HTG-C** is the clearing token between them
- Theo captures a **~2–3% spread**, and the accounting balances across all three entities at every step

The goal of the first deployment is **not** revenue. It is a clean, auditable ledger that proves the mechanics, so the raise (VC or crypto) is backed by a working system and real reconciliation, not a deck.

---

## 2. Entity Structure

Three books, three currencies, three jurisdictions. Keep them separate from day one.

| Entity | Holds | Currency | Role |
|---|---|---|---|
| **Theo BVI SPV** | USTRY/CETES/USDC reserve | USD | Issuer — mints/burns HTG-C, owns the reserve, books the spread |
| **LP Pool** (Stellar smart contract) | USDC ↔ HTG-C | USD/HTG-C | Swap liquidity for importer → supplier payments |
| **Theo Haiti Ops** | HTG bank float | HTG | Receives importer HTG, pays beneficiary cash, BRH compliance |

> **Invariant:** HTG-C is only ever minted against USD → USTRY (or USD → USDC in the test phase). Never against HTG. Every HTG-C in circulation is dollar-backed by construction. This is what makes over-collateralization self-enforcing.

---

## 3. The Matched Book

```
        USD IN                                  HTG OUT
   ┌──────────────┐                        ┌──────────────┐
   │     NGO      │   wires $10K       →    │ Beneficiaries│  paid HTG cash
   └──────┬───────┘                        └──────▲───────┘
          │              ┌─────────────┐          │
          └─────────────►│   HTG-C     │──────────┘
                         │  (clearing) │
          ┌─────────────►│             │──────────┐
          │              └─────────────┘          │
   ┌──────┴───────┐                        ┌──────▼───────┐
   │   Importer   │   deposits 1.3M HTG →  │   Supplier   │  paid USDC
   └──────────────┘                        └──────────────┘
        HTG IN                                  USD OUT
```

**In a matched cycle the NGO's dollars pay the importer's supplier; the importer's gourdes pay the NGO's beneficiary.** The Haiti bank is a net-zero HTG pass-through. Theo's spread is captured as USD in the SPV reserve. Nothing crosses the BRH wall.

### One matched $10K cycle — net positions

| Entity | Net change | Where |
|---|---|---|
| SPV | **+~$300** (spread, as reserve surplus) | Offshore, USD |
| Pool | $0 + fee share | On-chain |
| Haiti Ops | **$0** HTG (in = out) + Bons BRH on transit | Onshore, HTG |

---

## 4. The NGO Payroll-Float Mechanic (your free working capital)

NGOs run **1–2 week payroll cycles**. They pre-fund the full cycle up front, then disburse in tranches. That parked USD is the engine of the model:

```
Day 0:   NGO wires $10K for the biweekly payroll
         → sits in SPV reserve / pool as USDC
         → IS the swap liquidity for importer payments this cycle
Days 1-14: disbursed to beneficiaries in tranches as HTG cash
```

**The parked NGO float doubles as importer swap liquidity for free, for up to 2 weeks.** Importer payments up to the parked amount need no Theo capital — the NGO's dollars-in-transit fund the supplier payout, and the importer's HTG refills the bank float the beneficiaries drew down.

Your seed capital only has to cover:
1. **HTG-C inventory** — the bridge tokens importers draw
2. **HTG bank pre-stock** — gourdes to pay beneficiaries before the matching importer deposit lands
3. **Timing/size gaps** — when importer demand exceeds the parked NGO float in a given window

---

## 5. Float & Inventory Sizing

### Constants

```
FX rate:        130 HTG / USD       (set per actual market at launch)
OC ratio:       105%
Spread:         ~3% round trip      (e.g. 1% NGO in + 2% importer out — configurable)
Cycle length:   14 days (NGO payroll float)
```

### Per-payment capital — the direct answer

At the instant a payment settles, two assets must be present simultaneously: **HTG-C inventory** (to hand the importer) and **USDC liquidity** (to pay the supplier). The HTG bank must also be able to cover the beneficiary side.

| Component | $10K payment | $50K payment |
|---|---|---|
| HTG-C inventory | $10,000 | $50,000 |
| Reserve backing (105%) | $10,500 | $52,500 |
| USDC swap liquidity | $10,000 | $50,000 |
| HTG bank float (beneficiary) | 1.3M HTG (~$10K) | 6.5M HTG (~$50K) |

**Two funding scenarios:**

```
MATCHED (NGO float covers the USDC + importer HTG covers the bank):
   Your capital = reserve backing only
   → ~$10.5K  for a $10K payment
   → ~$52.5K  for a $50K payment

SELF-FUNDED (you bridge the full gap — no NGO overlap yet):
   Your capital = inventory backing + USDC + HTG float
   → ~$30K   for a $10K payment
   → ~$150K  for a $50K payment
```

### Concurrency

Multiply by the number of payments in flight within one 14-day cycle:

```
Inventory needed = (peak concurrent in-flight payments) × (per-payment requirement)
```

For testing you control the timing, so run payments **sequentially** and reuse the same inventory each cycle. One $50K of working inventory, turned over, validates the same mechanics as ten.

### What $150K seed buys (test phase, USDC reserve)

```
$150,000 seed
  − $30,000  HTG bank float + gas + ops buffer
  = $120,000 working inventory + swap liquidity

Matched (NGO funds the USDC side):
   $120K all → inventory → supports up to $120K concurrent in-flight,
   bounded by available NGO parked float

Self-funded (~2× per payment):
   ~$60K concurrent capacity → at 14-day turnover ≈ $120K/month ≈ $1.4M/yr
```

For **validation**, this is comfortable headroom: extensive $10K cycles plus at least one matched $50K cycle on $100K, or an unmatched $50K on $200K.

---

## 6. Two Yield Models (decide before the raise)

Once real USTRY/CETES backs the reserve, the bond yield (~7%) can be kept or shared with AMM LPs.

| | **Model A — Theo keeps yield** | **Model B — Theo shares yield** |
|---|---|---|
| LP earns | Swap fees only (~7% APY, no floor) | 7% bond floor + fees (~14% APY) |
| LP capital | Flighty — flees in slow months | Sticky — stays through slow months |
| Theo keeps | Spread + full bond yield | Spread + bond surplus only |
| Pool liquidity when you need it | At risk | Reliable |

**The spread (~$3M/yr at $120M volume) dwarfs the bond yield (~$350K). The spread is always 100% Theo's.** The bond yield is a cheap loyalty payment.

**Recommendation:** Launch with **Model B** to attract and retain LP surge capacity while young; migrate toward **Model A** as retained spread grows your own float and LP dependence falls.

---

## 7. Seed Capital Deployment Plan ($100–200K)

### Simplify the test: USDC as reserve first

Do **not** onboard Etherfuse/USTRY for the first flows. Hold **USDC 1:1 as the reserve.** You forgo the ~7% bond yield (immaterial at test scale) and remove an integration dependency. The payment mechanics, fee capture, and reconciliation are identical. Add USTRY/CETES in Phase 2 once the cycle is proven.

### Allocation (example, $150K)

```
$100,000  USDC reserve  → backs HTG-C inventory + swap liquidity
$ 35,000  HTG bank float (≈ 4.5M HTG) → beneficiary redemptions
$ 10,000  USDC pool buffer + Stellar gas/fees
$  5,000  contingency
```

---

## 8. Phased Launch Roadmap

| Phase | Capital | Payment size | Reserve | Goal |
|---|---|---|---|---|
| **0 — Testnet** | $0 | test tokens | mock | Deploy vault + AMM on Stellar testnet; run full cycle; spreadsheet ties out to on-chain |
| **1 — Tiny mainnet** | $10–20K | $10K | USDC 1:1 | 1 friendly NGO + 1 friendly importer; validate mint → swap → redeem → fee → reconcile |
| **2 — Scaled test** | $100–200K | $50K | add USTRY/CETES | Multiple cycles; measure float turnover, fee accrual, OC stability; build the track record |
| **3 — Raise** | — | — | — | Pitch VC/crypto with a live system + audited reconciliation, not a deck |

**Gate between phases:** do not advance until the prior phase's ledger reconciles to the on-chain state with zero unexplained variance.

---

## 9. Accounting Ledger — What to Track

Maintain one reconciliation sheet per cycle, three columns (SPV / Pool / Haiti), tying to on-chain balances:

```
Per cycle, record:
  - NGO USD received            (SPV in)
  - HTG-C minted                (SPV liability up, reserve up)
  - HTG-C disbursed to NGO      (to SDP → beneficiaries)
  - HTG paid to beneficiaries   (Haiti bank out)
  - HTG deposited by importer   (Haiti bank in)
  - HTG-C drawn by importer     (inventory out)
  - HTG-C swapped in pool       (pool HTG-C in, USDC out)
  - USDC paid to supplier       (pool out)
  - HTG-C burned                (SPV liability down, reserve down)
  - Spread captured             (SPV surplus, USD)
  - Bons BRH accrued            (Haiti, HTG)

Reconciliation checks (must hold every cycle):
  ✓ Reserve ≥ 1.05 × HTG-C outstanding         (OC invariant)
  ✓ HTG-C minted − HTG-C burned = HTG-C outstanding
  ✓ Haiti bank: HTG in − HTG out = float delta  (≈ 0 when matched)
  ✓ SPV USD: NGO in − supplier out = spread     (the profit)
  ✓ On-chain balances = ledger balances          (zero variance)
```

---

## 10. Where Value Accumulates (read before relocating)

| Stream | Currency | Location |
|---|---|---|
| FX spread (~3%) | USD | **Offshore (SPV)** |
| Bond yield | USD | **Offshore (SPV)** |
| Depreciation carry | USD | **Offshore (SPV)** |
| Bons BRH on transit float | HTG | Haiti (small) |

**Financial value books offshore in USD — by design, and to your protection.** It sits beyond BRH reach, capital controls, and HTG depreciation. Warehousing value in Haiti as HTG nets **−3%/yr** (≈7% Bons BRH − ≈10% depreciation): keep the Haiti float lean and fast-cycling.

**Move to Haiti for operational control and economic substance** (importer relationships, bank ops, BRH compliance, market position) — not because the carry settles there. If you later want a Haiti profit center, build a **deliberate local HTG lending book** (rates 20–30%+ can beat depreciation); don't get there by accidentally leaving FX exposure on the balance sheet.

---

## 11. Path to Raise

By end of Phase 2 you can show investors:

1. **A live system** clearing real NGO ↔ importer payments on Stellar
2. **An audited reconciliation** — every cycle ties out, OC never broke
3. **Measured unit economics** — actual spread captured, actual float turnover
4. **A scaling model** — float and inventory sizing validated against real cycles
5. **A clear use of proceeds** — raise grows the reserve/inventory, which directly grows throughput (the model in §5)

The pitch becomes: *"Working corridor, proven accounting, $X anchor volume committed. Capital buys float; float is throughput; throughput is spread."*

---

## Appendix — Quick Reference

```
Per $10K payment:   $10.5K inventory backing (matched) | ~$30K (self-funded)
Per $50K payment:   $52.5K inventory backing (matched) | ~$150K (self-funded)
$150K seed:         ~$120K working capital → $1.4M/yr self-funded throughput
NGO float:          14-day parked USD = free importer swap liquidity
Test reserve:       USDC 1:1 (skip USTRY until Phase 2)
Yield model:        Start B (share), migrate to A (keep) as float grows
Value location:     Offshore USD (protect it); Haiti lean (operate it)
```
