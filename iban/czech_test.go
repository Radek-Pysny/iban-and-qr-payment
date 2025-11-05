package iban_test

import (
	"testing"

	"github.com/Radek-Pysny/iban-and-qr-payment/iban"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	accountNumberPrefix string
	accountNumberBase   string
	bankCode            string
	checkDigits         string
	bankName            string
	humanReadable       string
	machineReadable     string
	czechNational       string
	fullCzechNational   string
}

func Test_CzechIbanFromText(t *testing.T) {

	testCases := []struct {
		title         string
		input         string
		expectedError error
		expected      expected
	}{
		{
			title:         "wrong-length--empty",
			input:         "",
			expectedError: iban.ErrWrongLength,
		},
		{
			title:         "wrong-length--10",
			input:         "01234567890",
			expectedError: iban.ErrWrongLength,
		},
		{
			title:         "wrong-length--16",
			input:         "01234567890ABCDEF",
			expectedError: iban.ErrWrongLength,
		},
		{
			title:         "wrong-length--23--one-less",
			input:         "CZ540300000000011798322",
			expectedError: iban.ErrWrongLength,
		},
		{
			title:         "wrong-length--23--one-more",
			input:         "CZ5403000000000117983223Z",
			expectedError: iban.ErrWrongLength,
		},
		{
			title:         "wrong-prefix",
			input:         "DE5403000000000117983223",
			expectedError: iban.ErrWrongPrefix,
		},
		{
			title:         "wrong-check-digit",
			input:         "CZ9903000000000117983223",
			expectedError: iban.ErrWrongCheckDigits,
		},
		{
			title: "sconto-machine",
			input: "CZ5403000000000117983223",
			expected: expected{
				accountNumberPrefix: "000000",
				accountNumberBase:   "0117983223",
				bankCode:            "0300",
				checkDigits:         "54",
				bankName:            "ČSOB, a.s.",
				humanReadable:       "CZ54 0300 0000 0001 1798 3223",
				machineReadable:     "CZ5403000000000117983223",
				czechNational:       "117983223/0300",
				fullCzechNational:   "000000-0117983223/0300",
			},
		},
		{
			title: "sconto-human",
			input: "CZ54 0300 0000 0001 1798 3223",
			expected: expected{
				accountNumberPrefix: "000000",
				accountNumberBase:   "0117983223",
				bankCode:            "0300",
				checkDigits:         "54",
				bankName:            "ČSOB, a.s.",
				humanReadable:       "CZ54 0300 0000 0001 1798 3223",
				machineReadable:     "CZ5403000000000117983223",
				czechNational:       "117983223/0300",
				fullCzechNational:   "000000-0117983223/0300",
			},
		},
		{
			title: "cnb-sample-1-machine",
			input: "CZ6508000000192000145399",
			expected: expected{
				accountNumberPrefix: "000019",
				accountNumberBase:   "2000145399",
				bankCode:            "0800",
				bankName:            "Česká spořitelna, a.s.",
				checkDigits:         "65",
				humanReadable:       "CZ65 0800 0000 1920 0014 5399",
				machineReadable:     "CZ6508000000192000145399",
				czechNational:       "19-2000145399/0800",
				fullCzechNational:   "000019-2000145399/0800",
			},
		},
		{
			title: "cnb-sample-1-human",
			input: "CZ65 0800 0000 1920 0014 5399",
			expected: expected{
				accountNumberPrefix: "000019",
				accountNumberBase:   "2000145399",
				bankCode:            "0800",
				checkDigits:         "65",
				bankName:            "Česká spořitelna, a.s.",
				humanReadable:       "CZ65 0800 0000 1920 0014 5399",
				machineReadable:     "CZ6508000000192000145399",
				czechNational:       "19-2000145399/0800",
				fullCzechNational:   "000019-2000145399/0800",
			},
		},
		{
			title: "cnb-sample-2-machine",
			input: "CZ6907101781240000004159",
			expected: expected{
				accountNumberPrefix: "178124",
				accountNumberBase:   "0000004159",
				bankCode:            "0710",
				checkDigits:         "69",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ69 0710 1781 2400 0000 4159",
				machineReadable:     "CZ6907101781240000004159",
				czechNational:       "178124-4159/0710",
				fullCzechNational:   "178124-0000004159/0710",
			},
		},
		{
			title: "cnb-sample-2-human",
			input: "CZ69 0710 1781 2400 0000 4159",
			expected: expected{
				accountNumberPrefix: "178124",
				accountNumberBase:   "0000004159",
				bankCode:            "0710",
				checkDigits:         "69",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ69 0710 1781 2400 0000 4159",
				machineReadable:     "CZ6907101781240000004159",
				czechNational:       "178124-4159/0710",
				fullCzechNational:   "178124-0000004159/0710",
			},
		},
		{
			title: "ocp-machine",
			input: "CZ3727000000001388010153",
			expected: expected{
				accountNumberPrefix: "000000",
				accountNumberBase:   "1388010153",
				bankCode:            "2700",
				checkDigits:         "37",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ37 2700 0000 0013 8801 0153",
				machineReadable:     "CZ3727000000001388010153",
				czechNational:       "1388010153/2700",
				fullCzechNational:   "000000-1388010153/2700",
			},
		},
		{
			title: "ocp-human",
			input: "CZ37 2700 0000 0013 8801 0153",
			expected: expected{
				accountNumberPrefix: "000000",
				accountNumberBase:   "1388010153",
				bankCode:            "2700",
				checkDigits:         "37",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ37 2700 0000 0013 8801 0153",
				machineReadable:     "CZ3727000000001388010153",
				czechNational:       "1388010153/2700",
				fullCzechNational:   "000000-1388010153/2700",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned, err := iban.CzechIbanFromText(tc.input)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			assertExpectedCzechIBAN(t, &tc.expected, returned)
		})
	}
}

