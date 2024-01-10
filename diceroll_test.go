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
	regexp.MustCompile(`\[[1-4] [1-4] [1-4] [1-4]\]`),
	DiceRoll{4, 4, 1}}
var values10d10 diceRollTestValues = diceRollTestValues{
	`10d10`,
	regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
	DiceRoll{10, 10, 0}}
var values1d6Minus1 diceRollTestValues = diceRollTestValues{
	`1d6-1`,
	regexp.MustCompile(`\[[1-6]\]`),
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

// Test simple rolling with valid values
func TestPerformRollWithValidDiceValues(t *testing.T) {
	for i := range validDiceRollsValues {
		diceRoll := validDiceRollsValues[i].diceRoll
		if sum := PerformRollAndSum(diceRoll.DiceAmmount, diceRoll.DiceSize, diceRoll.Modifier); sum <= 0 {
			t.Fatalf("Valid diceRoll %s returned invalid result %d", diceRoll.String(), sum)
		}
		if result, diceErr := PerformRoll(diceRoll.DiceAmmount, diceRoll.DiceSize, diceRoll.Modifier); diceErr != nil {
			t.Fatalf("Unexpected dice roll error: %s", diceErr.Error())
		} else {
			validateDiceRollResult(*result, validDiceRollsValues[i], t)
		}
	}
}

// Test simple rolling with invalid values
func TestPerformRollWithInvalidDiceValues(t *testing.T) {
	for i := range invalidDiceRollsValues {
		diceRoll := invalidDiceRollsValues[i].diceRoll
		if sum := PerformRollAndSum(diceRoll.DiceAmmount, diceRoll.DiceSize, diceRoll.Modifier); sum > 0 {
			t.Fatalf("Invalid diceRoll %s returned valid result %d", diceRoll.String(), sum)
		}
		if _, diceErr := PerformRoll(diceRoll.DiceAmmount, diceRoll.DiceSize, diceRoll.Modifier); diceErr == nil {
			t.Fatalf("Invalid dice roll %s did not generate an error", diceRoll.String())

		}
	}
}

// Test valid DiceRoll results
func TestDiceRollPerformRollWithValidDiceRolls(t *testing.T) {
	for i := range validDiceRollsValues {
		if result, diceErr := validDiceRollsValues[i].diceRoll.PerformRoll(); diceErr == nil {
			validateDiceRollResult(*result, validDiceRollsValues[i], t)
		} else {
			t.Fatalf("Unexpected dice roll error: %s", diceErr.Error())
		}
	}
}

// Test rolling invalid DiceRolls
func TestDiceRollPerformRollWithInvalidDiceRolls(t *testing.T) {
	for i := range invalidDiceRollsValues {
		if _, diceErr := invalidDiceRollsValues[i].diceRoll.PerformRoll(); diceErr == nil {
			t.Fatalf("Invalid dice roll %s did not generate an error", invalidDiceRollsValues[i].wantedDiceStr)
		}
	}
}

// Test rolling an array of valid DiceRoll
func TestPerformRollsWithValidDiceRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	if results, diceErrs := PerformRolls(diceRollsFromTestValues(validDiceRollsValues)...); len(diceErrs) == 0 {
		if len(results) == len(validDiceRollsValues) {
			for i := range results {
				validateDiceRollResult(results[i], validDiceRollsValues[i], t)
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
func TestPerformRollsWithInvalidDiceRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	if _, diceErrs := PerformRolls(diceRollsFromTestValues(invalidDiceRollsValues)...); len(diceErrs) > 0 {
		if len(diceErrs) != len(invalidDiceRollsValues) {
			// Missing errors, fail the test
			t.Fatalf("Error list length = %d, want match for %d", len(diceErrs), len(invalidDiceRollsValues))
		}
	}
}

// Validates DiceRoll string format
func validateDiceRollString(diceValues diceRollTestValues, t *testing.T) {
	if diceStr := diceValues.diceRoll.String(); !strings.EqualFold(diceValues.wantedDiceStr, diceStr) {
		t.Fatalf("DiceRoll = %s, must equal %s", diceStr, diceValues.wantedDiceStr)
	}
}

// Validates roll result matches expected format
func validateDiceRollResult(result DiceRollResult, diceValues diceRollTestValues, t *testing.T) {
	// validate result format
	if resultStr := result.String(); !diceValues.resultFormat.MatchString(result.String()) {
		t.Fatalf("Roll result = %s, want match for %#q", resultStr, diceValues.resultFormat)
	}

	// Validate sum
	sum := 0
	for i := range result.Dice {
		sum += result.Dice[i]
	}
	sum += diceValues.diceRoll.Modifier
	if sum <= 0 {
		sum = 1
	}
	if result.Sum != sum {
		t.Fatalf("DiceRoll %s result is %d, wanted be %d", diceValues.diceRoll.String(), result.Sum, sum)
	}
}

// Extracts and returns a DiceRoll array from a diceRollTestValues array
func diceRollsFromTestValues(testValues []diceRollTestValues) (diceRolls []DiceRoll) {
	for i := range testValues {
		diceRolls = append(diceRolls, testValues[i].diceRoll)
	}
	return
}
