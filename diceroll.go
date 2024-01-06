// Package diceroller generates dice roll results based on the values
// defined within a DiceRole object. It returns the results in a
// DiceRollResult struct.
package diceroller

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// A DiceRoll represents a dice rolling expression, such as 1d6 or 2d8+1.
type DiceRoll struct {
	DiceAmmount int // Ammount of dice to be rolled
	DiceSize    int // Size, or numer of faces, of the dice to be rolled
	Modifier    int // Value to be applied to the sum of rolled dices
}

// Max allowed value for DiceRoll definition.
const maxDiceRollValue int = 99999

// DiceRoll constructor. Validates values.
func NewDiceRoll(ammount int, size int, modifier int) (*DiceRoll, error) {
	diceRoll := DiceRoll{ammount, size, modifier}
	if diceErr := validateDiceRoll(diceRoll); diceErr == nil {
		return &diceRoll, nil
	} else {
		return nil, fmt.Errorf("invalid DiceRoll: %s", diceErr.Error())
	}
}

// Perform an array of DiceRoll. Returns a DiceRollResult array for valid DiceRolls and an error array for invalid ones.
func PerformRolls(diceRolls []DiceRoll) (results []DiceRollResult, diceErrs []error) {
	for i := range diceRolls {
		if result, diceErr := diceRolls[i].PerformRoll(); diceErr == nil {
			results = append(results, *result)
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return results, diceErrs
}

// Validates and performs a DiceRoll. Returns DiceRollResult if valid, an error if invalid.
func (diceRoll DiceRoll) PerformRoll() (*DiceRollResult, error) {
	// Validate DiceRoll
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		// Invalid DiceRoll, return error
		return nil, diceErr
	}

	// Generate rolls
	diceRollResult := NewDiceRollResult(diceRoll.String())
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceGen := getFreshRandomGenerator()
		diceRollResult.Dice = append(diceRollResult.Dice, diceGen.Intn(diceRoll.DiceSize)+1)
		diceRollResult.Sum += diceRollResult.Dice[i]
	}

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier

	// Minimum roll result if always 1
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	// Valid DiceRoll, return DiceRollResult
	return diceRollResult, nil
}

// Human readable DiceRoll string
func (diceRoll DiceRoll) String() string {
	// XdX format
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)

	// Add modifier when necessary
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += plusSymbol
		}
		strDiceRoll += fmt.Sprintf("%d", diceRoll.Modifier)
	}
	return strDiceRoll
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

// Validates ammount value for DiceRoll. Returns nil if valid, error if invalid.
func validateDiceAmmout(ammount int) error {
	if ammount > maxDiceRollValue || ammount <= 0 {
		return fmt.Errorf("invalid dice ammout %d", ammount)
	}
	return nil
}

// Validates size value for DiceRoll. Returns nil if valid, error if invalid.
func validateDiceSize(size int) error {
	if size > maxDiceRollValue || size <= 1 {
		return fmt.Errorf("invalid dice size %d", size)
	}
	return nil
}

// Validates modifier value for DiceRoll. Returns nil if valid, error if invalid.
func validateDiceModifier(modifier int) error {
	if int(math.Abs(float64(modifier))) > maxDiceRollValue {
		return fmt.Errorf("invalid dice modifier %d", modifier)
	}
	return nil
}

// Returns a formatted invalid dice roll error.
func formattedInvalidDiceRollError(diceRollStr string, diceErr error) error {
	return fmt.Errorf("dice roll %s: %s", diceRollStr, diceErr.Error())
}

// Seeds a fresh random generator.
func getFreshRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
