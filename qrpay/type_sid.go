package qrpay

import (
	"strconv"
	"strings"
)

type ShortInvoiceDescriptor struct {
	VersionMajor uint8
	VersionMinor uint8

	// ID = IDentifier
	InvoiceID string `sid:"ID;req;-;1;40;Unique identifier of invoice"`

	// DD = ?? Date
	Date string `sid:"DD;req;date;8;8;Date"`

	// AM = AMount
	Amount string `sid:"AM;opt-req;num-fp;1;18;Amount"`

	// TP = "Typ Plneni"
	FulfilmentType string `sid:"TP;opt;dig;1;1;Fulfilment type (0 = common, 1 = RPDP, 2 = mixed)"`

	// TD = "Typ Dokladu"
	InvoiceType string `sid:"TD;opt;dig;1;1;Invoice Type (0 = proforma invoice, ...)"`

	// SA = ??
	StoreApplied string `sid:"SA;opt;dig-bool;1;1;Store (advance payment) applied"`

	// MSG = MeSsaGe
	Message string `sid:"MSG;opt;-;1;40;Subject of invoice (textual description)"`

	// ON = Order Number
	OrderNumber string `sid:"ON;opt;-;1;20;Order number (due to which was invoice created)"`

	// VS = Variable Symbol
	VariableSymbol string `sid:"VS;opt;dig;1;10;VariableSymbol"`

	// VII = Vat Identification of Invoicer
	VatIdInv string `sid:"VII;opt;-;1;14;VAT identification number of invoicer"`

	// INI = Identification Number of Invoicer
	IdNumInv string `sid:"INI;opt;dig;1;8;Identification number of invoicer"`

	// VIR = Vat Identification of Recipient
	VatIdRec string `sid:"VIR;opt;-;1;14;VAT identification number of recipient"`

	// INR = Identification Number of Recipient
	IdNumRec string `sid:"INR;opt;dig;1;8;Identification number of recipient"`

	// DUZP = "Datum Uskutecnitelneho Zdanitelneho Plneni"
	TaxableTrxDate string `sid:"DUZP;opt;date;8;8;Date of taxable transaction"`

	// DPPD = "Datum Povinnosti Priznat Dan"
	TaxDeductionDate string `sid:"DPPD;opt;date;8;8;Date of tax deduction"`

	// DT = Due daTe
	DueDate string `sid:"DT;opt;date;8;8;Date"`

	// TB0 = Tax Base - level 0 (default tax)
	TaxBase0 string `sid:"TB0;opt;num-fp;1;18;Tax base - level 0"`

	// T0 = Tax amount - level 0
	Tax0 string `sid:"TB0;opt;num-fp;1;18;Tax amount - level 0"`

	// TB1 = Tax Base - level 1 (1st reduced tax)
	TaxBase1 string `sid:"TB1;opt;num-fp;1;18;Tax base - level 1"`

	// T1 = Tax amount - level 1
	Tax1 string `sid:"TB1;opt;num-fp;1;18;Tax amount - level 1"`

	// TB2 = Tax Base - level 2 (2nd reduced tax)
	TaxBase2 string `sid:"TB2;opt;num-fp;1;18;Tax base - level 2"`

	// T2 = Tax amount - level 2
	Tax2 string `sid:"TB2;opt;num-fp;1;18;Tax amount - level 2"`

	// NTB = No-Tax Base
	NoTaxBase string `sid:"NTB;opt;num-fp;1;18;No-tax base"`

	// CC = Currency Code
	CurrencyCode string `sid:"CC;opt;abc;3;3;ISO 4217 currency code"`

	// FX = eXchange rate
	ExchangeRate string `sid:"FX;opt;num-fp;1;18;Exchange rate between CZK and CurrencyCode field"`

	// FXA = eXchange rate amount
	ExchangeRateAmount string `sid:"FXA;opt;num-fp;1;18;Count of CurrencyCode units defined for exchange rate"`

	// ACC = ACCount
	Account string `sid:"ACC;req;-;24;46;IBAN[+BIC]"`

	// CRC32 checksum
	CRC32 string `sid:"CRC32;opt;hex;8;8;Hex CRC32 control sum"`

	// X-SW
	ProprietarySoftware string `sid:"X-SW;opt;-;1;30;Software used to create QR Faktura"`

	// X-URL
	ProprietaryURL string `sid:"X-URL;opt;-;1;70;URL for invoice (structured data)"`
}

