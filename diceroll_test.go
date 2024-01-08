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
var validDiceRollsValues = []diceRollTestValues{values4d4Plus1, values10d10, values1d6Minus1}

var values4d4Plus1 diceRollTestValues = diceRollTestValues{
	`4d4+1`,
	regexp.MustCompile(`\[[0-9] [0-9] [0-9] [0-9]\]`),
	DiceRoll{4, 4, 1}}
var values10d10 diceRollTestValues = diceRollTestValues{
	`10d10`,
	regexp.MustCompile(`\[[0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+ [0-9]+\]`),
	DiceRoll{10, 10, 0}}
var values1d6Minus1 diceRollTestValues = diceRollTestValues{
	`1d6-1`,
	regexp.MustCompile(`\[[0-9]\]`),
	DiceRoll{1, 6, -1}}

// Invalid Dice Rolls
var invalidDiceRollsValues = []diceRollTestValues{valuesBigAmmount, valuesBigSize, valuesBigModifier, valuesZeroAmmount, valuesZeroSize}

var valuesBigAmmount diceRollTestValues = diceRollTestValues{
	`123456d10`,
	nil,
	DiceRoll{123456, 10, 0}}
var valuesBigSize diceRollTestValues = diceRollTestValues{
	`10d123456+10`,
	nil,
	DiceRoll{10, 123456, 10}}
var valuesBigModifier diceRollTestValues = diceRollTestValues{
	`10d10-123456`,
	nil,
	DiceRoll{10, 10, -123456}}
var valuesZeroAmmount diceRollTestValues = diceRollTestValues{
	`0d8`,
	nil,
	DiceRoll{0, 8, 0}}
var valuesZeroSize diceRollTestValues = diceRollTestValues{
	`1d0`,
	nil,
	DiceRoll{1, 0, 0}}

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

// Test valid DiceRoll results
func TestValidDiceRollResult(t *testing.T) {
	for i := range validDiceRollsValues {
		rollAndvalidateDiceRollResult(validDiceRollsValues[i], t)
	}
}

// Test invalid DiceRolls
func TestInvalidDiceRollResult(t *testing.T) {
	for i := range invalidDiceRollsValues {
		validateInvalidDiceRoll(invalidDiceRollsValues[i], t)
	}
}

// Test rolling an array of valid DiceRoll
func TestPerformValidRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	if results, diceErrs := PerformRolls(diceRollsFromTestValues(validDiceRollsValues)); len(diceErrs) == 0 {
		if len(results) == len(validDiceRollsValues) {
			for i := range results {
				validateDiceRollResult(results[i], validDiceRollsValues[i].resultFormat, t)
			}
		} else {
			// Missing results, fail the test
			t.Fatalf("Result list length = %d, want match for %d", len(results), len(validDiceRollsValues))
		}
	} else {
		// Received errors, fail the test
		errStr := ""
		for i := range diceErrs {
			errStr += diceErrs[i].Error()
		}
		t.Fatalf("Unexpected dice roll errors: %s", errStr)
	}
}

// Test rolling an array of invalid DiceRoll
func TestPerformInvalidRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	if _, diceErrs := PerformRolls(diceRollsFromTestValues(invalidDiceRollsValues)); len(diceErrs) > 0 {
		if len(diceErrs) != len(invalidDiceRollsValues) {
			// Missing errors, fail the test
			t.Fatalf("Error list length = %d, want match for %d", len(diceErrs), len(invalidDiceRollsValues))
		}
	}
}

// Validate DiceRoll string format
func validateDiceRollString(diceValues diceRollTestValues, t *testing.T) {
	if diceStr := diceValues.diceRoll.String(); !strings.EqualFold(diceValues.wantedDiceStr, diceStr) {
		t.Fatalf("DiceRoll = %s, must equal %s", diceStr, diceValues.wantedDiceStr)
	}
}

// Validate DiceRollResult matches expected format
func rollAndvalidateDiceRollResult(diceValues diceRollTestValues, t *testing.T) {
	if result, diceErr := diceValues.diceRoll.PerformRoll(); diceErr == nil {
		validateDiceRollResult(*result, diceValues.resultFormat, t)
	} else {
		t.Fatalf("Unexpected dice roll error: %s", diceErr.Error())
	}
}

// Validate roll result matches expected format
func validateDiceRollResult(result DiceRollResult, format *regexp.Regexp, t *testing.T) {
	if resultStr := result.String(); !format.MatchString(result.String()) {
		t.Fatalf("Roll result = %s, want match for %#q", resultStr, format)
	}
}

// Validate invalid DiceRoll generates an error
func validateInvalidDiceRoll(diceValues diceRollTestValues, t *testing.T) {
	if _, diceErr := diceValues.diceRoll.PerformRoll(); diceErr == nil {
		t.Fatalf("Invalid dice roll did not generate an error: %s", diceValues.wantedDiceStr)
	}
}

// Extracts and returns a DiceRoll array from a diceRollTestValues array
func diceRollsFromTestValues(testValues []diceRollTestValues) (diceRolls []DiceRoll) {
	for i := range testValues {
		diceRolls = append(diceRolls, testValues[i].diceRoll)
	}
	return
}
