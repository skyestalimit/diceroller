package roller

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

// Test basic init and string representation
func TestValidDiceRollString(t *testing.T) {
	for i := range validDiceRollsValues {
		validateDiceRollString(validDiceRollsValues[i], t)
	}
}

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

// Test invalid DiceRoll results
func TestInvalidDiceRollResult(t *testing.T) {
	for i := range invalidDiceRollsValues {
		validateInvalidDiceRoll(invalidDiceRollsValues[i], t)
	}
}

// Test rolling an array of valid DiceRoll
func TestPerformValidRolls(t *testing.T) {
	// Build valid DiceRoll array
	validDiceRolls := make([]DiceRoll, 0)
	for i := range validDiceRollsValues {
		validDiceRolls = append(validDiceRolls, validDiceRollsValues[i].diceRoll)
	}

	// Perform roll on DiceRoll array
	if results, diceErrs := PerformRolls(validDiceRolls); len(diceErrs) == 0 {
		if len(results) == len(validDiceRolls) {
			for i := range results {
				validateDiceRollResult(results[i], validDiceRollsValues[i].resultFormat, t)
			}
		} else {
			// Missing results, fail the test
			t.Fatalf("Result list length = %d, want match for %d", len(results), len(validDiceRolls))
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
	// Build valid DiceRoll array
	invalidDiceRolls := make([]DiceRoll, 0)
	for i := range invalidDiceRollsValues {
		invalidDiceRolls = append(invalidDiceRolls, invalidDiceRollsValues[i].diceRoll)
	}

	// Perform roll on DiceRoll array
	if _, diceErrs := PerformRolls(invalidDiceRolls); len(diceErrs) > 0 {
		if len(diceErrs) != len(invalidDiceRolls) {
			// Missing results, fail the test
			t.Fatalf("Error list length = %d, want match for %d", len(diceErrs), len(invalidDiceRolls))
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

func validateInvalidDiceRoll(diceValues diceRollTestValues, t *testing.T) {
	if _, diceErr := diceValues.diceRoll.PerformRoll(); diceErr == nil {
		t.Fatalf("Invalid dice roll did not generate an error: %s", diceValues.wantedDiceStr)
	}
}
