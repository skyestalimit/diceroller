// Package diceroller generates either a sum or DiceRollResults out of
// RollArgs or DiceRolls. It's just that simple so get rolling!
package diceroller

import "math/rand"

// Straightforward rolling using RollArgs. Returns the sum, invalid RollArgs are worth 0.
func PerformRollArgsAndSum(rollArgs ...string) int {
	rollExpr, _ := ParseRollArgs(rollArgs...)
	return performRollsAndSum(rollExpr)
}

// Performs an array of RollArgs. Returns a DiceRollResult array for valid RollArgs and an error array for invalid ones.
func PerformRollArgs(rollArgs ...string) ([]DiceRollResult, []error) {
	rollExpr, argErrs := ParseRollArgs(rollArgs...)
	results, diceErrs := performRolls(rollExpr)
	return results, append(argErrs, diceErrs...)
}

// Performs an array of DiceRoll. Returns the sum, invalid DiceRolls are worth 0.
func PerformRollsAndSum(attribs *rollAttributes, diceRolls ...DiceRoll) int {
	results, _ := PerformRolls(attribs, diceRolls...)
	return DiceRollResultsSum(results...)
}

// Performs an array of DiceRoll. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func PerformRolls(attribs *rollAttributes, diceRolls ...DiceRoll) (results []DiceRollResult, diceErrs []error) {
	for i := range diceRolls {
		if result, diceErr := performRoll(attribs, diceRolls[i]); diceErr == nil {
			results = append(results, *result)
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return results, diceErrs
}

func performRollsAndSum(rollExpr rollingExpression) int {
	results, _ := performRolls(rollExpr)
	return DiceRollResultsSum(results...)
}

func performRolls(rollExpr rollingExpression) (results []DiceRollResult, diceErrs []error) {
	rollAttribs := rollExpr.attribs.(*rollAttributes)
	for i := range rollExpr.diceRolls {
		if result, diceErr := performRoll(rollAttribs, rollExpr.diceRolls[i].(DiceRoll)); diceErr == nil {
			results = append(results, *result)
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return results, diceErrs
}

// Validates and performs diceRoll. Returns a DiceRollResult if valid, an error if invalid.
func performRoll(attribs *rollAttributes, diceRoll DiceRoll) (*DiceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Generate rolls
	diceRollResult := generateRolls(attribs, diceRoll)

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier

	// Minimum roll result is always 1, even after applying negative modifiers
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	// Negative Sum if minus DiceRoll
	if !diceRoll.Plus {
		diceRollResult.Sum = -diceRollResult.Sum
	}

	// Valid DiceRoll, return DiceRollResult
	return diceRollResult, nil
}

func generateRolls(attribs *rollAttributes, diceRoll DiceRoll) *DiceRollResult {
	diceRollResult := newDiceRollResult(diceRoll.String())
	diceRollResult.attribs = attribs

	// Setup according to attribs
	diceAmmount := diceRoll.DiceAmmount
	if attribs != nil && attribs.hasAttrib(critAttrib) {
		diceAmmount = diceAmmount * 2
	}

	// Generate rolls
	for i := 0; i < diceAmmount; i++ {
		roll := rollDice(diceRoll.DiceSize)
		if attribs != nil && (attribs.hasAttrib(advantageAttrib) || attribs.hasAttrib(disadvantageAttrib)) {
			roll2 := rollDice(diceRoll.DiceSize)
			toDrop := 0
			roll, toDrop = advantageDisadvantage(attribs, roll, roll2)
			diceRollResult.Dropped = append(diceRollResult.Dropped, toDrop)
		}

		diceRollResult.Dice = append(diceRollResult.Dice, roll)
		diceRollResult.Sum += roll
	}

	return diceRollResult
}

func rollDice(diceSize int) int {
	return rand.Intn(diceSize) + 1
}

func advantageDisadvantage(attribs *rollAttributes, roll int, roll2 int) (int, int) {
	toDrop := 0
	if attribs.hasAttrib(advantageAttrib) {
		if roll == max(roll, roll2) {
			toDrop = roll2
		} else {
			toDrop, roll = roll, roll2
		}
	} else {
		if roll == min(roll, roll2) {
			toDrop = roll2
		} else {
			toDrop, roll = roll, roll2
		}
	}
	return roll, toDrop
}
