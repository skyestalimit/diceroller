// Package diceroller can generate either a sum or DiceRollResults out of
// rollArgs or DiceRolls. It's just that simple so get rolling!
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
	if diceErr := validateDiceRoll(diceRoll); diceErr == nil {
		return &diceRoll, nil
	} else {
		return nil, fmt.Errorf("invalid DiceRoll %s", diceErr.Error())
	}
}

// Performs the DiceRoll. Returns the sum if diceRoll valid, zero if invalid.
func (diceRoll DiceRoll) PerformRollAndSum() int {
	result, _ := diceRoll.PerformRoll()
	return result.Sum
}

// Validates and performs diceRoll. Returns a DiceRollResult if valid, an error if invalid.
func (diceRoll DiceRoll) PerformRoll() (*DiceRollResult, error) {
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

// Human readable DiceRoll string in XdY([+|-]Z) format.
// DiceRoll{2, 8, 1}.PerformRoll() would return "2d8+1"
func (diceRoll DiceRoll) String() string {
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)

	// Add modifier when necessary
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += plusSymbol
		}
		strDiceRoll += fmt.Sprint(diceRoll.Modifier)
	}
	return strDiceRoll
}

// Straightforward rolling. Rolls diceAmmount times a diceSize sized dice plus modifier.
// Returns the sum if valid, zero if invalid.
func PerformRollAndSum(diceAmmount int, diceSize int, modifier int) int {
	if result, _ := PerformRoll(diceAmmount, diceSize, modifier); result != nil {
		return result.Sum
	} else {
		return 0
	}
}

// Straightforward rolling. Rolls diceAmmount times a diceSize sized dice plus modifier.
// Returns a DiceRollResult if valid, an error if invalid.
func PerformRoll(diceAmmount int, diceSize int, modifier int) (*DiceRollResult, error) {
	if diceRoll, diceErr := NewDiceRoll(diceAmmount, diceSize, modifier); diceErr != nil {
		return nil, diceErr
	} else {
		return diceRoll.PerformRoll()
	}
}

// Straightforward rolling using rollArgs. Returns the sum, invalid rollArgs are worth 0.
func PerformRollArgsAndSum(rollArgs ...string) int {
	diceRolls, _ := ParseRollArgs(rollArgs...)
	return PerformRollsAndSum(diceRolls...)
}

// Performs an array of rollArgs. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func PerformRollArgs(rollArgs ...string) ([]DiceRollResult, []error) {
	diceRolls, argErrs := ParseRollArgs(rollArgs...)
	results, diceErrs := PerformRolls(diceRolls...)
	argErrs = append(argErrs, diceErrs...)
	return results, argErrs
}

// Performs an array of DiceRoll. Returns the sum, invalid DiceRolls are worth 0.
func PerformRollsAndSum(diceRolls ...DiceRoll) int {
	results, _ := PerformRolls(diceRolls...)
	return DiceRollResultsSum(results)
}

// Performs an array of DiceRoll. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func PerformRolls(diceRolls ...DiceRoll) (results []DiceRollResult, diceErrs []error) {
	for i := range diceRolls {
		if result, diceErr := diceRolls[i].PerformRoll(); diceErr == nil {
			results = append(results, *result)
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return results, diceErrs
}

// Validates diceRoll values. Returns nil if valid, error if invalid.
func validateDiceRoll(diceRoll DiceRoll) error {
	if diceErr := validateDiceAmmout(diceRoll.DiceAmmount); diceErr != nil {
		return formattedInvalidDiceRollError(diceRoll.String(), diceErr)
	}
	if diceErr := validateDiceSize(diceRoll.DiceSize); diceErr != nil {
		return formattedInvalidDiceRollError(diceRoll.String(), diceErr)
	}
	if diceErr := validateDiceModifier(diceRoll.Modifier); diceErr != nil {
		return formattedInvalidDiceRollError(diceRoll.String(), diceErr)
	}
	return nil
}

// Validates diceAmmount values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceAmmout(diceAmmount int) error {
	if diceAmmount > maxDiceRollValue || diceAmmount <= 0 {
		return fmt.Errorf("invalid dice ammout %d. %s", diceAmmount, bigNumberErrorMsg)
	} else if diceAmmount <= 0 {
		return fmt.Errorf("invalid dice ammout %d", diceAmmount)
	}
	return nil
}

// Validates diceSize values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceSize(diceSize int) error {
	if diceSize > maxDiceRollValue || diceSize <= 1 {
		return fmt.Errorf("invalid dice size %d. %s", diceSize, bigNumberErrorMsg)
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

// Returns a formatted invalid DiceRoll error.
func formattedInvalidDiceRollError(diceRollStr string, diceErr error) error {
	return fmt.Errorf("%s: %s", diceRollStr, diceErr.Error())
}
