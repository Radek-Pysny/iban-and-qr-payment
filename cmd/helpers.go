package cmd

import (
	"bufio"
	"io"
	"strings"
)

func processSingleLine(reader io.Reader, fn func(string) error) (any, error) {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return nil, errReadLn
	}

	line := strings.TrimSpace(scanner.Text())
	err := fn(line)
	if err != nil {
		return nil, err
	}

	return line, nil
}
