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
func NewDiceRoll(ammount int, size int, modifier int) (*DiceRoll, error) {
	diceRoll := DiceRoll{ammount, size, modifier}
	if diceErr := validateDiceRoll(diceRoll); diceErr == nil {
		return &diceRoll, nil
	} else {
		return nil, fmt.Errorf("invalid DiceRoll %s", diceErr.Error())
	}
}

// Performs the DiceRoll, returns the sum. Sum is 0 if DiceRoll is invalid.
func (diceRoll DiceRoll) PerformRollAndSum() int {
	result, _ := diceRoll.PerformRoll()
	return result.Sum
}

// Validates and performs the DiceRoll. Returns a DiceRollResult if valid, an error if invalid.
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

// Performs an array of DiceRoll, returns the sum. Invalid DiceRolls counts as 0.
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

// Validates DiceRoll values. Returns nil if valid, error if invalid.
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

// Validates ammount values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceAmmout(ammount int) error {
	if ammount > maxDiceRollValue || ammount <= 0 {
		return fmt.Errorf("invalid dice ammout %d. %s", ammount, bigNumberErrorMsg)
	} else if ammount <= 0 {
		return fmt.Errorf("invalid dice ammout %d", ammount)
	}
	return nil
}

// Validates size values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceSize(size int) error {
	if size > maxDiceRollValue || size <= 1 {
		return fmt.Errorf("invalid dice size %d. %s", size, bigNumberErrorMsg)
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
