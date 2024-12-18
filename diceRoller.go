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
	rollExprs, _ := parseRollArgs(rollArgs...)
	return performRollingExpressionsAndSum(rollExprs...)
}

// Performs an array of RollArgs. Returns a rollResult array for valid RollArgs and an error array for invalid ones.
func PerformRollArgs(rollArgs ...string) ([]rollResult, []error) {
	rollExprs, argErrs := parseRollArgs(rollArgs...)
	results, diceErrs := performRollingExpressions(rollExprs...)
	return results, append(argErrs, diceErrs...)
}

// Performs an array of DiceRoll. Returns the sum, invalid DiceRolls are worth 0.
func PerformDiceRollsAndSum(diceRolls ...DiceRoll) int {
	results, _ := PerformDiceRolls(diceRolls...)
	return RollResultsSum(results...)
}

// Performs an array of DiceRoll. Returns a rollResult array for valid DiceRolls and an error array for invalid ones.
func PerformDiceRolls(diceRolls ...DiceRoll) (results []rollResult, diceErrs []error) {
	return performRollingExpressions(*newRollingExpression(diceRolls...))
}

// Performs a rolling expression. Returns the sum, invalid DiceRolls are worth 0.
func performRollingExpressionsAndSum(rollExprs ...rollingExpression) int {
	results, _ := performRollingExpressions(rollExprs...)
	return RollResultsSum(results...)
}

// Performs a rolling expression. Returns a rollResult array for valid DiceRolls and an error array for invalid ones.
func performRollingExpressions(rollExprs ...rollingExpression) (results []rollResult, diceErrs []error) {
	wasCritHit := false
	for e := range rollExprs {
		rollExprResult := newRollResult()
		for i := range rollExprs[e].diceRolls {
			diceRoll := rollExprs[e].diceRolls[i]
			if wasCritHit {
				diceRoll.rollAttribs.setRollAttrib(critAttrib)
			}
			if result, diceErr := validateAndperformRoll(diceRoll); diceErr == nil {
				rollExprResult.results = append(rollExprResult.results, *result)
			} else {
				diceErrs = append(diceErrs, diceErr)
			}
		}

		wasCritHit = rollExprResult.detectScoredCritHit()

		results = append(results, *rollExprResult)
	}
	return results, diceErrs
}

// Validates and performs diceRoll. Returns a DiceRollResult if valid, an error if invalid.
func validateAndperformRoll(diceRoll DiceRoll) (*diceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Valid DiceRoll. Generate and return DiceRollResult
	return performRoll(diceRoll), nil
}

// Generates DiceRollResult and applies attribs.
func performRoll(diceRoll DiceRoll) *diceRollResult {
	diceRollResult := newDiceRollResult(diceRoll)

	// Generate rolls
	generateRolls(diceRoll, diceRollResult)

	// Drop High attrib
	if diceRoll.hasAttrib(dropHighAttrib) && len(diceRollResult.dice) > 1 {
		dropHigh(diceRollResult)
	}

	// Drop Low attrib
	if diceRoll.hasAttrib(dropLowAttrib) && len(diceRollResult.dice) > 1 {
		dropLow(diceRollResult)
	}

	// Apply modifier
	diceRollResult.sum += diceRoll.modifier

	// Half attrib
	if diceRoll.hasAttrib(halfAttrib) {
		diceRollResult.sum = halve(diceRollResult.sum)
	}

	// Minimum roll result is always 1, even after applying negative modifiers and half
	if diceRollResult.sum <= 0 {
		diceRollResult.sum = 1
	}

	// Negative Sum if minus DiceRoll
	if diceRoll.hasAttrib(minusAttrib) {
		diceRollResult.sum = -diceRollResult.sum
	}

	return diceRollResult
}

func generateRolls(diceRoll DiceRoll, diceRollResult *diceRollResult) {
	// Determine actual dice ammount to roll
	actualDiceAmmount := diceRoll.diceAmmount

	// Crit attrib
	if diceRoll.hasAttrib(critAttrib) {
		actualDiceAmmount = actualDiceAmmount * 2
	}

	// Generate rolls
	for i := 0; i < actualDiceAmmount; i++ {
		roll := rollDice(diceRoll.diceSize)

		// Advantage attrib
		if diceRoll.hasAttrib(advantageAttrib) {
			roll = advantage(roll, rollDice(diceRoll.diceSize), diceRollResult)
		}
		// Disadvantage attrib
		if diceRoll.hasAttrib(disadvantageAttrib) {
			roll = disadvantage(roll, rollDice(diceRoll.diceSize), diceRollResult)
		}

		diceRollResult.dice = append(diceRollResult.dice, roll)
		diceRollResult.sum += roll
	}
}

// Generates a single die roll.
func rollDice(diceSize int) int {
	return rand.Intn(diceSize) + 1
}

// Applies advantage logic. Returns the roll to keep and the roll to drop.
func advantage(roll int, roll2 int, diceRollResult *diceRollResult) (toKeep int) {
	toKeep, toDrop := roll, roll2 // Default return order, change if needed

	if roll != max(roll, roll2) {
		toKeep, toDrop = roll2, roll
	}

	diceRollResult.advDisDropped = append(diceRollResult.advDisDropped, toDrop)

	return toKeep
}

// Applies disavantage logic. Returns the roll to keep and the roll to drop.
func disadvantage(roll int, roll2 int, diceRollResult *diceRollResult) (toKeep int) {
	toKeep, toDrop := roll, roll2 // Default return order, change if needed

	if roll != min(roll, roll2) {
		toKeep, toDrop = roll2, roll
	}

	diceRollResult.advDisDropped = append(diceRollResult.advDisDropped, toDrop)

	return toKeep
}

// Applies drop high logic. Returns the index of the roll to drop.
func dropHigh(diceRollResult *diceRollResult) {
	drop := 0

	drop = slices.Max(diceRollResult.dice)
	dropIndex := slices.Index(diceRollResult.dice, drop)
	diceRollResult.highDropped = append(diceRollResult.highDropped, diceRollResult.dice[dropIndex])

	diceRollResult.sum -= diceRollResult.dice[dropIndex]
	diceRollResult.dice = slices.Delete(diceRollResult.dice, dropIndex, dropIndex+1)
}

// Applies drop low logic. Returns the index of the roll to drop.
func dropLow(diceRollResult *diceRollResult) {
	drop := 0

	drop = slices.Min(diceRollResult.dice)
	dropIndex := slices.Index(diceRollResult.dice, drop)
	diceRollResult.lowDropped = append(diceRollResult.lowDropped, diceRollResult.dice[dropIndex])

	diceRollResult.sum -= diceRollResult.dice[dropIndex]
	diceRollResult.dice = slices.Delete(diceRollResult.dice, dropIndex, dropIndex+1)
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
