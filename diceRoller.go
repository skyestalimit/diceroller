// Package diceroller generates either a sum or DiceRollResults out of
// RollArgs or DiceRolls. It's just that simple so get rolling!
package diceroller

import (
	"math"
	"math/rand"
	"slices"
)

// Straightforward rolling using RollArgs. Returns the sum, invalid RollArgs are worth 0.
func PerformRollArgsAndSum(rollArgs ...string) int {
	rollExprs, _ := ParseRollArgs(rollArgs...)
	return performRollingExpressionsAndSum(rollExprs...)
}

// Performs an array of RollArgs. Returns a DiceRollResult array for valid RollArgs and an error array for invalid ones.
func PerformRollArgs(rollArgs ...string) ([]RollingExpressionResult, []error) {
	rollExprs, argErrs := ParseRollArgs(rollArgs...)
	results, diceErrs := performRollingExpressions(rollExprs...)
	return results, append(argErrs, diceErrs...)
}

// Performs an array of DiceRoll. Returns the sum, invalid DiceRolls are worth 0.
func PerformRollsAndSum(diceRolls ...DiceRoll) int {
	results, _ := PerformRolls(diceRolls...)
	return DiceRollResultsSum(results...)
}

// Performs an array of DiceRoll. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func PerformRolls(diceRolls ...DiceRoll) (results []DiceRollResult, diceErrs []error) {
	for i := range diceRolls {
		if result, diceErr := performRoll(diceRolls[i]); diceErr == nil {
			results = append(results, *result)
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return results, diceErrs
}

// Performs a rolling expression. Returns the sum, invalid DiceRolls are worth 0.
func performRollingExpressionsAndSum(rollExprs ...rollingExpression) int {
	results, _ := performRollingExpressions(rollExprs...)
	return RollingExpressionResultSum(results...)
}

// Performs a rolling expression. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func performRollingExpressions(rollExprs ...rollingExpression) (results []RollingExpressionResult, diceErrs []error) {
	for e := range rollExprs {
		rollExprResult := newRollingExpressionResult()
		for i := range rollExprs[e].diceRolls {
			if result, diceErr := performRoll(rollExprs[e].diceRolls[i]); diceErr == nil {
				rollExprResult.Results = append(rollExprResult.Results, *result)
				rollExprResult.Sum += result.Sum
			} else {
				diceErrs = append(diceErrs, diceErr)
			}
		}
		results = append(results, *rollExprResult)
	}
	return results, diceErrs
}

// Validates and performs diceRoll. Returns a DiceRollResult if valid, an error if invalid.
func performRoll(diceRoll DiceRoll) (*DiceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Valid DiceRoll. Generate and return DiceRollResult
	return generateRolls(diceRoll), nil
}

// Generates DiceRollResult and applies attribs.
func generateRolls(diceRoll DiceRoll) *DiceRollResult {
	dndRollAttributes, _ := diceRoll.Attribs.(*dndRollAttributes)
	diceRollResult := newDiceRollResultWithAttribs(diceRoll.String(), dndRollAttributes)

	// Setup according to attribs
	diceAmmount := diceRoll.DiceAmmount
	var advDis rollAttribute = 0
	var dropHigh rollAttribute = 0
	var dropLow rollAttribute = 0
	var half rollAttribute = 0

	if dndRollAttributes != nil {
		for attrib := range dndRollAttributes.attribs {
			switch attrib {
			case critAttrib:
				// Crit attrib
				diceAmmount = diceAmmount * 2
			case advantageAttrib:
				advDis = advantageAttrib
			case disadvantageAttrib:
				advDis = disadvantageAttrib
			case dropHighAttrib:
				dropHigh = dropHighAttrib
			case dropLowAttrib:
				dropLow = dropLowAttrib
			case halfAttrib:
				half = halfAttrib
			}
		}
	}

	// Generate rolls
	for i := 0; i < diceAmmount; i++ {
		roll := rollDice(diceRoll.DiceSize)

		// Advantage / disadvantage attrib
		if advDis > 0 {
			roll = advantageDisadvantage(advDis, roll, rollDice(diceRoll.DiceSize), diceRollResult)
		}

		diceRollResult.Dice = append(diceRollResult.Dice, roll)
		diceRollResult.Sum += roll
	}

	// Drop High attrib
	if dropHigh > 0 && len(diceRollResult.Dice) > 1 {
		dropHighLow(dropHigh, diceRollResult)
	}

	// Drop Low attrib
	if dropLow > 0 && len(diceRollResult.Dice) > 1 {
		dropHighLow(dropLow, diceRollResult)
	}

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier

	// Half attrib
	if half == halfAttrib {
		diceRollResult.Sum = halve(diceRollResult.Sum)
	}

	// Minimum roll result is always 1, even after applying negative modifiers
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	// Negative Sum if minus DiceRoll
	if !diceRoll.Plus {
		diceRollResult.Sum = -diceRollResult.Sum
	}

	return diceRollResult
}

// Generates a single die roll.
func rollDice(diceSize int) int {
	return rand.Intn(diceSize) + 1
}

// Applies advantage or disavantage logic. Returns the roll to keep and the roll to drop.
func advantageDisadvantage(attrib rollAttribute, roll int, roll2 int, diceRollResult *DiceRollResult) (toKeep int) {
	toKeep, toDrop := roll, roll2 // Default return order, change if needed

	switch attrib {
	case advantageAttrib:
		if roll != max(roll, roll2) {
			toKeep, toDrop = roll2, roll
		}
	case disadvantageAttrib:
		if roll != min(roll, roll2) {
			toKeep, toDrop = roll2, roll
		}
	}

	diceRollResult.AdvDisDropped = append(diceRollResult.AdvDisDropped, toDrop)

	return toKeep
}

// Applies drop high or drop low logic. Returns the index of the roll to drop.
func dropHighLow(attrib rollAttribute, diceRollResult *DiceRollResult) {
	drop := 0

	switch attrib {
	case dropLowAttrib:
		drop = slices.Min(diceRollResult.Dice)
	case dropHighAttrib:
		drop = slices.Max(diceRollResult.Dice)
	}

	dropIndex := slices.Index(diceRollResult.Dice, drop)

	switch attrib {
	case dropLowAttrib:
		diceRollResult.LowDropped = append(diceRollResult.LowDropped, diceRollResult.Dice[dropIndex])
	case dropHighAttrib:
		diceRollResult.HighDropped = append(diceRollResult.HighDropped, diceRollResult.Dice[dropIndex])
	}

	diceRollResult.Sum -= diceRollResult.Dice[dropIndex]
	diceRollResult.Dice = slices.Delete(diceRollResult.Dice, dropIndex, dropIndex+1)
}

// Applies half logic. Rounds down.
func halve(sum int) int {
	halved := 0
	minus := false

	if sum < 0 {
		minus = true
	}

	halved = int(math.Abs(float64(sum))) / 2

	// Minimum result is 1 even if rounded down to 0
	if halved < 1 {
		halved = 1
	}

	if minus {
		halved = -halved
	}

	return halved
}
