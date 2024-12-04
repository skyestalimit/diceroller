package diceroller

import (
	"math/rand"
	"testing"
)

// Valid Roll Args
var ValidRollArgs = []string{
	"4d4+1",
	"10d10",
	"1d6-1",
	"10000d10000-10000",
	"1D8-00",
	"1d100+0",
	"20d12-9901"}

// Valid Roll Args attribs
var ValidRollArgsAttribs = []string{
	"roll",
	"hit",
	"dmg",
	"crit",
	"spell",
	"half",
	"adv",
	"dis",
	"droplow",
	"drophigh"}

// Invalid Roll Args
var InvalidRollArgs = []string{
	"9dd9-+1",
	"2d6+123456",
	"123456d12+12",
	"12d123456-100",
	"patate1",
	"sudo reboot",
	"1d4 1d4",
	"0d2",
	"1d0",
	"1b8",
	"1+8d8+1"}

// Invalid Roll Args attribs
var InvalidRollArgsAttribs = []string{
	"bonus",
	"damidge",
	"!@#$%^&*()",
	"sudo reboot",
	"check",
	"11",
	"\"hai\""}

func TestParseValidRollArgs(t *testing.T) {
	rollArgsArray := append(ValidRollArgs, ValidRollArgsAttribs...)
	_, errors := ParseRollArgs(rollArgsArray...)
	for i := range errors {
		validArgParsingError(errors[i], t)
	}
}

func TestParseInvalidRollArgs(t *testing.T) {
	rollArgsArray := append(InvalidRollArgs, InvalidRollArgsAttribs...)
	rollExprs, _ := ParseRollArgs(rollArgsArray...)
	for e := range rollExprs {
		for i := range rollExprs[e].diceRolls {
			invalidArgParsingError(rollExprs[e].diceRolls[i].String(), t)
		}
	}
}

func TestParseSingleValidRollArg(t *testing.T) {
	for i := range ValidRollArgs {
		if _, argErr := parseRollArg(ValidRollArgs[i]); argErr != nil {
			validArgParsingError(argErr, t)
		}
	}
}

func TestParseSingleInvalidRollArg(t *testing.T) {
	for i := range InvalidRollArgs {
		if _, argErr := parseRollArg(InvalidRollArgs[i]); argErr == nil {
			invalidArgParsingError(InvalidRollArgs[i], t)
		}
	}
}

func validArgParsingError(argErr error, t *testing.T) {
	t.Fatalf("Valid dice roll returned an error: %s", argErr.Error())
}

func invalidArgParsingError(diceRollStr string, t *testing.T) {
	t.Fatalf("Invalid roll %s did not generate an error!", diceRollStr)
}

func FuzzParseRollArgs(f *testing.F) {
	f.Add("8d4-1")
	f.Fuzz(func(t *testing.T, fuzzedRollArg string) {
		ParseRollArgs(fuzzedRollArg)
	})
}

func FuzzParseManyRollArgs(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, rollArgAmmount int) {
		rollArgsArray := make([]string, 0)
		for i := 0; i < rollArgAmmount; i++ {
			rollArg := ""
			argType := rand.Intn(4) + 1

			switch argType {
			case 1:
				rollArg = ValidRollArgs[rand.Intn(len(ValidRollArgs))]
			case 2:
				rollArg = InvalidRollArgs[rand.Intn(len(InvalidRollArgs))]
			case 3:
				rollArg = ValidRollArgsAttribs[rand.Intn(len(ValidRollArgsAttribs))]
			case 4:
				rollArg = InvalidRollArgsAttribs[rand.Intn(len(InvalidRollArgsAttribs))]
			}

			rollArgsArray = append(rollArgsArray, rollArg)
		}
		ParseRollArgs(rollArgsArray...)
	})
}
