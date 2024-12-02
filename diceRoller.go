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
		if result, diceErr := validateAndperformRoll(diceRolls[i]); diceErr == nil {
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
	wasCritHit := false
	for e := range rollExprs {
		rollExprResult := newRollingExpressionResult()
		for i := range rollExprs[e].diceRolls {
			diceRoll := rollExprs[e].diceRolls[i]
			if wasCritHit {
				diceRoll.Attribs.setRollAttrib(critAttrib)
			}
			if result, diceErr := validateAndperformRoll(diceRoll); diceErr == nil {
				rollExprResult.Results = append(rollExprResult.Results, *result)
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
func validateAndperformRoll(diceRoll DiceRoll) (*DiceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Valid DiceRoll. Generate and return DiceRollResult
	return performRoll(diceRoll), nil
}

// Generates DiceRollResult and applies attribs.
func performRoll(diceRoll DiceRoll) *DiceRollResult {
	diceRollResult := newDiceRollResult(diceRoll)
	dndRollAttributes, _ := diceRoll.Attribs.(*dndRollAttributes)

	// Generate rolls
	generateRolls(diceRoll, diceRollResult, dndRollAttributes)

	// Drop High attrib
	if dndRollAttributes.isDropHigh() && len(diceRollResult.Dice) > 1 {
		dropHigh(diceRollResult)
	}

	// Drop Low attrib
	if dndRollAttributes.isDropLow() && len(diceRollResult.Dice) > 1 {
		dropLow(diceRollResult)
	}

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier

	// Half attrib
	if dndRollAttributes.isHalf() {
		diceRollResult.Sum = halve(diceRollResult.Sum)
	}

	// Minimum roll result is always 1, even after applying negative modifiers and half
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	// Negative Sum if minus DiceRoll
	if !diceRoll.Plus {
		diceRollResult.Sum = -diceRollResult.Sum
	}

	return diceRollResult
}

func generateRolls(diceRoll DiceRoll, diceRollResult *DiceRollResult, dndRollAttributes *dndRollAttributes) {
	// Determine actual dice ammount to roll
	actualDiceAmmount := diceRoll.DiceAmmount
	if dndRollAttributes.isCrit() {
		actualDiceAmmount = actualDiceAmmount * 2
	}

	// Generate rolls
	for i := 0; i < actualDiceAmmount; i++ {
		roll := rollDice(diceRoll.DiceSize)

		// Advantage attrib
		if dndRollAttributes.isAdvantage() {
			roll = advantage(roll, rollDice(diceRoll.DiceSize), diceRollResult)
		}
		// Disadvantage attrib
		if dndRollAttributes.isDisadvantage() {
			roll = disadvantage(roll, rollDice(diceRoll.DiceSize), diceRollResult)
		}

		diceRollResult.Dice = append(diceRollResult.Dice, roll)
		diceRollResult.Sum += roll
	}
}

// Generates a single die roll.
func rollDice(diceSize int) int {
	return rand.Intn(diceSize) + 1
}

// Applies advantage logic. Returns the roll to keep and the roll to drop.
func advantage(roll int, roll2 int, diceRollResult *DiceRollResult) (toKeep int) {
	toKeep, toDrop := roll, roll2 // Default return order, change if needed

	if roll != max(roll, roll2) {
		toKeep, toDrop = roll2, roll
	}

	diceRollResult.AdvDisDropped = append(diceRollResult.AdvDisDropped, toDrop)

	return toKeep
}

// Applies disavantage logic. Returns the roll to keep and the roll to drop.
func disadvantage(roll int, roll2 int, diceRollResult *DiceRollResult) (toKeep int) {
	toKeep, toDrop := roll, roll2 // Default return order, change if needed

	if roll != min(roll, roll2) {
		toKeep, toDrop = roll2, roll
	}

	diceRollResult.AdvDisDropped = append(diceRollResult.AdvDisDropped, toDrop)

	return toKeep
}

// Applies drop high logic. Returns the index of the roll to drop.
func dropHigh(diceRollResult *DiceRollResult) {
	drop := 0

	drop = slices.Max(diceRollResult.Dice)
	dropIndex := slices.Index(diceRollResult.Dice, drop)
	diceRollResult.HighDropped = append(diceRollResult.HighDropped, diceRollResult.Dice[dropIndex])

	diceRollResult.Sum -= diceRollResult.Dice[dropIndex]
	diceRollResult.Dice = slices.Delete(diceRollResult.Dice, dropIndex, dropIndex+1)
}

// Applies drop low logic. Returns the index of the roll to drop.
func dropLow(diceRollResult *DiceRollResult) {
	drop := 0

	drop = slices.Min(diceRollResult.Dice)
	dropIndex := slices.Index(diceRollResult.Dice, drop)
	diceRollResult.LowDropped = append(diceRollResult.LowDropped, diceRollResult.Dice[dropIndex])

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
