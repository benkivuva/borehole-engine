# Borehole Edge-Scoring Engine

High-performance Go-based mobile infrastructure for offline credit scoring in Kenya. Parses M-Pesa, T-Kash, and Fuliza SMS messages to generate XGBoost-compatible feature vectors.

## Features

- **SMS Parsing** - M-Pesa (2026 UA series), T-Kash, Fuliza with named capture regex
- **15-Feature Vector** - Gambling index, utility ratio, income regularity, Fuliza usage
- **Mobile-Optimized** - Zero-allocation techniques, minimal GC pressure
- **Go 1.22+ API** - Modern ServeMux with method routing

## Quick Start

```bash
# Build
go build ./...

# Run API server
go run cmd/api/main.go

# Test scoring endpoint
curl -X POST http://localhost:8080/v1/score \
  -H "Content-Type: application/json" \
  -d '{"logs": ["UA1234ABCD Confirmed. You have received Ksh1,500.00 from JOHN DOE"]}'
```

## Project Structure

```
borehole-engine/
├── cmd/api/          # API entrypoint
│   └── main.go       # HTTP server with POST /v1/score
├── pkg/
│   ├── parser/       # SMS parsing
│   │   ├── parser.go     # Transaction types, ParseLogs()
│   │   └── patterns.go   # Pre-compiled regex patterns
│   └── engine/       # Feature engineering
│       └── engine.go     # 15-feature Vectorize()
└── go.mod
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
  "score": 0.659,
  "features": [1500, 0, 1500, ...],
  "txn_count": 2
}
```

### `GET /health`

Returns server health status.

## Feature Vector

| Index | Feature | Description |
|-------|---------|-------------|
| 0 | total_income | Sum of received amounts |
| 1 | total_expenses | Sum of sent/paid amounts |
| 2 | net_flow | Income - expenses |
| 6 | gambling_index | Betting spend ratio |
| 7 | utility_ratio | Utility payments ratio |
| 8 | fuliza_usage | Fuliza borrowed / income |

## License

MIT
