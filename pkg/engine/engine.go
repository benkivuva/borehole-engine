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
const featureCount = 22

// Vectorize transforms transactions into a 22-element feature vector.
// Features are deterministic for XGBoost reproducibility.
//
// Feature indices:
//   - 0: total_income         - Sum of all received amounts
//   - 1: total_expenses       - Sum of all sent/paid amounts
//   - 2: net_flow             - income - expenses
//   - 3: avg_txn_amount       - Mean transaction value
//   - 4: txn_count            - Total transaction count
//   - 5: income_regularity    - Coefficient of variation for income
//   - 6: gambling_index       - Gambling spend / total expenses
//   - 7: utility_ratio        - Utility payments / total expenses
//   - 8: fuliza_usage         - Fuliza borrowed / total income
//   - 9: fuliza_repay_rate    - Fuliza repaid / Fuliza borrowed
//   - 10: p2p_ratio           - P2P sends / total expenses
//   - 11: max_single_txn      - Largest single transaction
//   - 12: balance_volatility  - Std dev of transaction amounts
//   - 13: days_active         - Unique days with transactions (simulated)
//   - 14: avg_daily_volume    - Total volume / days active
//   - 15: hustler_balance     - Latest Hustler Fund debt/balance
//   - 16: okoa_frequency      - Count of Okoa Jahazi occurrences
//   - 17: airtel_volume       - Total Airtel Money transaction volume
//   - 18: lender_diversity    - Count of unique digital lenders
//   - 19: emergency_reliance  - (Okoa + Fuliza) / Total Income
//   - 20: savings_rate        - MMF deposits / Total Income
//   - 21: bank_activity       - Count of bank transactions
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
		// New metrics
		hustlerBalance float64
		okoaCount      float64
		airtelVolume   float64
		mmfDeposits    float64
		bankTxnCount   float64
		okoaAmount     float64
		// Tracking
		amounts       = make([]float64, 0, len(txns))
		incomeAmounts = make([]float64, 0, len(txns)/2)
		lenders       = make(map[string]bool) // For lender diversity
	)

	for _, txn := range txns {
		amounts = append(amounts, txn.Amount)

		if txn.Amount > maxTxn {
			maxTxn = txn.Amount
		}

		switch txn.Type {
		// Income types
		case parser.TxnMPesaReceived, parser.TxnTKashReceived, parser.TxnAirtelReceived:
			totalIncome += txn.Amount
			incomeAmounts = append(incomeAmounts, txn.Amount)
			if txn.Type == parser.TxnAirtelReceived {
				airtelVolume += txn.Amount
			}

		// Expense types (P2P)
		case parser.TxnMPesaSent, parser.TxnTKashSent, parser.TxnAirtelSent:
			totalExpenses += txn.Amount
			p2pSends += txn.Amount
			if txn.Type == parser.TxnAirtelSent {
				airtelVolume += txn.Amount
			}

		// Paybill / Buy Goods
		case parser.TxnMPesaPaybill, parser.TxnMPesaBuyGoods:
			totalExpenses += txn.Amount
			utilitySpend += txn.Amount * 0.3 // Heuristic

		// Fuliza
		case parser.TxnFulizaLoan:
			fulizaBorrowed += txn.Amount
			totalIncome += txn.Amount

		case parser.TxnFulizaRepay:
			fulizaRepaid += txn.Amount
			totalExpenses += txn.Amount

		// Hustler Fund
		case parser.TxnHustlerLoan:
			totalIncome += txn.Amount
			if txn.Balance > hustlerBalance {
				hustlerBalance = txn.Balance
			}
			if txn.Amount > 0 && hustlerBalance == 0 {
				hustlerBalance = txn.Amount
			}

		case parser.TxnHustlerRepay:
			totalExpenses += txn.Amount

		// Okoa Jahazi
		case parser.TxnOkoaReceived:
			okoaCount++
			okoaAmount += txn.Amount
			totalIncome += txn.Amount

		case parser.TxnOkoaDebt:
			okoaCount++
			if txn.Balance > 0 {
				okoaAmount = txn.Balance
			}

		// Digital Lenders
		case parser.TxnDigitalLoan:
			totalIncome += txn.Amount
			if txn.Lender != "" {
				lenders[txn.Lender] = true
			}

		case parser.TxnDigitalRepay:
			totalExpenses += txn.Amount
			if txn.Lender != "" {
				lenders[txn.Lender] = true
			}

		// MMF Savings
		case parser.TxnMMFDeposit:
			mmfDeposits += txn.Amount
			totalExpenses += txn.Amount // Savings reduce available balance

		case parser.TxnMMFWithdraw:
			totalIncome += txn.Amount

		// Bank activity
		case parser.TxnBankDeposit:
			bankTxnCount++
			totalExpenses += txn.Amount

		case parser.TxnBankWithdraw:
			bankTxnCount++
			totalIncome += txn.Amount

		// Gambling
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
	features[13] = math.Min(float64(len(txns)), 30)

	// Feature 14: Average Daily Volume
	features[14] = safeDiv(sum(amounts), features[13])

	// Feature 15: Hustler Balance (latest debt amount)
	features[15] = hustlerBalance

	// Feature 16: Okoa Frequency (count of Okoa occurrences)
	features[16] = okoaCount

	// Feature 17: Airtel Volume (total Airtel transaction volume)
	features[17] = airtelVolume

	// Feature 18: Lender Diversity (unique digital lenders)
	features[18] = float64(len(lenders))

	// Feature 19: Emergency Reliance Ratio ((Okoa + Fuliza) / Income)
	emergencyBorrowing := okoaAmount + fulizaBorrowed
	features[19] = safeDiv(emergencyBorrowing, totalIncome)

	// Feature 20: Savings Rate (MMF deposits / Income)
	features[20] = safeDiv(mmfDeposits, totalIncome)

	// Feature 21: Bank Activity Count
	features[21] = bankTxnCount

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
