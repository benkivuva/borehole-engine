package parser

import "regexp"

// Pre-compiled regex patterns for Kenyan mobile money SMS formats.
// These are global but immutable, safe for concurrent use.
// Named capture groups are used for readable extraction.

// =============================================================================
// M-Pesa 2026 UA series patterns
// =============================================================================
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

// =============================================================================
// Fuliza patterns
// =============================================================================
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

// =============================================================================
// T-Kash patterns (Telkom Kenya)
// =============================================================================
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

// =============================================================================
// Airtel Money patterns
// =============================================================================
var (
	// airtelReceivedPattern matches: "Transaction ID: AM12345678. You have received Ksh1,000.00 from..."
	airtelReceivedPattern = regexp.MustCompile(
		`(?i)Transaction\s+ID[:\s]*(?P<refcode>AM[A-Z0-9]+).*[Yy]ou\s+have\s+received\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+from\s+(?P<sender>[A-Z\s]+)`,
	)

	// airtelSentPattern matches: "Transaction ID: AM12345678. Ksh500.00 sent to..."
	airtelSentPattern = regexp.MustCompile(
		`(?i)Transaction\s+ID[:\s]*(?P<refcode>AM[A-Z0-9]+).*(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+sent\s+to\s+(?P<recipient>[A-Z\s]+)`,
	)

	// airtelGenericPattern matches generic Airtel Money keyword
	airtelGenericPattern = regexp.MustCompile(`(?i)Airtel\s*Money`)
)

// =============================================================================
// Hustler Fund patterns (Government of Kenya)
// =============================================================================
var (
	// hustlerLoanPattern matches: "Hustler Fund. You have been disbursed Ksh500.00..."
	hustlerLoanPattern = regexp.MustCompile(
		`(?i)Hustler\s+Fund.*(?:disbursed|received)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// hustlerRepayPattern matches: "Hustler Fund. You have repaid Ksh200.00..."
	hustlerRepayPattern = regexp.MustCompile(
		`(?i)Hustler\s+Fund.*repaid\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// hustlerBalancePattern matches: "Hustler Fund. Your loan balance is Ksh300.00..."
	hustlerBalancePattern = regexp.MustCompile(
		`(?i)Hustler\s+Fund.*(?:balance|limit)\s+(?:is\s+)?(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)
)

// =============================================================================
// Okoa Jahazi patterns (Safaricom airtime credit)
// =============================================================================
var (
	// okoaReceivedPattern matches: "You have received Ksh50 Okoa Jahazi..."
	okoaReceivedPattern = regexp.MustCompile(
		`(?i)(?:received|got)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+Okoa\s+Jahazi`,
	)

	// okoaDebtPattern matches: "Your Okoa debt is Ksh50..."
	okoaDebtPattern = regexp.MustCompile(
		`(?i)Okoa\s+(?:Jahazi\s+)?debt\s+(?:is\s+)?(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// okoaRepayPattern matches: "Okoa Jahazi. You have repaid Ksh50..."
	okoaRepayPattern = regexp.MustCompile(
		`(?i)Okoa\s+Jahazi.*repaid\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)
)

// =============================================================================
// Digital Lenders patterns (Tala, Branch, Zenka, etc.)
// =============================================================================
var (
	// digitalLenderPattern matches SMS from major Kenyan digital lenders
	digitalLenderPattern = regexp.MustCompile(
		`(?i)(Tala|Branch|Zenka|Zash|Okolea|KCB-MPESA|Fuliza|Timiza|Berry|Kashway)`,
	)

	// loanDisbursementPattern matches: "You have received Ksh5,000.00 from Tala..."
	loanDisbursementPattern = regexp.MustCompile(
		`(?i)(?:received|disbursed)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+(?:from\s+)?(?P<lender>Tala|Branch|Zenka|Zash|Okolea)`,
	)

	// loanRepaymentPattern matches: "Ksh1,000.00 received by Tala..."
	loanRepaymentPattern = regexp.MustCompile(
		`(?i)(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+(?:paid|received\s+by)\s+(?P<lender>Tala|Branch|Zenka|Zash|Okolea)`,
	)
)

// =============================================================================
// MMF Savings patterns (M-Shwari, KCB M-Pesa, Mali, Stawi)
// =============================================================================
var (
	// mshwariDepositPattern matches: "M-Shwari. You have deposited Ksh1,000.00..."
	mshwariDepositPattern = regexp.MustCompile(
		`(?i)M-Shwari.*(?:deposited|saved|transferred)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// mshwariWithdrawPattern matches: "M-Shwari. You have withdrawn Ksh500.00..."
	mshwariWithdrawPattern = regexp.MustCompile(
		`(?i)M-Shwari.*(?:withdrawn|transferred)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// kcbMpesaPattern matches KCB M-Pesa savings
	kcbMpesaSavePattern = regexp.MustCompile(
		`(?i)KCB\s*M-?PESA.*(?:deposited|saved|transferred)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// maliPattern matches Mali (Safaricom MMF)
	maliSavePattern = regexp.MustCompile(
		`(?i)Mali.*(?:deposited|invested|saved)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// stawiPattern matches Stawi (NCBA-Safaricom)
	stawiSavePattern = regexp.MustCompile(
		`(?i)Stawi.*(?:deposited|saved)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)

	// genericMMFPattern matches any MMF-related keywords
	mmfPattern = regexp.MustCompile(
		`(?i)(M-Shwari|KCB\s*M-?PESA|Mali|Stawi|Lock\s+Savings)`,
	)
)

// =============================================================================
// Bank Transfer patterns
// =============================================================================
var (
	// bankTransferPattern matches transfers to/from banks
	bankTransferPattern = regexp.MustCompile(
		`(?i)(KCB|Equity|Co-?op(?:erative)?|NCBA|Stanbic|Absa|DTB|I&M|Family\s+Bank|Bank\s+of\s+Africa)`,
	)

	// bankDepositPattern matches: "Deposited Ksh5,000.00 to Equity Bank..."
	bankDepositPattern = regexp.MustCompile(
		`(?i)(?:deposited|transferred|sent)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+(?:to\s+)?(?P<bank>KCB|Equity|Co-?op|NCBA|Stanbic|Absa)`,
	)

	// bankWithdrawPattern matches: "Withdrawn Ksh2,000.00 from Equity Bank..."
	bankWithdrawPattern = regexp.MustCompile(
		`(?i)(?:withdrawn|received)\s+(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)\s+(?:from\s+)?(?P<bank>KCB|Equity|Co-?op|NCBA|Stanbic|Absa)`,
	)
)

// =============================================================================
// Gambling platform patterns
// =============================================================================
var (
	// gamblingPattern matches any mention of major Kenyan betting platforms
	gamblingPattern = regexp.MustCompile(
		`(?i)(Betika|SportPesa|Mozzart|Odibets|Betway|1xBet|Betin|Dafabet|22Bet|Helabet)`,
	)

	// amountPattern is a generic pattern to extract amounts from any SMS
	amountPattern = regexp.MustCompile(
		`(?:Ksh|KES)\s*(?P<amt>[\d,]+\.?\d*)`,
	)
)

// =============================================================================
// Utility company patterns
// =============================================================================
var (
	// utilityPattern matches common Kenyan utility providers
	utilityPattern = regexp.MustCompile(
		`(?i)(KPLC|Kenya\s+Power|Nairobi\s+Water|Safaricom\s+Home|Zuku|DSTV|GOtv|StarTimes)`,
	)
)
