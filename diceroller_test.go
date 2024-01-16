package diceroller

import (
	"testing"
)

// Test valid DiceRoll results
func TestPerformRollWithValidDiceRolls(t *testing.T) {
	for i := range validDiceRollsValues {
		result, diceErr := performRoll(nil, validDiceRollsValues[i].diceRoll)
		if diceErr != nil {
			t.Fatalf("Unexpected dice roll error: %s", diceErr.Error())
		}
		validateDiceRollResult(*result, validDiceRollsValues[i], t)
	}
}

// Test rolling invalid DiceRolls
func TestPerformRollWithInvalidDiceRolls(t *testing.T) {
	for i := range invalidDiceRollsValues {
		_, diceErr := performRoll(nil, invalidDiceRollsValues[i].diceRoll)
		if diceErr == nil {
			t.Fatalf("Invalid dice roll %s did not generate an error", invalidDiceRollsValues[i].wantedDiceStr)
		}
	}
}

// Test rolling an array of valid DiceRoll
func TestPerformRollsWithValidDiceRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	results, diceErrs := PerformRolls(nil, diceRollsFromTestValues(validDiceRollsValues)...)

	// No errors should be received for a valid DiceRoll array
	if len(diceErrs) != 0 {
		// Received errors, fail the test
		errStr := ""
		for i := range diceErrs {
			errStr += diceErrs[i].Error()
		}
		t.Fatalf("Unexpected dice roll errors: %s", errStr)

	}

	// One result per valid DiceRoll should be received
	if len(results) != len(validDiceRollsValues) {
		// Missing results, fail the test
		t.Fatalf("Result list length = %d, want match for %d", len(results), len(validDiceRollsValues))
	}

	// Validate result array
	for i := range results {
		validateDiceRollResult(results[i], validDiceRollsValues[i], t)
	}
}

// Test rolling an array of invalid DiceRoll
func TestPerformRollsWithInvalidDiceRolls(t *testing.T) {
	// Perform rolls on DiceRoll array
	_, diceErrs := PerformRolls(nil, diceRollsFromTestValues(invalidDiceRollsValues)...)

	// One error per invalid DiceRoll should be received
	if len(diceErrs) != len(invalidDiceRollsValues) {
		t.Fatalf("Error list length = %d, want match for %d", len(diceErrs), len(invalidDiceRollsValues))
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

	// Test sending the whole array
	if sum := PerformRollArgsAndSum(validRollArgs...); sum <= 0 {
		t.Fatalf("Valid RollArgs result = %d, want > 0", sum)
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

// Validates roll result matches expected format
func validateDiceRollResult(result DiceRollResult, diceValues diceRollTestValues, t *testing.T) {
	// validate result format
	if resultStr := result.String(); !diceValues.resultFormat.MatchString(result.String()) {
		t.Fatalf("Roll result = %s, want match for %#q", resultStr, diceValues.resultFormat)
	}

	// Reproduce sum calculations
	sum := 0
	for i := range result.Dice {
		sum += result.Dice[i]
	}
	sum += diceValues.diceRoll.Modifier
	if sum <= 0 {
		sum = 1
	}
	if !diceValues.diceRoll.Plus {
		sum = -sum
	}

	// Validate sum
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
	f.Add(6, 2, 4, true)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int, plus bool) {
		PerformRollsAndSum(nil, *newDiceRoll(diceAmmount, diceSize, modifier, plus))
	})
}

func FuzzPerformRolls(f *testing.F) {
	f.Add(12, 18, 11, true)
	f.Fuzz(func(t *testing.T, diceAmmount int, diceSize int, modifier int, plus bool) {
		PerformRolls(nil, *newDiceRoll(diceAmmount, diceSize, modifier, plus))
	})
}
