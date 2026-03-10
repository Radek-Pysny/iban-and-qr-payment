package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Radek-Pysny/iban-and-qr-payment/qrpay"
	"github.com/spf13/cobra"
)

// cmdQrParse represents the qr-parse command
var cmdQrParse = &cobra.Command{
	Use:   "qr-parse",
	Short: "Parse QR code from stdin or from the given file (using zbarimg tool).",
	Long:  `...`,
	Run:   startQrParse,
}

var cfgQrParse struct {
	inputLine string

	// TODO: idea about continuous reading from stdin
}

func init() {
	cmd := cmdQrParse
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVar(
		&cfgQrParse.inputLine,
		"input-line",
		"",
		"Input line instead of reading it from stdin.",
	)
}

func startQrParse(_ *cobra.Command, _ []string) {
	cfg := cfgQrParse

	if cfg.inputLine != "" {
		err := processLineQrParse(cfg.inputLine)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return
	}

	_, err := processSingleLine(os.Stdin, processLineQrParse)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func processLineQrParse(line string) error {
	if line == "" {
		return errEmptyInput
	}

	if qrpay.IsShortPaymentDescriptor(line) {
		spd, err := qrpay.ShortPaymentDescriptorFromText(line, qrpay.ModeFlagsIgnoreFailures, nil)
		if err != nil {
			return err
		}

		fmt.Println(spd.VerboseString())
	} else if qrpay.IsShortInvoiceDescriptor(line) {
		sid, err := qrpay.ShortInvoiceDescriptorFromText(line, qrpay.ModeFlagsIgnoreFailures, nil)
		if err != nil {
			return err
		}

		fmt.Println(sid.VerboseString())
	} else {
		return errors.New("line is not SPD nor SID")
	}

	return nil
}
