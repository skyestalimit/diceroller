package diceroller

import (
	"regexp"
	"strings"
	"testing"
)

type diceRollTestValues struct {
	wantedDiceStr string
	resultFormat  *regexp.Regexp
	diceRoll      DiceRoll
}

// Valid Dice Rolls
var validDiceRollsValues = []diceRollTestValues{
	{
		`4d4+1`,
		regexp.MustCompile(`\[[1-4] [1-4] [1-4] [1-4]\]`),
		DiceRoll{4, 4, 1, newRollAttributes()},
	},
	{
		`10d10`,
		regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
		DiceRoll{10, 10, 0, newRollAttributes()},
	},
	{
		`1d6-1`,
		regexp.MustCompile(`\[[1-6]\]`),
		DiceRoll{1, 6, -1, newRollAttributes()},
	},
	{
		`4d12-12345`,
		regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
		DiceRoll{4, 12, -12345, newRollAttributes()},
	},
	{
		`-5d6-1`,
		regexp.MustCompile(`\[[1-6] [1-6] [1-6] [1-6] [1-6]\]`),
		DiceRoll{5, 6, -1, newRollAttributes(minusAttrib)},
	},
	{
		`1d2-4`,
		regexp.MustCompile(`\[[1-2]\]`),
		DiceRoll{1, 2, -4, newRollAttributes()},
	},
}

// Invalid Dice Rolls
var invalidDiceRollsValues = []diceRollTestValues{
	{
		`123456d10`,
		nil,
		DiceRoll{123456, 10, 0, newRollAttributes()},
	},
	{
		`10d123456+10`,
		nil,
		DiceRoll{10, 123456, 10, newRollAttributes()},
	},
	{
		`10d10-123456`,
		nil,
		DiceRoll{10, 10, -123456, newRollAttributes()},
	},
	{
		`0d8`,
		nil,
		DiceRoll{0, 8, 0, newRollAttributes()},
	},
	{
		`-1d0`,
		nil,
		DiceRoll{1, 0, 0, newRollAttributes(minusAttrib)},
	},
}

// Test basic init and string representation for valid DiceRolls
func TestValidDiceRollString(t *testing.T) {
	for i := range validDiceRollsValues {
		validateDiceRollString(validDiceRollsValues[i], t)
	}
}

// Test basic init and string representation for invalid DiceRolls
func TestInvalidDiceRollString(t *testing.T) {
	for i := range invalidDiceRollsValues {
		validateDiceRollString(invalidDiceRollsValues[i], t)
	}
}

// Test basic init and string representation for invalid DiceRolls
func TestValidRoll(t *testing.T) {
	for i := range validDiceRollsValues {
		diceRoll := validDiceRollsValues[i].diceRoll
		roll := diceRoll.Roll()
		if roll == 0 {
			t.Fatalf("Valid DiceRoll %s rolled %d, wanted > 1", diceRoll.String(), roll)
		}
	}
}

// Validates DiceRoll string format
func validateDiceRollString(diceValues diceRollTestValues, t *testing.T) {
	if diceStr := diceValues.diceRoll.String(); !strings.EqualFold(diceValues.wantedDiceStr, diceStr) {
		t.Fatalf("DiceRoll = %s, must equal %s", diceStr, diceValues.wantedDiceStr)
	}
}

func FuzzNewDiceRoll(f *testing.F) {
	f.Add(2, 8, 1)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int) {
		NewDiceRollWithAttribs(diceAmmount, diceSize, modifier, nil)
	})
}
