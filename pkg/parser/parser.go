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
	TxnMPesaReceived
	TxnMPesaSent
	TxnMPesaPaybill
	TxnMPesaBuyGoods
	TxnFulizaLoan
	TxnFulizaRepay
	TxnTKashReceived
	TxnTKashSent
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
func parseSingleLog(log string) (Transaction, error) {
	txn := Transaction{
		Type:    TxnUnknown,
		RawText: log,
	}

	// Try each pattern in order of likelihood
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

	// Check for gambling platforms
	if gamblingPattern.MatchString(log) {
		txn.Type = TxnGambling
		// Try to extract amount from gambling messages
		if match := amountPattern.FindStringSubmatch(log); match != nil {
			txn.Amount = parseAmount(match[1])
		}
		return txn, nil
	}

	return txn, fmt.Errorf("no pattern matched for log")
}

// parseAmount converts Kenyan SMS amount format to float64.
// Handles formats like "Ksh1,500.00", "Ksh 1500", "1,234.56"
func parseAmount(s string) float64 {
	if s == "" {
		return 0
	}

	// Remove common prefixes and whitespace
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "Ksh")
	s = strings.TrimPrefix(s, "KES")
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