func Test_CzechIbanFromNationalFormat(t *testing.T) {
	testCases := []struct {
		title         string
		input         string
		expectedError error
		expected      expected
	}{
		{
			title:         "no-slash",
			input:         "1900023346",
			expectedError: iban.ErrCzNoSepBetweenAccountNumberAndBankCode,
		},
		{
			title:         "too-long-prefix",
			input:         "1234567-01/0100",
			expectedError: iban.ErrCzTooLargeAccountNumberPrefix,
		},
		{
			title:         "too-long-base-1",
			input:         "12345678901/0100",
			expectedError: iban.ErrCzTooLargeAccountNumberBase,
		},
		{
			title:         "too-long-base-2",
			input:         "123456-12345678901/0100",
			expectedError: iban.ErrCzTooLargeAccountNumberBase,
		},
		{
			title:         "zero-base-1",
			input:         "0/0100",
			expectedError: iban.ErrCzZeroAccountNumberBase,
		},
		{
			title:         "zero-base-2",
			input:         "123456-0/0100",
			expectedError: iban.ErrCzZeroAccountNumberBase,
		},
		{
			title:         "zero-bank-code-1",
			input:         "1234567890/0",
			expectedError: iban.ErrCzZeroBankCode,
		},
		{
			title:         "zero-bank-code-2",
			input:         "123456-1234567890/0",
			expectedError: iban.ErrCzZeroBankCode,
		},
		{
			title:         "too-long-bank-code-1",
			input:         "1234567890/12345",
			expectedError: iban.ErrCzTooLargeBankCode,
		},
		{
			title:         "too-long-bank-code-2",
			input:         "123456-1234567890/12345",
			expectedError: iban.ErrCzTooLargeBankCode,
		},
		{
			title: "sconto",
			input: "0117983223/0300",
			expected: expected{
				accountNumberPrefix: "000000",
				accountNumberBase:   "0117983223",
				bankCode:            "0300",
				checkDigits:         "54",
				bankName:            "ČSOB, a.s.",
				humanReadable:       "CZ54 0300 0000 0001 1798 3223",
				machineReadable:     "CZ5403000000000117983223",
				czechNational:       "117983223/0300",
				fullCzechNational:   "000000-0117983223/0300",
			},
		},
		{
			title: "cnb-sample-1",
			input: "19-2000145399/0800",
			expected: expected{
				accountNumberPrefix: "000019",
				accountNumberBase:   "2000145399",
				bankCode:            "0800",
				bankName:            "Česká spořitelna, a.s.",
				checkDigits:         "65",
				humanReadable:       "CZ65 0800 0000 1920 0014 5399",
				machineReadable:     "CZ6508000000192000145399",
				czechNational:       "19-2000145399/0800",
				fullCzechNational:   "000019-2000145399/0800",
			},
		},
		{
			title: "cnb-sample-2",
			input: "178124-4159/0710",
			expected: expected{
				accountNumberPrefix: "178124",
				accountNumberBase:   "0000004159",
				bankCode:            "0710",
				checkDigits:         "69",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ69 0710 1781 2400 0000 4159",
				machineReadable:     "CZ6907101781240000004159",
				czechNational:       "178124-4159/0710",
				fullCzechNational:   "178124-0000004159/0710",
			},
		},
		{
			title: "cnb-sample-2",
			input: "178124-4159/0710",
			expected: expected{
				accountNumberPrefix: "178124",
				accountNumberBase:   "0000004159",
				bankCode:            "0710",
				checkDigits:         "69",
				bankName:            "Česká národní banka",
				humanReadable:       "CZ69 0710 1781 2400 0000 4159",
				machineReadable:     "CZ6907101781240000004159",
				czechNational:       "178124-4159/0710",
				fullCzechNational:   "178124-0000004159/0710",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned, err := iban.CzechIbanFromNationalFormat(tc.input)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			assertExpectedCzechIBAN(t, &tc.expected, returned)
		})
	}
}

func assertExpectedCzechIBAN(t *testing.T, expected *expected, returned iban.CzechIBAN) {
	t.Helper()

	assert.Equal(t, "CZ", returned.CountryCode())
	assert.Equal(t, expected.accountNumberPrefix, returned.AccountNumberPrefix())
	assert.Equal(t, expected.accountNumberBase, returned.AccountNumberBase())
	assert.Equal(t, expected.bankCode, returned.BankCode())
	assert.Equal(t, expected.checkDigits, returned.CheckDigits())

	assert.Equal(t, expected.humanReadable, returned.HumanFormString())
	assert.Equal(t, expected.machineReadable, returned.MachineFormString())
	assert.Equal(t, expected.fullCzechNational, returned.FullCzechNationalString())
	assert.Equal(t, expected.czechNational, returned.CzechNationalString())
}
