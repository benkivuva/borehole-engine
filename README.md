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

## Mobile Application (React Native)

The repository includes a complete React Native application in the `MobileApp/` directory that demonstrates the engine in action.

### Key Features
*   **Encrypted Storage**: Uses **SQLCipher** (AES-256) to store credit score history.
*   **Android Keystore**: Encryption keys are generated and stored in the hardware-backed Android KeyStore.
*   **Privacy First**: "Nuke Data" feature securely wipes all local data via `VACUUM`.
*   **Real-time Dashboard**: Visualizes the 20-dimension feature vector (Cash In vs Out, Debt Ratios).

### Quick Start (Mobile)

Prerequisites: Node.js, Android Studio, Java 17+.

```bash
# 1. Install Dependencies
cd MobileApp
npm install

# 2. Build Go-Mobile Library (If changing engine code)
cd ..
gomobile bind -v -target=android -androidapi 21 -o MobileApp/android/app/libs/borehole.aar ./pkg/mobile

# 3. Run on Android Emulator
cd MobileApp
npm run android
```

## Project Structure

```
borehole-engine/
├── cmd/api/          # API entrypoint
│   └── main.go       # HTTP server with POST /v1/score
├── pkg/
│   ├── parser/       # SMS parsing (Regex patterns)
│   ├── engine/       # Feature engineering & ML Logic
│   └── mobile/       # Go-Mobile bindings (JNI Bridge)
└── MobileApp/        # React Native Application
    ├── android/      # Native Android project (Gradle)
    ├── src/
    │   └── storage/  # Encrypted Database (SQLCipher)
    ├── App.tsx       # Main UI & Logic
    └── BoreholeBridge.ts # JS Bridge to Go Engine
```

## API

### `POST /v1/score`

**Request:**
```json
{
  "logs": [
    "UA1234ABCD Confirmed. You have received Ksh1,500.00 from JOHN DOE",
    "Fuliza M-PESA. You have borrowed Ksh2,000.00"
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

## Feature Vector (20 Dimensions)
| Index | Feature | Description |
|-------|---------|-------------|
| 0-5 | Financial Health | Income, Expenses, Profitability, Txn Count, Max Txn, Consistency |
| 6 | `gambling_index` | Betting spend / Total Expenses |
| 7 | `utility_ratio` | Utility spend / Total Expenses |
| 8-9 | `fuliza_stats` | Usage Ratio & Repayment Rate |
| 13 | `hustler_balance` | Latest Hustler Fund debt |
| 14 | `okoa_frequency` | Okoa Jahazi usage count |
| 15 | `airtel_volume` | Total Airtel Money volume |
| 16 | `lender_diversity` | Count of unique digital lenders |
| 17 | `emergency_reliance` | (Okoa + Fuliza) / Income |
| 18 | `savings_rate` | MMF deposits / Income |
| 19 | `bank_activity` | Count of bank transactions |

https://github.com/benkivuva/borehole-engine

## License

MIT
