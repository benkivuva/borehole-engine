// Package parser provides high-performance SMS parsing for Kenyan mobile money transactions.
// Optimized for mobile CPU with minimal allocations.
package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TransactionType represents the category of a mobile money transaction.
type TransactionType int

const (
	TxnUnknown TransactionType = iota
	// M-Pesa types
	TxnMPesaReceived
	TxnMPesaSent
	TxnMPesaPaybill
	TxnMPesaBuyGoods
	// Fuliza types
	TxnFulizaLoan
	TxnFulizaRepay
	// T-Kash types
	TxnTKashReceived
	TxnTKashSent
	// Airtel Money types
	TxnAirtelReceived
	TxnAirtelSent
	// Hustler Fund types
	TxnHustlerLoan
	TxnHustlerRepay
	// Okoa Jahazi types
	TxnOkoaReceived
	TxnOkoaDebt
	// Digital Lender types
	TxnDigitalLoan
	TxnDigitalRepay
	// MMF Savings types
	TxnMMFDeposit
	TxnMMFWithdraw
	// Bank types
	TxnBankDeposit
	TxnBankWithdraw
	// Other types
	TxnGambling
	TxnUtility
)

// String returns the string representation of a TransactionType.
func (t TransactionType) String() string {
	switch t {
	case TxnMPesaReceived:
		return "MPESA_RECEIVED"
	case TxnMPesaSent:
		return "MPESA_SENT"
	case TxnMPesaPaybill:
		return "MPESA_PAYBILL"
	case TxnMPesaBuyGoods:
		return "MPESA_BUYGOODS"
	case TxnFulizaLoan:
		return "FULIZA_LOAN"
	case TxnFulizaRepay:
		return "FULIZA_REPAY"
	case TxnTKashReceived:
		return "TKASH_RECEIVED"
	case TxnTKashSent:
		return "TKASH_SENT"
	case TxnAirtelReceived:
		return "AIRTEL_RECEIVED"
	case TxnAirtelSent:
		return "AIRTEL_SENT"
	case TxnHustlerLoan:
		return "HUSTLER_LOAN"
	case TxnHustlerRepay:
		return "HUSTLER_REPAY"
	case TxnOkoaReceived:
		return "OKOA_RECEIVED"
	case TxnOkoaDebt:
		return "OKOA_DEBT"
	case TxnDigitalLoan:
		return "DIGITAL_LOAN"
	case TxnDigitalRepay:
		return "DIGITAL_REPAY"
	case TxnMMFDeposit:
		return "MMF_DEPOSIT"
	case TxnMMFWithdraw:
		return "MMF_WITHDRAW"
	case TxnBankDeposit:
		return "BANK_DEPOSIT"
	case TxnBankWithdraw:
		return "BANK_WITHDRAW"
	case TxnGambling:
		return "GAMBLING"
	case TxnUtility:
		return "UTILITY"
	default:
		return "UNKNOWN"
	}
}

// Transaction represents a parsed mobile money transaction.
// Fields are optimized for zero-copy where possible.
type Transaction struct {
	Type      TransactionType
	RefCode   string
	Amount    float64
	Balance   float64
	Timestamp time.Time
	Recipient string
	Sender    string
	Lender    string // For digital lender identification
	RawText   string
}

// ScoreResult contains the credit scoring output.
type ScoreResult struct {
	Score    float64   `json:"score"`
	Features []float64 `json:"features"`
	TxnCount int       `json:"txn_count"`
}

// Parser defines the interface for parsing SMS logs.
type Parser interface {
	ParseLogs(ctx context.Context, logs []string) ([]Transaction, error)
}

// DefaultParser implements the Parser interface with optimized parsing.
type DefaultParser struct{}

// NewParser creates a new Parser instance.
func NewParser() Parser {
	return &DefaultParser{}
}

// ParseLogs parses a slice of SMS logs into transactions.
// It uses context for cancellation support and pre-allocates slices
// to minimize garbage collection on mobile devices.
func (p *DefaultParser) ParseLogs(ctx context.Context, logs []string) ([]Transaction, error) {
	if len(logs) == 0 {
		return []Transaction{}, nil
	}

	// Pre-allocate to minimize allocations
	txns := make([]Transaction, 0, len(logs))

	for i, log := range logs {
		// Check context cancellation every 100 logs to balance
		// responsiveness with performance
		if i%100 == 0 {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("parsing cancelled at log %d: %w", i, ctx.Err())
			default:
			}
		}

		txn, err := parseSingleLog(log)
		if err != nil {
			// Skip unparseable logs - common in real SMS data
			continue
		}
		txns = append(txns, txn)
	}

	return txns, nil
}

