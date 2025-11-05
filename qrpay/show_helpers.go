package qrpay

import (
	"strings"
)

func appendLine(sb *strings.Builder, description string, value string) {
	sb.WriteByte(' ')
	sb.WriteString(description)
	sb.WriteByte(':')
	sb.WriteByte(' ')
	sb.WriteString(value)
	sb.WriteByte('\n')
}

func showDate(text string) string {
	if len(text) == 8 {
		return text[:4] + "-" + text[4:6] + "-" + text[6:]
	}

	return text
}
