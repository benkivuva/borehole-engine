package parser

import "regexp"

// Pre-compiled regex patterns for Kenyan mobile money SMS formats.
// These are global but immutable, safe for concurrent use.
// Named capture groups are used for readable extraction.

// M-Pesa 2026 UA series patterns
var (
	// mpesaReceivedPattern matches: "UA1234ABCD Confirmed. You have received Ksh1,500.00 from JOHN DOE 0712345678..."
	mpesaReceivedPattern = regexp.MustCompile(
		`(?i)(?P<refcode>UA[A-Z0-9]{8,10})\s+[Cc]onfirmed\.?\s+[Yy]ou\s+have\s+received\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+from\s+(?P<sender>[A-Z\s]+\d*)`,
	)

	// mpesaSentPattern matches: "UA1234ABCD Confirmed. Ksh500.00 sent to JANE DOE 0798765432..."
	mpesaSentPattern = regexp.MustCompile(
		`(?i)(?P<refcode>UA[A-Z0-9]{8,10})\s+[Cc]onfirmed\.?\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+sent\s+to\s+(?P<recipient>[A-Z\s]+\d*)`,
	)

	// mpesaPaybillPattern matches: "UA1234ABCD Confirmed. Ksh1,000.00 paid to KPLC. Account Number 12345..."
	mpesaPaybillPattern = regexp.MustCompile(
		`(?i)(?P<refcode>UA[A-Z0-9]{8,10})\s+[Cc]onfirmed\.?\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+paid\s+to\s+(?P<account>[A-Z0-9\s]+)`,
	)

	// mpesaBuyGoodsPattern matches: "UA1234ABCD Confirmed. Ksh200.00 paid to SUPERMARKET Till Number 123456..."
	mpesaBuyGoodsPattern = regexp.MustCompile(
		`(?i)(?P<refcode>UA[A-Z0-9]{8,10})\s+[Cc]onfirmed\.?\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+paid\s+to\s+(?P<merchant>[A-Z\s]+)\s*[Tt]ill`,
	)
)

// Fuliza patterns
var (
	// fulizaLoanPattern matches: "Fuliza M-PESA. You have borrowed Ksh2,000.00..."
	fulizaLoanPattern = regexp.MustCompile(
		`(?i)Fuliza.*[Yy]ou\s+have\s+borrowed\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// fulizaRepayPattern matches: "Fuliza M-PESA. You have repaid Ksh500.00..."
	fulizaRepayPattern = regexp.MustCompile(
		`(?i)Fuliza.*[Yy]ou\s+have\s+repaid\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)`,
	)
)

// T-Kash patterns (Telkom Kenya)
var (
	// tkashReceivedPattern matches: "T-Kash: You have received Ksh1,000.00 from JOHN DOE..."
	tkashReceivedPattern = regexp.MustCompile(
		`(?i)T-Kash.*[Yy]ou\s+have\s+received\s+Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+from\s+(?P<sender>[A-Z\s]+)`,
	)

	// tkashSentPattern matches: "T-Kash: Ksh500.00 sent to JANE DOE..."
	tkashSentPattern = regexp.MustCompile(
		`(?i)T-Kash.*Ksh\s*(?P<amt>[\d,]+\.?\d*)\s+sent\s+to\s+(?P<recipient>[A-Z\s]+)`,
	)
)

// Gambling platform patterns (betting companies in Kenya)
var (
	// gamblingPattern matches any mention of major Kenyan betting platforms
	gamblingPattern = regexp.MustCompile(
		`(?i)(Betika|SportPesa|Mozzart|Odibets|Betway|1xBet|Betin|Dafabet|22Bet|Helabet)`,
	)

	// amountPattern is a generic pattern to extract amounts from any SMS
	amountPattern = regexp.MustCompile(
		`Ksh\s*(?P<amt>[\d,]+\.?\d*)`,
	)
)

// Utility company patterns
var (
	// utilityPattern matches common Kenyan utility providers
	utilityPattern = regexp.MustCompile(
		`(?i)(KPLC|Kenya\s+Power|Nairobi\s+Water|Safaricom\s+Home|Zuku|DSTV|GOtv|StarTimes)`,
	)
)
