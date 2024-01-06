package diceroller

import (
	"testing"
)

// Valid Roll Args
var validRollArgs = []string{
	"4d4+1",
	"10d10",
	"1d6-1",
	"10000d10000-10000",
	"1D8-00"}

// Invalid Roll Args
var invalidRollArgs = []string{
	"9dd9-+1",
	"2d6+123456",
	"123456d12+12",
	"12d123456-100",
	"patate1",
	"sudo reboot",
	"1d4 1d4",
	"0d1",
	"1d0",
	"1b8",
	"1+8d8+1"}

func TestParseValidRollArgs(t *testing.T) {
	_, errors := ParseRollArgs(validRollArgs)
	for i := range errors {
		validArgParsingError(errors[i], t)
	}
}

func TestParseInvalidRollArgs(t *testing.T) {
	diceRolls, _ := ParseRollArgs(invalidRollArgs)
	for i := range diceRolls {
		invalidArgParsingError(diceRolls[i].String(), t)
	}
}

func TestParseSingleValidRollArg(t *testing.T) {
	for i := range validRollArgs {
		if _, argErr := ParseRollArg(validRollArgs[i]); argErr != nil {
			validArgParsingError(argErr, t)
		}
	}
}

func TestParseSingleInvalidRollArg(t *testing.T) {
	for i := range invalidRollArgs {
		if _, argErr := ParseRollArg(invalidRollArgs[i]); argErr == nil {
			invalidArgParsingError(invalidRollArgs[i], t)
		}
	}
}

func validArgParsingError(argErr error, t *testing.T) {
	t.Fatalf("Valid dice roll returned an error: %s", argErr.Error())
}

func invalidArgParsingError(diceRollStr string, t *testing.T) {
	t.Fatalf("Invalid roll %s did not generate an error!", diceRollStr)
}
