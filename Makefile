.PHONY: unit-tests unit-tests-verbose

unit-tests:
	go test ./... -count=1

unit-tests-verbose:
	go test ./... -count=1 -v
