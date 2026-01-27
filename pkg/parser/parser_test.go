package parser

import (
	"context"
	"testing"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Ksh with comma", "Ksh1,500.00", 1500.00},
		{"Ksh no comma", "Ksh500", 500.00},
		{"Ksh with space", "Ksh 1,234.56", 1234.56},
		{"KES uppercase", "KES2,000.00", 2000.00},
		{"KES with space", "KES 3,500", 3500.00},
		{"lowercase ksh", "ksh100", 100.00},
		{"plain number", "5000.50", 5000.50},
		{"number with comma", "10,000", 10000.00},
		{"empty string", "", 0},
		{"invalid", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAmount(tt.input)
			if result != tt.expected {
				t.Errorf("parseAmount(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseSingleLog_MPesa(t *testing.T) {
	tests := []struct {
		name        string
		log         string
		wantType    TransactionType
		wantAmount  float64
		wantRefCode string
	}{
		{
			name:        "M-Pesa received",
			log:         "UA1234ABCDEF Confirmed. You have received Ksh1,500.00 from JOHN DOE 0712345678",
			wantType:    TxnMPesaReceived,
			wantAmount:  1500.00,
			wantRefCode: "UA1234ABCDEF",
		},
		{
			name:        "M-Pesa sent",
			log:         "UA5678EFGHIJ Confirmed. Ksh500.00 sent to JANE DOE 0798765432",
			wantType:    TxnMPesaSent,
			wantAmount:  500.00,
			wantRefCode: "UA5678EFGHIJ",
		},
		{
			name:        "M-Pesa paybill",
			log:         "UA9999XYZABC Confirmed. Ksh1,000.00 paid to KPLC Account 12345",
			wantType:    TxnMPesaPaybill,
			wantAmount:  1000.00,
			wantRefCode: "UA9999XYZABC",
		},
		{
			name:        "M-Pesa received (QKJ prefix)",
			log:         "QKJ3XPYC5T Confirmed. You have received Ksh15,000.00 from SARAH JANE",
			wantType:    TxnMPesaReceived,
			wantAmount:  15000.00,
			wantRefCode: "QKJ3XPYC5T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
			if txn.RefCode != tt.wantRefCode {
				t.Errorf("RefCode = %v, want %v", txn.RefCode, tt.wantRefCode)
			}
		})
	}
}

func TestParseSingleLog_Fuliza(t *testing.T) {
	tests := []struct {
		name       string
		log        string
		wantType   TransactionType
		wantAmount float64
	}{
		{
			name:       "Fuliza loan",
			log:        "Fuliza M-PESA. You have borrowed Ksh2,000.00 from your limit",
			wantType:   TxnFulizaLoan,
			wantAmount: 2000.00,
		},
		{
			name:       "Fuliza repay",
			log:        "Fuliza M-PESA. You have repaid Ksh500.00",
			wantType:   TxnFulizaRepay,
			wantAmount: 500.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
		})
	}
}

func TestParseSingleLog_Airtel(t *testing.T) {
	tests := []struct {
		name       string
		log        string
		wantType   TransactionType
		wantAmount float64
	}{
		{
			name:       "Airtel received",
			log:        "Transaction ID: AM12345678. You have received Ksh1,000.00 from JOHN DOE",
			wantType:   TxnAirtelReceived,
			wantAmount: 1000.00,
		},
		{
			name:       "Airtel sent",
			log:        "Transaction ID: AM87654321. Ksh500.00 sent to JANE DOE",
			wantType:   TxnAirtelSent,
			wantAmount: 500.00,
		},
		{
			name:       "Airtel Money generic",
			log:        "Airtel Money: Your transaction of Ksh200.00 was successful",
			wantType:   TxnAirtelReceived,
			wantAmount: 200.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
		})
	}
}

func TestParseSingleLog_Hustler(t *testing.T) {
	tests := []struct {
		name       string
		log        string
		wantType   TransactionType
		wantAmount float64
	}{
		{
			name:       "Hustler Fund loan",
			log:        "Hustler Fund. You have been disbursed Ksh500.00 to your account",
			wantType:   TxnHustlerLoan,
			wantAmount: 500.00,
		},
		{
			name:       "Hustler Fund repay",
			log:        "Hustler Fund. You have repaid Ksh200.00",
			wantType:   TxnHustlerRepay,
			wantAmount: 200.00,
		},
		{
			name:       "Hustler Fund sent repayment",
			log:        "Confirmed. You have sent Ksh2,000.00 to Hustler Fund on 20/1/26.",
			wantType:   TxnHustlerRepay,
			wantAmount: 2000.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
		})
	}
}

func TestParseSingleLog_Okoa(t *testing.T) {
	tests := []struct {
		name        string
		log         string
		wantType    TransactionType
		wantAmount  float64
		wantBalance float64
	}{
		{
			name:       "Okoa received",
			log:        "You have received Ksh50 Okoa Jahazi airtime credit",
			wantType:   TxnOkoaReceived,
			wantAmount: 50.00,
		},
		{
			name:        "Okoa debt",
			log:         "Your Okoa debt is Ksh50. Please repay",
			wantType:    TxnOkoaDebt,
			wantBalance: 50.00,
		},
		{
			name:        "Okoa received + debt combined",
			log:         "You have received Ksh 100.00 Okoa Jahazi. Your Okoa debt is Ksh 110.00.",
			wantType:    TxnOkoaReceived,
			wantAmount:  100.00,
			wantBalance: 110.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if tt.wantAmount > 0 && txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
			if tt.wantBalance > 0 && txn.Balance != tt.wantBalance {
				t.Errorf("Balance = %v, want %v", txn.Balance, tt.wantBalance)
			}
		})
	}
}

