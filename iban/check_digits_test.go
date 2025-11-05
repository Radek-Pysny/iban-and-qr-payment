package iban_test

import (
	"strconv"
	"testing"

	"github.com/Radek-Pysny/iban-and-qr-payment/iban"
	"github.com/stretchr/testify/require"
)

func Test_IterateIbanForCheckDigits(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "GB82WEST12345698765432",
			expected: "3214282912345698765432161182",
		},
		{
			input:    "CH750034633221115556T",
			expected: "003463322111555629121775",
		},
		{
			input:    "CZ3727000000001388010153",
			expected: "27000000001388010153123537",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			acc := ""
			for char, ok := range iban.IterateIbanForCheckDigits(tc.input) {
				require.True(t, ok)

				acc += string(char + '0')
			}

			require.Equal(t, tc.expected, acc)
		})
	}
}

func Test_ValidateCheckDigits(t *testing.T) {
	testCases := []struct {
		input       string
		expected    bool
		expectedErr error
	}{
		{
			input:    "GB82WEST12345698765432",
			expected: true,
		},
		{
			input:    "IE64IRCE92050112345678",
			expected: true,
		},
		{
			input:    "BI1320001100010000123456789",
			expected: true,
		},
		{
			input:    "CZ3727000000001388010153",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			returned, err := iban.ValidateCheckDigits(tc.input)

			if tc.expectedErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, returned)
		})
	}
}

func Test_CalculateCheckDigitsRemainder(t *testing.T) {
	testCases := []struct {
		input       string
		expected    uint8
		expectedErr error
	}{
		{
			input:    "GB82WEST12345698765432",
			expected: 1,
		},
		{
			input:    "IE64IRCE92050112345678",
			expected: 1,
		},
		{
			input:    "BI1320001100010000123456789",
			expected: 1,
		},
		{
			input:    "CZ3727000000001388010153",
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			returned, err := iban.CalculateCheckDigitsRemainder(tc.input)

			if tc.expectedErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, returned)
		})
	}
}

func Test_IterateIbanCheckDigitsReminderStepByStep(t *testing.T) {
	testCases := []struct {
		input    string
		expected []uint32
	}{
		{
			input:    "GB82WEST12345698765432",
			expected: []uint32{321_428_291, 702_345_698, 297_654_321, 2_461_182, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			require.NotZero(t, len(tc.expected))

			var (
				index    = 0
				returned uint32
				err      error
			)
			for returned, err = range iban.IterateIbanCheckDigitsReminderStepByStep(tc.input) {
				require.NoError(t, err, "%d: error", index)
				require.Less(t, index, len(tc.expected), "%d: too many results", index)
				require.Equal(
					t,
					strconv.FormatUint(uint64(tc.expected[index]), 10),
					strconv.FormatUint(uint64(returned), 10),
					"%d: result",
					index,
				)

				index++
			}

			require.Equal(t, len(tc.expected), index)
			require.Equal(t, tc.expected[len(tc.expected)-1], returned)
		})
	}
}
