package qrpay_test

import (
	"errors"
	"testing"

	"github.com/Radek-Pysny/iban-and-qr-payment/qrpay"
	"github.com/stretchr/testify/require"
)

func Test_ShortPaymentDescriptorFromText_valid(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected qrpay.ShortPaymentDescriptor
	}{
		{
			title: "fixed-sconto-1",
			input: "SPD*1.0*ACC:CZ5403000000000117983223*AM:2638.00*CC:CZK*" +
				"MSG:PLATBA OBJ C.: 25548273*DT:20250825*X-VS:8250411*X-SS:82432257*",
			expected: qrpay.ShortPaymentDescriptor{
				VersionMajor:   1,
				VersionMinor:   0,
				Account:        "CZ5403000000000117983223",
				Amount:         "2638.00",
				CurrencyCode:   "CZK",
				DueTo:          "20250825",
				Message:        "PLATBA OBJ C.: 25548273",
				VariableSymbol: "8250411",
				SpecificSymbol: "82432257",
			},
		},
		{
			title: "ocp-1",
			input: "SPD*1.0*ACC:CZ3727000000001388010153*AM:50000*CC:CZK*RN:OCNI CENTRUM PRAHA, A.S.*X-VS:8605031223*" +
				"X-INV:SID%2A1.0%2AID:250771%2ADD:20251013%2AMSG:RELEXSMILE%2AINI:26169231%2AX-SW:MEDICUS%2ADT:20251031",
			expected: qrpay.ShortPaymentDescriptor{
				VersionMajor:   1,
				VersionMinor:   0,
				Account:        "CZ3727000000001388010153",
				Amount:         "50000",
				CurrencyCode:   "CZK",
				RecipientName:  "OCNI CENTRUM PRAHA, A.S.",
				VariableSymbol: "8605031223",
				QrFaktura:      "SID*1.0*ID:250771*DD:20251013*MSG:RELEXSMILE*INI:26169231*X-SW:MEDICUS*DT:20251031",
				SID: &qrpay.ShortInvoiceDescriptor{
					VersionMajor:        1,
					VersionMinor:        0,
					InvoiceID:           "250771",
					Date:                "20251013",
					Message:             "RELEXSMILE",
					IdNumInv:            "26169231",
					DueDate:             "20251031",
					ProprietarySoftware: "MEDICUS",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned, err := qrpay.ShortPaymentDescriptorFromText(tc.input, qrpay.ModeFlagsStandard, nil)

			require.NoError(t, err)
			require.Equal(t, &tc.expected, returned)
		})
	}
}

func Test_ShortPaymentDescriptorFromText_invalid(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected *qrpay.ShortPaymentDescriptorParseError
	}{
		{
			title: "sconto-1--empty-X-KS",
			input: "SPD*1.0*ACC:CZ5403000000000117983223*AM:2638.00*CC:CZK*" +
				"MSG:PLATBA OBJ C.: 25548273*DT:20250825*X-VS:8250411*X-KS:*X-SS:82432257*",
			expected: qrpay.NewParseError("X-KS", "", "length check", errors.New("value is shorter than 1")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := qrpay.ShortPaymentDescriptorFromText(tc.input, qrpay.ModeFlagsStandard, nil)

			require.NotNil(t, err)

			var qrpayErr *qrpay.ShortPaymentDescriptorParseError
			require.True(t, errors.As(err, &qrpayErr))
			errors.As(err, &qrpayErr)

			require.True(t, tc.expected.Equal(qrpayErr))
		})
	}
}
