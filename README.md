# Borehole Edge-Scoring Engine

High-performance Go-based mobile infrastructure for offline credit scoring in Kenya. Parses M-Pesa, Airtel Money, T-Kash, Hustler Fund, and bank SMS messages to generate XGBoost-compatible feature vectors.

## Features

- **Automated SMS Scraping (New)**
  - One-tap financial health scan directly from Android inbox.
  - Smart keyword filtering (Confirmed, M-PESA, Airtel, HustlerFund, etc.).
- **Multi-Provider SMS Parsing**
  - **Mobile Money**: M-Pesa (Supports any alphanumeric series), Airtel Money, T-Kash
  - **Credit Products**: Fuliza, Hustler Fund (Advanced repayment detection), Okoa Jahazi (Debt snapshotting)
  - **Digital Lenders**: Tala, Branch, Zenka, Zash, Okolea
  - **Savings & Banking**: M-Shwari, KCB M-Pesa, Mali, Stawi, Bank Transfers (Equity, KCB, Co-op, NCBA)
- **22-Feature Vector** - Financial health, lender diversity, emergency reliance, savings rate
- **Local-Only Processing** - 100% offline, privacy-first engine. No SMS logs ever leave the device.
- **Go 1.25+ Infrastructure** - Zero-allocation techniques and high-performance routing.


## Quick Start

```bash
# Build
go build ./...

# Run API server
go run cmd/api/main.go

# Test scoring endpoint
curl -X POST http://localhost:8080/v1/score \
  -H "Content-Type: application/json" \
  -d '{"logs": ["UA1234ABCD Confirmed. You have received Ksh1,500.00 from JOHN DOE", "Hustler Fund. You have been disbursed Ksh500.00"]}'
```

## Project Structure

```
borehole-engine/
├── cmd/api/          # API entrypoint
│   └── main.go       # HTTP server with POST /v1/score
├── pkg/
│   ├── parser/       # SMS parsing
│   │   ├── parser.go     # Transaction types, ParseLogs()
│   │   └── patterns.go   # Pre-compiled regex patterns (M-Pesa, Airtel, Hustler, etc.)
│   └── engine/       # Feature engineering
│       └── engine.go     # 22-feature Vectorize()
└── go.mod
```

## API

### `POST /v1/score`

**Request:**
```json
{
  "logs": [
    "UA1234ABCD Confirmed. You have received Ksh1,500.00 from JOHN DOE",
    "Fuliza M-PESA. You have borrowed Ksh2,000.00",
    "Transaction ID: AM1234. You have received Ksh1,000 from Jane",
    "You have received Ksh5,000 from Tala"
  ]
}
```

**Response:**
```json
{
  "score": 0.72,
  "features": [2500, 0, 2500, ...],
  "txn_count": 4
}
```

## Feature Vector

| Index | Feature | Description |
|-------|---------|-------------|
| 0-2 | Income/Expense/Net | Financial flow metrics |
| 6 | gambling_index | Betting spend ratio |
| 8 | fuliza_usage | Fuliza borrowed / income |
| 15 | hustler_balance | Latest Hustler Fund debt |
| 16 | okoa_frequency | Okoa Jahazi usage count |
| 17 | airtel_volume | Total Airtel Money volume |
| 18 | lender_diversity | Count of unique digital lenders |
| 19 | emergency_reliance | (Okoa + Fuliza) / Income |
| 20 | savings_rate | MMF deposits / Income |
| 21 | bank_activity | Count of bank transactions |

## License

MIT
