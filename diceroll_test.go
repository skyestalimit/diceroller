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
		*newDiceRoll(4, 4, 1, true),
	},
	{
		`10d10`,
		regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
		*newDiceRoll(10, 10, 0, true),
	},
	{
		`1d6-1`,
		regexp.MustCompile(`\[[1-6]\]`),
		*newDiceRoll(1, 6, -1, true),
	},
	{
		`4d12-12345`,
		regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
		*newDiceRoll(4, 12, -12345, true),
	},
	{
		`-5d6-1`,
		regexp.MustCompile(`\[[1-6] [1-6] [1-6] [1-6] [1-6]\]`),
		*newDiceRoll(5, 6, -1, false),
	},
}

// Invalid Dice Rolls
var invalidDiceRollsValues = []diceRollTestValues{
	{
		`123456d10`,
		nil,
		DiceRoll{123456, 10, 0, true, newDnDRollAttributes()},
	},
	{
		`10d123456+10`,
		nil,
		DiceRoll{10, 123456, 10, true, newDnDRollAttributes()},
	},
	{
		`10d10-123456`,
		nil,
		DiceRoll{10, 10, -123456, true, newDnDRollAttributes()},
	},
	{
		`0d8`,
		nil,
		DiceRoll{0, 8, 0, true, newDnDRollAttributes()},
	},
	{
		`1d0`,
		nil,
		DiceRoll{1, 0, 0, true, newDnDRollAttributes()},
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
	f.Add(2, 8, 1, true)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int, plus bool) {
		NewDiceRollWithAttribs(diceAmmount, diceSize, modifier, plus, nil)
	})
}
