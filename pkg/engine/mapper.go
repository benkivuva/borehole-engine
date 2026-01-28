package engine

import (
	"borehole/core/pkg/parser"
	"math"
)

const (
	FeatureCount = 20
)

// MapFeatures transforms raw transactions into a 20-dimension feature vector.
// This is decoupled from the inference engine to allow independent testing/evolution.
func MapFeatures(txns []parser.Transaction) []float64 {
	features := make([]float64, FeatureCount)
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
		hustlerBalance float64
		okoaCount      float64
		airtelVolume   float64
		mmfDeposits    float64
		bankTxnCount   float64
		okoaAmount     float64
		amounts        = make([]float64, 0, len(txns))
		incomeAmounts  = make([]float64, 0, len(txns)/2)
		lenders        = make(map[string]bool)
	)

	for _, txn := range txns {
		amounts = append(amounts, txn.Amount)
		if txn.Amount > maxTxn {
			maxTxn = txn.Amount
		}

		switch txn.Type {
		case parser.TxnMPesaReceived, parser.TxnTKashReceived, parser.TxnAirtelReceived:
			totalIncome += txn.Amount
			incomeAmounts = append(incomeAmounts, txn.Amount)
			if txn.Type == parser.TxnAirtelReceived {
				airtelVolume += txn.Amount
			}
		case parser.TxnMPesaSent, parser.TxnTKashSent, parser.TxnAirtelSent:
			totalExpenses += txn.Amount
			p2pSends += txn.Amount
			if txn.Type == parser.TxnAirtelSent {
				airtelVolume += txn.Amount
			}
		case parser.TxnMPesaPaybill, parser.TxnMPesaBuyGoods:
			totalExpenses += txn.Amount
			utilitySpend += txn.Amount * 0.3
		case parser.TxnFulizaLoan:
			fulizaBorrowed += txn.Amount
			totalIncome += txn.Amount
		case parser.TxnFulizaRepay:
			fulizaRepaid += txn.Amount
			totalExpenses += txn.Amount
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
		case parser.TxnOkoaReceived:
			okoaCount++
			totalIncome += txn.Amount
			if txn.Balance > 0 {
				okoaAmount = txn.Balance
			} else {
				okoaAmount += txn.Amount
			}
		case parser.TxnOkoaDebt:
			okoaCount++
			if txn.Balance > 0 {
				okoaAmount = txn.Balance
			} else if txn.Amount > 0 {
				okoaAmount += txn.Amount
			}
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
		case parser.TxnMMFDeposit:
			mmfDeposits += txn.Amount
			totalExpenses += txn.Amount
		case parser.TxnMMFWithdraw:
			totalIncome += txn.Amount
		case parser.TxnBankDeposit:
			bankTxnCount++
			totalExpenses += txn.Amount
		case parser.TxnBankWithdraw:
			bankTxnCount++
			totalIncome += txn.Amount
		case parser.TxnGambling:
			gamblingSpend += txn.Amount
			totalExpenses += txn.Amount
		}
	}

	// 20-Dimension Mapping
	features[0] = totalIncome
	features[1] = totalExpenses
	features[2] = safeDiv(totalIncome, totalExpenses) // Profitability Ratio
	features[3] = float64(len(txns))
	features[4] = maxTxn
	features[5] = coefficientOfVariation(incomeAmounts)
	features[6] = safeDiv(gamblingSpend, totalExpenses)
	features[7] = safeDiv(utilitySpend, totalExpenses)
	features[8] = safeDiv(fulizaBorrowed, totalIncome)
	features[9] = safeDiv(fulizaRepaid, fulizaBorrowed)
	features[10] = safeDiv(p2pSends, totalExpenses)
	features[11] = stdDev(amounts)
	features[12] = math.Min(float64(len(txns)), 30) // Days Active Approx
	features[13] = hustlerBalance
	features[14] = okoaCount
	features[15] = airtelVolume
	features[16] = float64(len(lenders))
	features[17] = safeDiv(okoaAmount+fulizaBorrowed, totalIncome) // Emergency Reliance
	features[18] = safeDiv(mmfDeposits, totalIncome)               // Savings Rate
	features[19] = bankTxnCount

	return features
}

// Utility functions moved from engine.go for modularity

func safeDiv(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func sum(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return sum(values) / float64(len(values))
}

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
