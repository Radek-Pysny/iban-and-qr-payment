package qrpay

import (
	"strconv"
	"strings"
)

type ShortPaymentDescriptor struct {
	VersionMajor uint8
	VersionMinor uint8

	Account             string   `spd:"ACC;req;-;24;46;IBAN[+BIC]"`
	AlternativeAccounts []string `spd:"ALT-ACC;opt;-;24;93;IBAN+[BIC[,IBAN+[BIC]]]"`
	Amount              string   `spd:"AM;opt;num-fp;1;10;Amount"`
	CurrencyCode        string   `spd:"CC;opt;abc;3;3;ISO 4217 currency code"`
	SenderReference     string   `spd:"RF;opt;dig;1;16;A sender reference number"`
	RecipientName       string   `spd:"RN;opt;-;1;35;Recipient name"`
	DueTo               string   `spd:"DT;opt;date;8;8;YYYYMMDD due date"`
	PaymentType         string   `spd:"PT;opt;-;1;3;Payment type code"`
	Message             string   `spd:"MSG;opt;-;1;60;Message for recipient"`
	CRC32               string   `spd:"CRC32;opt;hex;8;8;Hex CRC32 control sum"`
	NotificationChannel string   `spd:"NT;opt;-;1;1;Notification channel (P = SMS (on Phone), E = Email)"`
	NotificationAddress string   `spd:"NTA;opt;-;1;320;Notification address"`

	// Common Czech extensions
	PaymentRetryPeriod string `spd:"X-PER;opt;dig;1;2;Count of days for retrials (max. 30)"`
	VariableSymbol     string `spd:"X-VS;opt;dig;1;10;Variable symbol (numeric only)"`
	SpecificSymbol     string `spd:"X-SS;opt;dig;1;10;Specific symbol (numeric only)"`
	ConstantSymbol     string `spd:"X-KS;opt;dig;1;10;Constant symbol (numeric only)"`
	InternalID         string `spd:"X-ID;opt;-;1;20;Internal ID"`
	URL                string `spd:"X-URL;opt;-;1;140;Custom URL"`

	// Czech ``QR Faktura'' extension
	QrFaktura string `spd:"X-INV;opt;ast;1;538;Czech-specific QR invoice encoded fields"`
	SID       *ShortInvoiceDescriptor
}

func (spd *ShortPaymentDescriptor) VerboseString() string {
	sb := strings.Builder{}
	sb.Grow(10_000)

	sb.WriteString("SPD v")
	sb.WriteString(strconv.Itoa(int(spd.VersionMajor)))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(int(spd.VersionMinor)))
	sb.WriteByte('\n')

	if s := spd.Account; s != "" {
		appendLine(&sb, "ACC (account)", s)
	}
	if len(spd.AlternativeAccounts) > 0 {
		appendLine(&sb, "ALT-ACC (alternative accounts)", strings.Join(spd.AlternativeAccounts, ", "))
	}
	if s := spd.Amount; s != "" {
		appendLine(&sb, "AM (amount)", s)
	}
	if s := spd.CurrencyCode; s != "" {
		appendLine(&sb, "CC (currency code)", s)
	}
	if s := spd.SenderReference; s != "" {
		appendLine(&sb, "RF (sender reference)", s)
	}
	if s := spd.RecipientName; s != "" {
		appendLine(&sb, "RN (recipient name)", s)
	}
	if s := spd.RecipientName; s != "" {
		appendLine(&sb, "RN (recipient name)", s)
	}
	if s := spd.DueTo; s != "" {
		appendLine(&sb, "DT (due to date)", showDate(s))
	}
	if s := spd.PaymentType; s != "" {
		appendLine(&sb, "PT (payment type)", s)
	}
	if s := spd.Message; s != "" {
		appendLine(&sb, "MSG (message)", s)
	}
	if s := spd.CRC32; s != "" {
		appendLine(&sb, "CRC32", s)
	}
	if s := spd.NotificationChannel; s != "" {
		appendLine(&sb, "NT (notification channel)", s)
	}
	if s := spd.NotificationAddress; s != "" {
		appendLine(&sb, "NTA (notification address)", s)
	}

	// Common Czech extensions
	if s := spd.PaymentRetryPeriod; s != "" {
		appendLine(&sb, "X-PER (retry period)", s)
	}
	if s := spd.VariableSymbol; s != "" {
		appendLine(&sb, "X-VS (variable symbol)", s)
	}
	if s := spd.SpecificSymbol; s != "" {
		appendLine(&sb, "X-SS (specific symbol)", s)
	}
	if s := spd.ConstantSymbol; s != "" {
		appendLine(&sb, "X-KS (constant)", s)
	}
	if s := spd.InternalID; s != "" {
		appendLine(&sb, "X-ID (internal ID)", s)
	}
	if s := spd.URL; s != "" {
		appendLine(&sb, "X-URL (custom URL)", s)
	}

	if sid := spd.SID; sid != nil {
		sb.WriteByte('\n')
		appendLine(&sb, "X-INV (invoice)", "...")
		sid.buildVerboseString(&sb)
	}

	return sb.String()
}
