package roller

import (
	"regexp"
	"strings"
	"testing"
)

type diceRollTestValues struct {
	asString     string
	resultFormat *regexp.Regexp
	diceRoll     DiceRoll
}

// Valid Dice Rolls
var validDiceRolls = []DiceRoll{values4d4Plus1.diceRoll, values10d10.diceRoll, values1d6Minus1.diceRoll}

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
var invalidDiceRolls = []DiceRoll{valuesBigAmmount.diceRoll, valuesBigSize.diceRoll, valuesBigModifier.diceRoll, valuesZeroAmmount.diceRoll, valuesZeroSize.diceRoll}

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
	validateDiceRollString(values4d4Plus1.diceRoll, values4d4Plus1.asString, t)
	validateDiceRollString(values10d10.diceRoll, values10d10.asString, t)
	validateDiceRollString(values1d6Minus1.diceRoll, values1d6Minus1.asString, t)
}

func TestInvalidDiceRollString(t *testing.T) {
	validateDiceRollString(valuesBigAmmount.diceRoll, valuesBigAmmount.asString, t)
	validateDiceRollString(valuesBigSize.diceRoll, valuesBigSize.asString, t)
	validateDiceRollString(valuesBigModifier.diceRoll, valuesBigModifier.asString, t)
	validateDiceRollString(valuesZeroAmmount.diceRoll, valuesZeroAmmount.asString, t)
	validateDiceRollString(valuesZeroSize.diceRoll, valuesZeroSize.asString, t)
}

// Test valid DiceRoll results
func TestValidDiceRollResult(t *testing.T) {
	rollAndvalidateDiceRollResult(values4d4Plus1, t)
	rollAndvalidateDiceRollResult(values10d10, t)
	rollAndvalidateDiceRollResult(values1d6Minus1, t)
}

// Test invalid DiceRoll results
func TestInvalidDiceRollResult(t *testing.T) {
	validateInvalidDiceRoll(valuesBigAmmount, t)
	validateInvalidDiceRoll(valuesBigSize, t)
	validateInvalidDiceRoll(valuesBigModifier, t)
	validateInvalidDiceRoll(valuesZeroAmmount, t)
	validateInvalidDiceRoll(valuesZeroSize, t)
}

// Test rolling an array of valid DiceRoll
func TestPerformValidRolls(t *testing.T) {
	// Perform roll on DiceRoll array
	results, diceErrs := PerformRolls(validDiceRolls)
	if len(diceErrs) == 0 {
		if len(results) == len(validDiceRolls) {
			validateDiceRollResult(results[0], values4d4Plus1.resultFormat, t)
			validateDiceRollResult(results[1], values10d10.resultFormat, t)
			validateDiceRollResult(results[2], values1d6Minus1.resultFormat, t)
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
	// Perform roll on DiceRoll array
	_, diceErrs := PerformRolls(invalidDiceRolls)
	if len(diceErrs) > 0 {
		if len(diceErrs) != len(invalidDiceRolls) {
			// Missing results, fail the test
			t.Fatalf("Error list length = %d, want match for %d", len(diceErrs), len(invalidDiceRolls))
		}
	}
}

// Validate DiceRoll string format
func validateDiceRollString(diceRoll DiceRoll, asString string, t *testing.T) {
	if diceStr := diceRoll.String(); !strings.EqualFold(asString, diceStr) {
		t.Fatalf("DiceRoll = %s, must equal %s", diceStr, asString)
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
		t.Fatalf("Invalid dice roll did not generate an error: %s", diceValues.asString)
	}
}
