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
var validDiceRollsValues = []diceRollTestValues{values4d4Plus1, values10d10, values1d6Minus1, values4d12Minus2500}

var values4d4Plus1 diceRollTestValues = diceRollTestValues{
	`4d4+1`,
	regexp.MustCompile(`\[[1-4] [1-4] [1-4] [1-4]\]`),
	*newDiceRoll(4, 4, 1)}
var values10d10 diceRollTestValues = diceRollTestValues{
	`10d10`,
	regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
	*newDiceRoll(10, 10, 0)}
var values1d6Minus1 diceRollTestValues = diceRollTestValues{
	`1d6-1`,
	regexp.MustCompile(`\[[1-6]\]`),
	*newDiceRoll(1, 6, -1)}
var values4d12Minus2500 diceRollTestValues = diceRollTestValues{
	`4d12-2500`,
	regexp.MustCompile(`\[([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9]) ([1]?[0-9])\]`),
	*newDiceRoll(4, 12, -2500)}

// Invalid Dice Rolls
var invalidDiceRollsValues = []diceRollTestValues{valuesBigAmmount, valuesBigSize, valuesBigModifier, valuesZeroAmmount, valuesZeroSize}

var valuesBigAmmount diceRollTestValues = diceRollTestValues{
	`123456d10`,
	nil,
	*newDiceRoll(123456, 10, 0)}
var valuesBigSize diceRollTestValues = diceRollTestValues{
	`10d123456+10`,
	nil,
	*newDiceRoll(10, 123456, 10)}
var valuesBigModifier diceRollTestValues = diceRollTestValues{
	`10d10-123456`,
	nil,
	*newDiceRoll(10, 10, -123456)}
var valuesZeroAmmount diceRollTestValues = diceRollTestValues{
	`0d8`,
	nil,
	*newDiceRoll(0, 8, 0)}
var valuesZeroSize diceRollTestValues = diceRollTestValues{
	`1d0`,
	nil,
	*newDiceRoll(1, 0, 0)}

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
func TestDiceRollPerformRollWithValidDiceRolls(t *testing.T) {
	for i := range validDiceRollsValues {
		if result, diceErr := performRoll(validDiceRollsValues[i].diceRoll); diceErr == nil {
			validateDiceRollResult(*result, validDiceRollsValues[i], t)
		} else {
			t.Fatalf("Unexpected dice roll error: %s", diceErr.Error())
		}
	}
}

// Test rolling invalid DiceRolls
func TestDiceRollPerformRollWithInvalidDiceRolls(t *testing.T) {
	for i := range invalidDiceRollsValues {
		if _, diceErr := performRoll(invalidDiceRollsValues[i].diceRoll); diceErr == nil {
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

// Test rolling with valid RollArgs
func TestPerformRollArgsWithValidRollArgs(t *testing.T) {
	// Test valid RollArgs individually
	for i := range validRollArgs {
		if sum := PerformRollArgsAndSum(validRollArgs[i]); sum <= 0 {
			t.Fatalf("Valid RollArg %s result = %d, want > 0", validRollArgs[i], sum)
		}
	}

	// Test the full array
	if sum := PerformRollArgsAndSum(validRollArgs...); sum <= 0 {
		t.Fatalf("Valid RollArgs result = %d, want > 0", sum)
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

// Test rolling with invalid RollArgs
func TestPerformRollArgsWithInvalidRollArgs(t *testing.T) {
	// Test valid RollArgs individually
	for i := range invalidRollArgs {
		if sum := PerformRollArgsAndSum(invalidRollArgs[i]); sum > 0 {
			t.Fatalf("Invalid RollArg %s result = %d, want 0", invalidRollArgs[i], sum)
		}
	}

	// Test the full array
	if sum := PerformRollArgsAndSum(invalidRollArgs...); sum > 0 {
		t.Fatalf("Invalid RollArgs result = %d, want 0", sum)
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

func FuzzPerformRollArgsAndSum(f *testing.F) {
	f.Add("8d4-1")
	f.Fuzz(func(t *testing.T, fuzzedRollArg string) {
		PerformRollArgsAndSum(fuzzedRollArg)
	})
}

func FuzzTestPerformRollArgs(f *testing.F) {
	f.Add("6d6+66")
	f.Fuzz(func(t *testing.T, fuzzedRollArg string) {
		PerformRollArgs(fuzzedRollArg)
	})
}

func FuzzTestPerformRollsAndSum(f *testing.F) {
	f.Add(6, 2, 4)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int) {
		PerformRollsAndSum(*newDiceRoll(diceAmmount, diceSize, modifier))
	})
}

func FuzzPerformRolls(f *testing.F) {
	f.Add(12, 18, 11)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int) {
		PerformRolls(*newDiceRoll(diceAmmount, diceSize, modifier))
	})
}

func FuzzNewDiceRoll(f *testing.F) {
	f.Add(2, 8, 1)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int) {
		NewDiceRoll(diceAmmount, diceSize, modifier)
	})
}

func FuzzDiceRollResultSum(f *testing.F) {
	f.Add("2d6+1")
	f.Fuzz(func(t *testing.T, diceStr string) {
		newDiceRollResult(diceStr)
		DiceRollResultsSum(*newDiceRollResult(diceStr))
	})
}