func TestParseSingleLog_MMF(t *testing.T) {
	tests := []struct {
		name          string
		log           string
		wantType      TransactionType
		wantAmount    float64
		wantRecipient string
	}{
		{
			name:          "M-Shwari deposit",
			log:           "M-Shwari. You have deposited Ksh1,000.00 to your savings",
			wantType:      TxnMMFDeposit,
			wantAmount:    1000.00,
			wantRecipient: "M-Shwari",
		},
		{
			name:       "M-Shwari withdraw",
			log:        "M-Shwari. You have withdrawn Ksh500.00",
			wantType:   TxnMMFWithdraw,
			wantAmount: 500.00,
		},
		{
			name:          "KCB M-Pesa save",
			log:           "KCB M-PESA. You have deposited Ksh2,000.00",
			wantType:      TxnMMFDeposit,
			wantAmount:    2000.00,
			wantRecipient: "KCB M-Pesa",
		},
		{
			name:          "Mali invest",
			log:           "Mali. You have invested Ksh500.00 in money market fund",
			wantType:      TxnMMFDeposit,
			wantAmount:    500.00,
			wantRecipient: "Mali",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
			if tt.wantRecipient != "" && txn.Recipient != tt.wantRecipient {
				t.Errorf("Recipient = %v, want %v", txn.Recipient, tt.wantRecipient)
			}
		})
	}
}

func TestParseSingleLog_DigitalLender(t *testing.T) {
	tests := []struct {
		name       string
		log        string
		wantType   TransactionType
		wantAmount float64
		wantLender string
	}{
		{
			name:       "Tala loan",
			log:        "You have received Ksh5,000.00 from Tala",
			wantType:   TxnDigitalLoan,
			wantAmount: 5000.00,
			wantLender: "Tala",
		},
		{
			name:       "Branch disbursement",
			log:        "Disbursed Ksh3,000.00 from Branch to your M-Pesa",
			wantType:   TxnDigitalLoan,
			wantAmount: 3000.00,
			wantLender: "Branch",
		},
		{
			name:       "Zenka repayment",
			log:        "Ksh1,000.00 paid to Zenka successfully",
			wantType:   TxnDigitalRepay,
			wantAmount: 1000.00,
			wantLender: "Zenka",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", txn.Type, tt.wantType)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
			if txn.Lender != tt.wantLender {
				t.Errorf("Lender = %v, want %v", txn.Lender, tt.wantLender)
			}
		})
	}
}

func TestParseSingleLog_Gambling(t *testing.T) {
	tests := []struct {
		name       string
		log        string
		wantAmount float64
	}{
		{
			name:       "Betika",
			log:        "Betika: Your bet of Ksh100.00 has been placed",
			wantAmount: 100.00,
		},
		{
			name:       "SportPesa",
			log:        "SportPesa: Win! You have received Ksh500.00",
			wantAmount: 500.00,
		},
		{
			name:       "Mozzart",
			log:        "Mozzart Bet: Deposit of Ksh200.00 confirmed",
			wantAmount: 200.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, err := parseSingleLog(tt.log)
			if err != nil {
				t.Fatalf("parseSingleLog() error = %v", err)
			}
			if txn.Type != TxnGambling {
				t.Errorf("Type = %v, want %v", txn.Type, TxnGambling)
			}
			if txn.Amount != tt.wantAmount {
				t.Errorf("Amount = %v, want %v", txn.Amount, tt.wantAmount)
			}
		})
	}
}

func TestParseLogs(t *testing.T) {
	parser := NewParser()
	ctx := context.Background()

	logs := []string{
		"UA1234ABCDEF Confirmed. You have received Ksh1,500.00 from JOHN DOE 0712345678",
		"Fuliza M-PESA. You have borrowed Ksh2,000.00",
		"Invalid log message that won't match",
		"Hustler Fund. You have been disbursed Ksh500.00",
		"M-Shwari. You have deposited Ksh1,000.00 to your savings",
	}

	txns, err := parser.ParseLogs(ctx, logs)
	if err != nil {
		t.Fatalf("ParseLogs() error = %v", err)
	}

	// Should parse 4 valid transactions (skip the invalid one)
	if len(txns) != 4 {
		t.Errorf("ParseLogs() returned %d transactions, want 4", len(txns))
	}

	// Verify types
	expectedTypes := []TransactionType{
		TxnMPesaReceived,
		TxnFulizaLoan,
		TxnHustlerLoan,
		TxnMMFDeposit,
	}

	for i, expected := range expectedTypes {
		if txns[i].Type != expected {
			t.Errorf("txns[%d].Type = %v, want %v", i, txns[i].Type, expected)
		}
	}
}

func TestParseLogs_ContextCancellation(t *testing.T) {
	parser := NewParser()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	logs := make([]string, 200) // Large enough to trigger check
	for i := range logs {
		logs[i] = "UA1234ABCDEF Confirmed. You have received Ksh100.00 from TEST"
	}

	_, err := parser.ParseLogs(ctx, logs)
	if err == nil {
		t.Error("ParseLogs() should return error on cancelled context")
	}
}

func TestTransactionType_String(t *testing.T) {
	tests := []struct {
		txnType  TransactionType
		expected string
	}{
		{TxnMPesaReceived, "MPESA_RECEIVED"},
		{TxnAirtelSent, "AIRTEL_SENT"},
		{TxnHustlerLoan, "HUSTLER_LOAN"},
		{TxnOkoaReceived, "OKOA_RECEIVED"},
		{TxnMMFDeposit, "MMF_DEPOSIT"},
		{TxnDigitalLoan, "DIGITAL_LOAN"},
		{TxnBankDeposit, "BANK_DEPOSIT"},
		{TxnGambling, "GAMBLING"},
		{TxnUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.txnType.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
