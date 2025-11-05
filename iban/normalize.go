package iban

import (
	"strings"
	"unicode"

	"github.com/Radek-Pysny/iban-and-qr-payment/stringutils"
)

func NormalizeIbanToMachineForm(iban string) string {
	sb := strings.Builder{}
	sb.Grow(len(iban))

	for _, ch := range iban {
		if unicode.IsSpace(ch) {
			continue
		}

		sb.WriteRune(ch)
	}

	return sb.String()
}

func NormalizeIbanToHumanForm(iban string) string {
	return stringutils.SeparateGroups(NormalizeIbanToMachineForm(iban), 4, " ")
}