// parseSingleLog parses a single SMS message into a Transaction.
// Uses keyword-based fast path before regex matching for performance.
func parseSingleLog(log string) (Transaction, error) {
	txn := Transaction{
		Type:    TxnUnknown,
		RawText: log,
	}

	// Convert to uppercase once for keyword checking
	logUpper := strings.ToUpper(log)

	// Fast keyword-based routing to avoid unnecessary regex matching
	switch {
	case strings.Contains(logUpper, "AIRTEL") || strings.Contains(logUpper, "AM1"):
		return parseAirtel(log, txn)

	case strings.Contains(logUpper, "HUSTLER"):
		return parseHustler(log, txn)

	case strings.Contains(logUpper, "OKOA"):
		return parseOkoa(log, txn)

	case strings.Contains(logUpper, "M-SHWARI") || strings.Contains(logUpper, "MALI") ||
		strings.Contains(logUpper, "STAWI") || strings.Contains(logUpper, "KCB M-PESA"):
		return parseMMF(log, txn)

	case strings.Contains(logUpper, "TALA") || strings.Contains(logUpper, "BRANCH") ||
		strings.Contains(logUpper, "ZENKA") || strings.Contains(logUpper, "ZASH") ||
		strings.Contains(logUpper, "OKOLEA"):
		return parseDigitalLender(log, txn)

	case strings.Contains(logUpper, "T-KASH"):
		return parseTKash(log, txn)

	case strings.Contains(logUpper, "FULIZA"):
		return parseFuliza(log, txn)

	default:
		// Fall through to M-Pesa and other patterns
		return parseMPesaAndOthers(log, txn)
	}
}

// parseAirtel handles Airtel Money transactions.
func parseAirtel(log string, txn Transaction) (Transaction, error) {
	if match := airtelReceivedPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnAirtelReceived
		txn.RefCode = getNamedGroup(airtelReceivedPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(airtelReceivedPattern, match, "amt"))
		txn.Sender = getNamedGroup(airtelReceivedPattern, match, "sender")
		return txn, nil
	}

	if match := airtelSentPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnAirtelSent
		txn.RefCode = getNamedGroup(airtelSentPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(airtelSentPattern, match, "amt"))
		txn.Recipient = getNamedGroup(airtelSentPattern, match, "recipient")
		return txn, nil
	}

	// Generic Airtel detection with amount extraction
	if airtelGenericPattern.MatchString(log) {
		if match := amountPattern.FindStringSubmatch(log); match != nil {
			txn.Type = TxnAirtelReceived // Default to received
			txn.Amount = parseAmount(getNamedGroup(amountPattern, match, "amt"))
			return txn, nil
		}
	}

	return txn, fmt.Errorf("no Airtel pattern matched")
}

// parseHustler handles Hustler Fund transactions.
func parseHustler(log string, txn Transaction) (Transaction, error) {
	if match := hustlerLoanPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnHustlerLoan
		txn.Amount = parseAmount(getNamedGroup(hustlerLoanPattern, match, "amt"))
		txn.Lender = "Hustler Fund"
		return txn, nil
	}

	if match := hustlerRepayPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnHustlerRepay
		txn.Amount = parseAmount(getNamedGroup(hustlerRepayPattern, match, "amt"))
		txn.Lender = "Hustler Fund"
		return txn, nil
	}

	if match := hustlerBalancePattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnHustlerLoan
		txn.Balance = parseAmount(getNamedGroup(hustlerBalancePattern, match, "amt"))
		txn.Lender = "Hustler Fund"
		return txn, nil
	}

	return txn, fmt.Errorf("no Hustler pattern matched")
}

