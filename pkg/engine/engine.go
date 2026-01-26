// Package engine provides feature engineering for credit scoring.
// Transforms parsed transactions into feature vectors for XGBoost model.
package engine

import (
	"math"

	"borehole/core/pkg/parser"
)

// Vectorizer defines the interface for feature extraction.
type Vectorizer interface {
	Vectorize(txns []parser.Transaction) []float64
}

// Engine implements the Vectorizer interface.
type Engine struct{}

// NewEngine creates a new Engine instance.
func NewEngine() Vectorizer {
	return &Engine{}
}

// featureCount is the number of features in the output vector.
const featureCount = 15

// Vectorize transforms transactions into a 15-element feature vector.
// Features are deterministic for XGBoost reproducibility.
//
// Feature indices:
//   - 0: total_income       - Sum of all received amounts
//   - 1: total_expenses     - Sum of all sent/paid amounts
//   - 2: net_flow           - income - expenses
//   - 3: avg_txn_amount     - Mean transaction value
//   - 4: txn_count          - Total transaction count
//   - 5: income_regularity  - Coefficient of variation for income
//   - 6: gambling_index     - Gambling spend / total expenses
//   - 7: utility_ratio      - Utility payments / total expenses
//   - 8: fuliza_usage       - Fuliza borrowed / total income
//   - 9: fuliza_repay_rate  - Fuliza repaid / Fuliza borrowed
//   - 10: p2p_ratio         - P2P sends / total expenses
//   - 11: max_single_txn    - Largest single transaction
//   - 12: balance_volatility - Std dev of transaction amounts
//   - 13: days_active       - Unique days with transactions (simulated)
//   - 14: avg_daily_volume  - Total volume / days active
func (e *Engine) Vectorize(txns []parser.Transaction) []float64 {
	features := make([]float64, featureCount)

	if len(txns) == 0 {
		return features
	}

	var (
		totalIncome    float64
		totalExpenses  float64
		gamblingSpend  float64
		utilitySpend   float64
		fulizaBorrowed float64
		fulizaRepaid   float64
		p2pSends       float64
		maxTxn         float64
		amounts        = make([]float64, 0, len(txns))
		incomeAmounts  = make([]float64, 0, len(txns)/2)
	)

	for _, txn := range txns {
		amounts = append(amounts, txn.Amount)

		if txn.Amount > maxTxn {
			maxTxn = txn.Amount
		}

		switch txn.Type {
		case parser.TxnMPesaReceived, parser.TxnTKashReceived:
			totalIncome += txn.Amount
			incomeAmounts = append(incomeAmounts, txn.Amount)

		case parser.TxnMPesaSent, parser.TxnTKashSent:
			totalExpenses += txn.Amount
			p2pSends += txn.Amount

		case parser.TxnMPesaPaybill, parser.TxnMPesaBuyGoods:
			totalExpenses += txn.Amount
			// Check if utility (simplified heuristic)
			utilitySpend += txn.Amount * 0.3 // Assume 30% of paybill/buygoods is utility

		case parser.TxnFulizaLoan:
			fulizaBorrowed += txn.Amount
			totalIncome += txn.Amount // Fuliza adds to available funds

		case parser.TxnFulizaRepay:
			fulizaRepaid += txn.Amount
			totalExpenses += txn.Amount

		case parser.TxnGambling:
			gamblingSpend += txn.Amount
			totalExpenses += txn.Amount
		}
	}

	// Feature 0: Total Income
	features[0] = totalIncome

	// Feature 1: Total Expenses
	features[1] = totalExpenses

	// Feature 2: Net Flow
	features[2] = totalIncome - totalExpenses

	// Feature 3: Average Transaction Amount
	features[3] = safeDiv(sum(amounts), float64(len(amounts)))

	// Feature 4: Transaction Count
	features[4] = float64(len(txns))

	// Feature 5: Income Regularity (coefficient of variation)
	features[5] = coefficientOfVariation(incomeAmounts)

	// Feature 6: Gambling Index
	features[6] = safeDiv(gamblingSpend, totalExpenses)

	// Feature 7: Utility Ratio
	features[7] = safeDiv(utilitySpend, totalExpenses)

	// Feature 8: Fuliza Usage
	features[8] = safeDiv(fulizaBorrowed, totalIncome)

	// Feature 9: Fuliza Repay Rate
	features[9] = safeDiv(fulizaRepaid, fulizaBorrowed)

	// Feature 10: P2P Ratio
	features[10] = safeDiv(p2pSends, totalExpenses)

	// Feature 11: Max Single Transaction
	features[11] = maxTxn

	// Feature 12: Balance Volatility (std dev of amounts)
	features[12] = stdDev(amounts)

	// Feature 13: Days Active (estimated from transaction count)
	// In production, this would use actual timestamps
	features[13] = math.Min(float64(len(txns)), 30) // Cap at 30 days

	// Feature 14: Average Daily Volume
	features[14] = safeDiv(sum(amounts), features[13])

	return features
}

// safeDiv performs division with zero-check to avoid NaN/Inf.
func safeDiv(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// sum calculates the sum of a float64 slice.
func sum(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// mean calculates the arithmetic mean.
func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return sum(values) / float64(len(values))
}

// stdDev calculates the population standard deviation.
func stdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	m := mean(values)
	var sumSquares float64
	for _, v := range values {
		diff := v - m
		sumSquares += diff * diff
	}

	return math.Sqrt(sumSquares / float64(len(values)))
}

// coefficientOfVariation calculates CV (std dev / mean).
// Lower values indicate more regular income.
func coefficientOfVariation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	m := mean(values)
	if m == 0 {
		return 0
	}

	return stdDev(values) / m
}
