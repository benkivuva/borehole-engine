# Borehole Edge-Scoring Engine 🛡️📱

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![React Native](https://img.shields.io/badge/React_Native-0.73-61DAFB?style=flat&logo=react)](https://reactnative.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Privacy](https://img.shields.io/badge/Privacy-100%25_Offline-green)](README.md)

**Borehole** is a decentralized, privacy-first financial infrastructure that enables **offline credit scoring** for the unbanked in emerging markets. 

It parses unstructured financial SMS logs (M-Pesa, Airtel Money, Banks) directly on the user's device, generates a 20-dimensional risk vector, and calculates a credit score using an embedded **Go-based Inference Engine**.

Most importantly, it generates **Cryptographically Verifiable Claims** (Ed25519), allowing users to prove their creditworthiness to lenders without revealing their raw transaction history.

---

## 🌟 Key Features

### 🧠 Edge ML Inference
*   **Zero-Latency**: 100% On-device scoring using a compiled decision tree engine.
*   **Privacy-Preserving**: Raw SMS logs never leave the phone.
*   **Robust**: Hardcoded logic guarantees stability even in unstable environments.

### 🔐 Digital Trust (New!)
*   **Ed25519 Signing**: Every score is cryptographically signed by the engine.
*   **QR Verification**: Users can share a QR code containing their *Verified Score* and *Signature*. Lenders can verify authenticity offline.
*   **Anonymous**: The verifying lender sees the score, not the bank statements.

### 💼 Financial Infrastructure
*   **Universal Parser**: Supports M-Pesa, Airtel, T-Kash, Hustler Fund, Fuliza, and major banks.
*   **Encrypted Vault**: History is stored in a refined SQLCipher database (AES-256), protected by the Android Keystore.
*   **Feature Vector**: Extracts 20+ signals including *Gambling Ratio*, *Emergency Loan Reliance*, and *Community Lending Diversity*.

---

## 🏗️ Architecture

```mermaid
graph TD
    SMS[Raw SMS Logs] -->|Regex Parser| Vector[Feature Vector (20D)]
    Vector -->|Go Engine| Score[Credit Probability]
    Score -->|Security Module| Signed[Ed25519 Signed Certificate]
    Signed -->|App Bridge| UI[React Native Dashboard]
    UI -->|AES-256| DB[(Encrypted SQLCipher)]
    UI -->|QR Code| Lender[Digital Verification]
```

---

## 🚀 Quick Start

### Prerequisites
*   **Go** 1.20+
*   **Node.js** 18+
*   **Android Studio** (with NDK)

### 1. Build the Engine (Backend/Test)
```bash
go build ./...
go test ./pkg/...
```

### 2. Run the Mobile App
The mobile app includes the compiled Go engine as a native library.

```bash
# Install JS dependencies
cd MobileApp
npm install

# (Optional) Recompile the Go Engine Native Library
cd ..
gomobile bind -v -target=android -androidapi 21 -o MobileApp/android/app/libs/borehole.aar ./pkg/mobile

# Launch on Android Emulator
cd MobileApp
npm run android
```

### 3. Simulating Transactions
Use `adb` to inject fake M-Pesa SMS messages into the emulator to test the scoring logic.

```powershell
# Simulate a large deposit (High Income Signal)
adb emu sms send M-PESA "RC9999ZZ Confirmed. You have received Ksh75,000.00 from ELON MUSK on 28/1/26 at 1:00 PM."
```

---

## 📊 Feature Vector Specifications

| Index | Feature Family | Description |
|-------|----------------|-------------|
| 0-5   | **Cash Flow**  | Income, Expenses, Net Flow, Txn Frequency, Max Txn Size |
| 6-7   | **Risk Flags** | Gambling Index (% of spend), Utility Payments Ratio |
| 8-9   | **Liquidity**  | Fuliza (Overdraft) Usage & Repayment Rate |
| 13-17 | **Ecosystem**  | Hustler Fund, Okoa Jahazi, Airtel Money Volume, Lender Diversity |
| 18-19 | **Stability**  | Savings Rate (MMF/M-Shwari), Banking Activity |

---

## 🛡️ Security Model

1.  **Input Integrity**: The engine only accepts system-level SMS broadcasts (on supported devices) or user-pasted verify-copy.
2.  **Execution Integrity**: The scoring logic is compiled native code, making it resistant to simple runtime tampering.
3.  **Storage Confidentiality**: All persistence is handled by `react-native-sqlite-storage` with **SQLCipher**, keyed by the hardware-backed **Android Keystore**.

---

## License
MIT License. Open Infrastructure for Financial Inclusion.