// parseOkoa handles Okoa Jahazi transactions.
func parseOkoa(log string, txn Transaction) (Transaction, error) {
	if match := okoaReceivedPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnOkoaReceived
		txn.Amount = parseAmount(getNamedGroup(okoaReceivedPattern, match, "amt"))
		return txn, nil
	}

	if match := okoaDebtPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnOkoaDebt
		txn.Balance = parseAmount(getNamedGroup(okoaDebtPattern, match, "amt"))
		return txn, nil
	}

	if match := okoaRepayPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnOkoaDebt
		txn.Amount = parseAmount(getNamedGroup(okoaRepayPattern, match, "amt"))
		return txn, nil
	}

	return txn, fmt.Errorf("no Okoa pattern matched")
}

// parseMMF handles Money Market Fund savings (M-Shwari, KCB M-Pesa, Mali, Stawi).
func parseMMF(log string, txn Transaction) (Transaction, error) {
	// M-Shwari
	if match := mshwariDepositPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMMFDeposit
		txn.Amount = parseAmount(getNamedGroup(mshwariDepositPattern, match, "amt"))
		txn.Recipient = "M-Shwari"
		return txn, nil
	}
	if match := mshwariWithdrawPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMMFWithdraw
		txn.Amount = parseAmount(getNamedGroup(mshwariWithdrawPattern, match, "amt"))
		txn.Sender = "M-Shwari"
		return txn, nil
	}

	// KCB M-Pesa
	if match := kcbMpesaSavePattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMMFDeposit
		txn.Amount = parseAmount(getNamedGroup(kcbMpesaSavePattern, match, "amt"))
		txn.Recipient = "KCB M-Pesa"
		return txn, nil
	}

	// Mali
	if match := maliSavePattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMMFDeposit
		txn.Amount = parseAmount(getNamedGroup(maliSavePattern, match, "amt"))
		txn.Recipient = "Mali"
		return txn, nil
	}

	// Stawi
	if match := stawiSavePattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMMFDeposit
		txn.Amount = parseAmount(getNamedGroup(stawiSavePattern, match, "amt"))
		txn.Recipient = "Stawi"
		return txn, nil
	}

	// Generic MMF with amount extraction
	if mmfPattern.MatchString(log) {
		if match := amountPattern.FindStringSubmatch(log); match != nil {
			txn.Type = TxnMMFDeposit
			txn.Amount = parseAmount(getNamedGroup(amountPattern, match, "amt"))
			return txn, nil
		}
	}

	return txn, fmt.Errorf("no MMF pattern matched")
}

// parseDigitalLender handles digital loan app transactions (Tala, Branch, etc.).
func parseDigitalLender(log string, txn Transaction) (Transaction, error) {
	if match := loanDisbursementPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnDigitalLoan
		txn.Amount = parseAmount(getNamedGroup(loanDisbursementPattern, match, "amt"))
		txn.Lender = getNamedGroup(loanDisbursementPattern, match, "lender")
		return txn, nil
	}

	if match := loanRepaymentPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnDigitalRepay
		txn.Amount = parseAmount(getNamedGroup(loanRepaymentPattern, match, "amt"))
		txn.Lender = getNamedGroup(loanRepaymentPattern, match, "lender")
		return txn, nil
	}

	// Generic lender detection
	if digitalLenderPattern.MatchString(log) {
		if match := amountPattern.FindStringSubmatch(log); match != nil {
			// Infer loan or repay based on keywords
			logUpper := strings.ToUpper(log)
			if strings.Contains(logUpper, "REPAY") || strings.Contains(logUpper, "PAID") {
				txn.Type = TxnDigitalRepay
			} else {
				txn.Type = TxnDigitalLoan
			}
			txn.Amount = parseAmount(getNamedGroup(amountPattern, match, "amt"))
			// Extract lender name
			if lender := digitalLenderPattern.FindString(log); lender != "" {
				txn.Lender = lender
			}
			return txn, nil
		}
	}

	return txn, fmt.Errorf("no digital lender pattern matched")
}

// parseTKash handles T-Kash transactions.
func parseTKash(log string, txn Transaction) (Transaction, error) {
	if match := tkashReceivedPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnTKashReceived
		txn.Amount = parseAmount(getNamedGroup(tkashReceivedPattern, match, "amt"))
		txn.Sender = getNamedGroup(tkashReceivedPattern, match, "sender")
		return txn, nil
	}

	if match := tkashSentPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnTKashSent
		txn.Amount = parseAmount(getNamedGroup(tkashSentPattern, match, "amt"))
		txn.Recipient = getNamedGroup(tkashSentPattern, match, "recipient")
		return txn, nil
	}

	return txn, fmt.Errorf("no T-Kash pattern matched")
}

