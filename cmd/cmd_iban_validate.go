package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Radek-Pysny/iban-and-qr-payment/iban"
	"github.com/spf13/cobra"
)

// cmdIbanValidate represents the iban-validate command
var cmdIbanValidate = &cobra.Command{
	Use:   "iban-validate",
	Short: "Validate IBAN code from stdin.",
	Long:  `...`,
	Run:   startIbanValidate,
}

type configIbanValidate struct {
	inputLine string
	quiet     bool
}

func (cfg configIbanValidate) stdoutPrintLn(args ...any) {
	if cfg.quiet {
		return
	}

	_, _ = fmt.Println(args...)
}

func (cfg configIbanValidate) stderrPrintLn(args ...any) {
	if cfg.quiet {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, args...)
}

var (
	cfgIbanValidate configIbanValidate

	errInvalidIBAN = errors.New("invalid IBAN")
)

func init() {
	cmd := cmdIbanValidate
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVar(
		&cfgIbanValidate.inputLine,
		"input-line",
		"",
		"Input line instead of reading it from stdin.",
	)

	cmd.Flags().BoolVarP(
		&cfgIbanValidate.quiet,
		"quiet",
		"q",
		false,
		"Suppress prints (communicate only via exit code aka $?).",
	)
}

func startIbanValidate(_ *cobra.Command, _ []string) {
	cfg := cfgIbanValidate

	if cfg.inputLine != "" {
		err := cfg.processLineIbanValidate(cfg.inputLine)
		cfg.processResult(cfg.inputLine, err)

		return
	}

	anyLine, err := processSingleLine(os.Stdin, cfg.processLineIbanValidate)
	cfg.processResult(anyLine, err)
}

func (cfg configIbanValidate) processLineIbanValidate(line string) error {
	if line == "" {
		return errEmptyInput
	}

	valid, err := iban.ValidateCheckDigits(iban.NormalizeIbanToMachineForm(line))
	switch {
	case err != nil:
		return err

	case !valid:
		return errInvalidIBAN

	default:
		return nil
	}
}

func (cfg configIbanValidate) processResult(input any, err error) {
	var inputIBAN string
	if input, ok := input.(string); ok {
		inputIBAN = `"` + iban.NormalizeIbanToHumanForm(input) + `"`
	}

	if err != nil {
		if errors.Is(err, errInvalidIBAN) {
			cfg.stdoutPrintLn(inputIBAN, "is not valid IBAN")
		} else if errors.Is(err, errInternal) {
			_, _ = fmt.Fprintln(os.Stderr, err.Error()) // ignore quiet mode
		} else {
			cfg.stderrPrintLn(inputIBAN, "is not valid IBAN:", err.Error())
		}
		os.Exit(1)
	}

	cfg.stdoutPrintLn(inputIBAN, "is valid IBAN")
}
