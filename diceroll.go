// Package diceroller generates either a sum or DiceRollResults out of
// RollArgs or DiceRolls. It's just that simple so get rolling!
package diceroller

import (
	"fmt"
	"math"
	"math/rand"
)

// A DiceRoll represents a dice rolling expression, such as 1d6 or 2d8+1.
type DiceRoll struct {
	DiceAmmount int // Ammount of dice to be rolled
	DiceSize    int // Size, or number of faces, of the dice to be rolled
	Modifier    int // Value to be applied to the sum of rolled dices
}

// Max allowed values for DiceRoll to avoid long run times and overflow.
const maxDiceRollValue int = 99999

// Ridiculous error message to send back for ridiculously big values.
const bigNumberErrorMsg = "This is a dice roller, not a Pi calculator"

// DiceRoll constructor, validates values.
func NewDiceRoll(diceAmmount int, diceSize int, modifier int) (*DiceRoll, error) {
	diceRoll := DiceRoll{diceAmmount, diceSize, modifier}
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		return nil, fmt.Errorf("invalid DiceRoll %s", diceErr.Error())
	}
	return &diceRoll, nil
}

// Human readable DiceRoll string, such as "2d8+1".
func (diceRoll DiceRoll) String() string {
	// XdY format
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)

	// Add modifier when necessary
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += "+"
		}
		strDiceRoll += fmt.Sprint(diceRoll.Modifier)
	}

	return strDiceRoll
}

// Straightforward rolling using RollArgs. Returns the sum, invalid RollArgs are worth 0.
func PerformRollArgsAndSum(rollArgs ...string) int {
	diceRolls, _ := ParseRollArgs(rollArgs...)
	return PerformRollsAndSum(diceRolls...)
}

// Performs an array of RollArgs. Returns a DiceRollResult array for valid RollArgs and an error array for invalid ones.
func PerformRollArgs(rollArgs ...string) ([]DiceRollResult, []error) {
	diceRolls, argErrs := ParseRollArgs(rollArgs...)
	results, diceErrs := PerformRolls(diceRolls...)
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

// Validates and performs diceRoll. Returns a DiceRollResult if valid, an error if invalid.
func performRoll(diceRoll DiceRoll) (*DiceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Generate rolls
	diceRollResult := NewDiceRollResult(diceRoll.String())
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceRollResult.Dice = append(diceRollResult.Dice, rand.Intn(diceRoll.DiceSize)+1)
		diceRollResult.Sum += diceRollResult.Dice[i]
	}

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier

	// Minimum roll result is always 1, even after applying negative modifiers
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	// Valid DiceRoll, return DiceRollResult
	return diceRollResult, nil
}

// Validates diceRoll values. Returns nil if valid, error if invalid.
func validateDiceRoll(diceRoll DiceRoll) error {
	if diceErr := validateDiceAmmout(diceRoll.DiceAmmount); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceSize(diceRoll.DiceSize); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceModifier(diceRoll.Modifier); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	return nil
}

// Validates diceAmmount values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceAmmout(diceAmmount int) error {
	if diceAmmount > maxDiceRollValue {
		return fmt.Errorf("invalid dice ammout %d. %s", diceAmmount, bigNumberErrorMsg)
	} else if diceAmmount <= 0 {
		return fmt.Errorf("invalid dice ammout %d", diceAmmount)
	}
	return nil
}

// Validates diceSize values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceSize(diceSize int) error {
	if diceSize > maxDiceRollValue {
		return fmt.Errorf("invalid dice size %d. %s", diceSize, bigNumberErrorMsg)
	} else if diceSize <= 1 {
		return fmt.Errorf("invalid dice size %d", diceSize)
	}
	return nil
}

// Validates modifier values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceModifier(modifier int) error {
	if int(math.Abs(float64(modifier))) > maxDiceRollValue {
		return fmt.Errorf("invalid dice modifier %d. %s", modifier, bigNumberErrorMsg)
	}
	return nil
}
