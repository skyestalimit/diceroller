package diceroller

import (
	"math/rand"
	"testing"
)

// Valid Roll Args
var validRollArgs = []string{
	"4d4+1",
	"10d10",
	"1d6-1",
	"10000d10000-10000",
	"1D8-00",
	"1d100+0",
	"20d12-9901"}

// Valid Roll Args attribs
var validRollArgsAttribs = []string{
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
var invalidRollArgs = []string{
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
var invalidRollArgsAttribs = []string{
	"bonus",
	"damidge",
	"!@#$%^&*()",
	"sudo reboot",
	"check",
	"11",
	"\"hai\""}

func TestParseValidRollArgs(t *testing.T) {
	rollArgsArray := append(validRollArgs, validRollArgsAttribs...)
	_, errors := parseRollArgs(rollArgsArray...)
	for i := range errors {
		validArgParsingError(errors[i], t)
	}
}

func TestParseInvalidRollArgs(t *testing.T) {
	rollArgsArray := append(invalidRollArgs, invalidRollArgsAttribs...)
	rollExprs, _ := parseRollArgs(rollArgsArray...)
	for e := range rollExprs {
		for i := range rollExprs[e].diceRolls {
			invalidArgParsingError(rollExprs[e].diceRolls[i].String(), t)
		}
	}
}

func TestParseSingleValidRollArg(t *testing.T) {
	for i := range validRollArgs {
		if _, argErr := parseRollArg(validRollArgs[i]); argErr != nil {
			validArgParsingError(argErr, t)
		}
	}
}

func TestParseSingleInvalidRollArg(t *testing.T) {
	for i := range invalidRollArgs {
		if _, argErr := parseRollArg(invalidRollArgs[i]); argErr == nil {
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

func FuzzParseRollArgs(f *testing.F) {
	f.Add("8d4-1")
	f.Fuzz(func(t *testing.T, fuzzedRollArg string) {
		parseRollArgs(fuzzedRollArg)
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
				rollArg = validRollArgs[rand.Intn(len(validRollArgs))]
			case 2:
				rollArg = invalidRollArgs[rand.Intn(len(invalidRollArgs))]
			case 3:
				rollArg = validRollArgsAttribs[rand.Intn(len(validRollArgsAttribs))]
			case 4:
				rollArg = invalidRollArgsAttribs[rand.Intn(len(invalidRollArgsAttribs))]
			}

			rollArgsArray = append(rollArgsArray, rollArg)
		}
		parseRollArgs(rollArgsArray...)
	})
}