func (sid *ShortInvoiceDescriptor) VerboseString() string {
	sb := strings.Builder{}
	sb.Grow(10_000)
	sid.buildVerboseString(&sb)

	return sb.String()
}

func (sid *ShortInvoiceDescriptor) buildVerboseString(sb *strings.Builder) {
	sb.WriteString("SID v")
	sb.WriteString(strconv.Itoa(int(sid.VersionMajor)))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(int(sid.VersionMinor)))
	sb.WriteByte('\n')

	if s := sid.InvoiceID; s != "" {
		appendLine(sb, "ID (invoice ID)", s)
	}
	if s := sid.Date; s != "" {
		appendLine(sb, "DD (date)", showDate(s))
	}
	if s := sid.Amount; s != "" {
		appendLine(sb, "AM (amount)", s)
	}
	if s := sid.FulfilmentType; s != "" {
		appendLine(sb, "FT (fulfilment type)", s)
	}
	if s := sid.InvoiceType; s != "" {
		appendLine(sb, "TD (invoice type)", s)
	}
	if s := sid.StoreApplied; s != "" {
		appendLine(sb, "SA (store advance payment applied)", s)
	}
	if s := sid.Message; s != "" {
		appendLine(sb, "MSG (subject of invoice)", s)
	}
	if s := sid.OrderNumber; s != "" {
		appendLine(sb, "ON (order number)", s)
	}
	if s := sid.VariableSymbol; s != "" {
		appendLine(sb, "VS (variable symbol)", s)
	}
	if s := sid.VatIdInv; s != "" {
		appendLine(sb, "VII (VAT ID of invoicer)", s)
	}
	if s := sid.IdNumInv; s != "" {
		appendLine(sb, "INI (Identification number of invoicer)", s)
	}
	if s := sid.VatIdRec; s != "" {
		appendLine(sb, "VIR (VAT ID of recipient)", s)
	}
	if s := sid.IdNumRec; s != "" {
		appendLine(sb, "INR (Identification number of recipient)", s)
	}
	if s := sid.TaxableTrxDate; s != "" {
		appendLine(sb, "DUZP (taxable transaction date)", showDate(s))
	}
	if s := sid.TaxDeductionDate; s != "" {
		appendLine(sb, "DPPD (tax deduction date)", showDate(s))
	}
	if s := sid.DueDate; s != "" {
		appendLine(sb, "DT (due date)", showDate(s))
	}
	if s := sid.TaxBase0; s != "" {
		appendLine(sb, "TB0 (tax base - level 0)", s)
	}
	if s := sid.Tax0; s != "" {
		appendLine(sb, "T0 (tax amount - level 0)", s)
	}
	if s := sid.TaxBase1; s != "" {
		appendLine(sb, "TB1 (tax base - level 1)", s)
	}
	if s := sid.Tax1; s != "" {
		appendLine(sb, "T1 (tax amount - level 1)", s)
	}
	if s := sid.TaxBase2; s != "" {
		appendLine(sb, "TB2 (tax base - level 2)", s)
	}
	if s := sid.Tax2; s != "" {
		appendLine(sb, "T2 (tax amount - level 2)", s)
	}
	if s := sid.NoTaxBase; s != "" {
		appendLine(sb, "NTB (no-tax base)", s)
	}
	if s := sid.CurrencyCode; s != "" {
		appendLine(sb, "CC (currency code)", s)
	}
	if s := sid.ExchangeRate; s != "" {
		appendLine(sb, "FX (exchange rate)", s)
	}
	if s := sid.ExchangeRateAmount; s != "" {
		appendLine(sb, "FXA (exchange rate amount)", s)
	}
	if s := sid.Account; s != "" {
		appendLine(sb, "ACC (account)", s)
	}
	if s := sid.CRC32; s != "" {
		appendLine(sb, "CRC32", s)
	}
	if s := sid.ProprietarySoftware; s != "" {
		appendLine(sb, "X-SW (SW used to generate QR Faktura)", s)
	}
	if s := sid.ProprietaryURL; s != "" {
		appendLine(sb, "X-URL (URL for invoice or structured data)", s)
	}
}