// parseFuliza handles Fuliza loan transactions.
func parseFuliza(log string, txn Transaction) (Transaction, error) {
	if match := fulizaLoanPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnFulizaLoan
		txn.Amount = parseAmount(getNamedGroup(fulizaLoanPattern, match, "amt"))
		return txn, nil
	}

	if match := fulizaRepayPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnFulizaRepay
		txn.Amount = parseAmount(getNamedGroup(fulizaRepayPattern, match, "amt"))
		return txn, nil
	}

	return txn, fmt.Errorf("no Fuliza pattern matched")
}

// parseMPesaAndOthers handles M-Pesa, gambling, and other patterns.
func parseMPesaAndOthers(log string, txn Transaction) (Transaction, error) {
	// M-Pesa patterns
	if match := mpesaReceivedPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMPesaReceived
		txn.RefCode = getNamedGroup(mpesaReceivedPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(mpesaReceivedPattern, match, "amt"))
		txn.Sender = getNamedGroup(mpesaReceivedPattern, match, "sender")
		return txn, nil
	}

	if match := mpesaSentPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMPesaSent
		txn.RefCode = getNamedGroup(mpesaSentPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(mpesaSentPattern, match, "amt"))
		txn.Recipient = getNamedGroup(mpesaSentPattern, match, "recipient")
		return txn, nil
	}

	if match := mpesaPaybillPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMPesaPaybill
		txn.RefCode = getNamedGroup(mpesaPaybillPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(mpesaPaybillPattern, match, "amt"))
		txn.Recipient = getNamedGroup(mpesaPaybillPattern, match, "account")
		return txn, nil
	}

	if match := mpesaBuyGoodsPattern.FindStringSubmatch(log); match != nil {
		txn.Type = TxnMPesaBuyGoods
		txn.RefCode = getNamedGroup(mpesaBuyGoodsPattern, match, "refcode")
		txn.Amount = parseAmount(getNamedGroup(mpesaBuyGoodsPattern, match, "amt"))
		txn.Recipient = getNamedGroup(mpesaBuyGoodsPattern, match, "merchant")
		return txn, nil
	}

	// Check for gambling platforms
	if gamblingPattern.MatchString(log) {
		txn.Type = TxnGambling
		if match := amountPattern.FindStringSubmatch(log); match != nil {
			txn.Amount = parseAmount(getNamedGroup(amountPattern, match, "amt"))
		}
		return txn, nil
	}

	// Check for bank transfers
	if bankTransferPattern.MatchString(log) {
		if match := bankDepositPattern.FindStringSubmatch(log); match != nil {
			txn.Type = TxnBankDeposit
			txn.Amount = parseAmount(getNamedGroup(bankDepositPattern, match, "amt"))
			txn.Recipient = getNamedGroup(bankDepositPattern, match, "bank")
			return txn, nil
		}
		if match := bankWithdrawPattern.FindStringSubmatch(log); match != nil {
			txn.Type = TxnBankWithdraw
			txn.Amount = parseAmount(getNamedGroup(bankWithdrawPattern, match, "amt"))
			txn.Sender = getNamedGroup(bankWithdrawPattern, match, "bank")
			return txn, nil
		}
	}

	return txn, fmt.Errorf("no pattern matched for log")
}

// parseAmount converts Kenyan SMS amount format to float64.
// Handles formats like "Ksh1,500.00", "Ksh 1500", "KES 1,234.56"
func parseAmount(s string) float64 {
	if s == "" {
		return 0
	}

	// Remove common prefixes and whitespace
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "Ksh")
	s = strings.TrimPrefix(s, "ksh")
	s = strings.TrimPrefix(s, "KES")
	s = strings.TrimPrefix(s, "kes")
	s = strings.TrimSpace(s)

	// Remove commas (Kenyan format uses commas for thousands)
	s = strings.ReplaceAll(s, ",", "")

	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return amount
}

// getNamedGroup extracts a named capture group from regex match.
func getNamedGroup(re *regexp.Regexp, match []string, name string) string {
	for i, groupName := range re.SubexpNames() {
		if groupName == name && i < len(match) {
			return match[i]
		}
	}
	return ""
}
