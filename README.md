# Borehole Edge-Scoring Engine

High-performance Go-based mobile infrastructure for offline credit scoring in Kenya. Parses M-Pesa, Airtel Money, T-Kash, Hustler Fund, and bank SMS messages to generate XGBoost-compatible feature vectors.

## 🧠 Edge ML Inference (New!)

Borehole uses a **Probabilistic XGBoost Inference Engine** running directly on the mobile device (Go-based).

### The Pipeline: From SMS to Score

1.  **Extract (Parser)**: Zero-allocation Regex engine extracts amounts, dates, and types (e.g., M-Pesa, Hustler Fund).
2.  **Transform (Mapper)**: Transactions are aggregated into a fixed **20-dimensional feature vector** (e.g., *Income*, *GamblingRatio*, *LoanRepaymentRate*).
3.  **Inference (XGBoost)**: The vector is passed to an embedded Gradient Boosting model (via `dmitryikh/leaves`).
4.  **Activation (Probability)**: Raw tree margins are squashed via **Sigmoid** to a 0.0-1.0 risk probability.

### Core Features
*   **Performance**: Zero-allocation inference loop optimized for mobile ARM processors.
*   **Safety**: Robust fallback mechanism ensures the app never crashes even if the model file is corrupted (defaults to neutral score).
*   **Privacy**: 100% Offline. No data leaves the device.


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
| 0-5 | Financial Health | Income, Expenses, Profitability, Txn Count, Max Txn, Consistency |

| 6 | `gambling_index` | Betting spend / Total Expenses |
| 7 | `utility_ratio` | Utility spend / Total Expenses |
| 8 | `fuliza_usage` | Fuliza borrowed / Income |
| 9 | `fuliza_repay` | Fuliza repayment rate |
| 13 | `hustler_balance` | Latest Hustler Fund debt |
| 14 | `okoa_frequency` | Okoa Jahazi usage count |
| 15 | `airtel_volume` | Total Airtel Money volume |
| 16 | `lender_diversity` | Count of unique digital lenders |
| 17 | `emergency_reliance` | (Okoa + Fuliza) / Income |
| 18 | `savings_rate` | MMF deposits / Income |
| 19 | `bank_activity` | Count of bank transactions |


## License

MIT
