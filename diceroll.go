package diceroller

import (
	"fmt"
	"math"
)

// A DiceRoll represents a dice rolling expression, such as 1d6 or 2d8+1.
type DiceRoll struct {
	diceAmmount int             // Ammount of dice to be rolled
	diceSize    int             // Size, or number of faces, of the dice to be rolled
	modifier    int             // Value to be applied to the sum of rolled dices
	rollAttribs *rollAttributes // Contains rollAttributes affecting the rolls
}

// Max allowed DiceRoll values to avoid long run times and overflow.
const maxDiceRollValue int = 99999

// Ridiculous error message for ridiculously big values.
const bigNumberErrorMsg = "This is a dice roller, not a Pi calculator"

// DiceRoll constructor, validates values.
func NewDiceRoll(diceAmmount int, diceSize int, modifier int) (*DiceRoll, error) {
	return NewDiceRollWithAttribs(diceAmmount, diceSize, modifier, newRollAttributes())
}

// DiceRoll constructor with rollAttributes, validates values.
func NewDiceRollWithAttribs(diceAmmount int, diceSize int, modifier int, attribs *rollAttributes) (*DiceRoll, error) {
	if attribs == nil {
		attribs = newRollAttributes()
	}
	diceRoll := DiceRoll{diceAmmount, diceSize, modifier, attribs}
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		return nil, fmt.Errorf("invalid DiceRoll %s", diceErr.Error())
	}
	return &diceRoll, nil
}

// DiceRoll constructor, validates values but doesn't return errors. Can be useful for testing.
func newDiceRoll(diceAmmount int, diceSize int, modifier int) *DiceRoll {
	diceRoll, _ := NewDiceRoll(diceAmmount, diceSize, modifier)
	return diceRoll
}

// Performs the DiceRoll. Returns the sum if valid, zero if invalid.
func (diceRoll DiceRoll) Roll() int {
	result, _ := validateAndperformRoll(diceRoll)
	return result.sum
}

// Returns true if wanted is set. Provides nil protection that rollAttribute can't provide itself.
func (diceRoll DiceRoll) hasAttrib(wanted rollAttribute) bool {
	found := false
	if diceRoll.rollAttribs != nil {
		found = diceRoll.rollAttribs.hasAttrib(wanted)
	}
	return found
}

// Human readable DiceRoll string, such as "2d8+1".
func (diceRoll DiceRoll) String() string {
	strDiceRoll := ""

	// Add minus symbol if needed
	if diceRoll.hasAttrib(minusAttrib) {
		strDiceRoll += "-"
	}

	// XdY format
	strDiceRoll += fmt.Sprintf("%dd%d", diceRoll.diceAmmount, diceRoll.diceSize)

	// Add modifier when necessary
	if diceRoll.modifier != 0 {
		if diceRoll.modifier > 0 {
			strDiceRoll += "+"
		}
		strDiceRoll += fmt.Sprint(diceRoll.modifier)
	}

	return strDiceRoll
}

// Validates diceRoll values. Returns nil if valid, error if invalid.
func validateDiceRoll(diceRoll DiceRoll) error {
	if diceErr := validateDiceAmmout(diceRoll.diceAmmount); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceSize(diceRoll.diceSize); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceModifier(diceRoll.modifier); diceErr != nil {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error())
	}
	return nil
}

// Validates diceAmmount values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceAmmout(diceAmmount int) error {
	ammountErrStr := "invalid dice ammout"
	if diceAmmount > maxDiceRollValue {
		return fmt.Errorf("%s %d. %s", ammountErrStr, diceAmmount, bigNumberErrorMsg)
	} else if diceAmmount <= 0 {
		return fmt.Errorf("%s %d", ammountErrStr, diceAmmount)
	}
	return nil
}

// Validates diceSize values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceSize(diceSize int) error {
	sizeErrStr := "invalid dice size"
	if diceSize > maxDiceRollValue {
		return fmt.Errorf("%s %d. %s", sizeErrStr, diceSize, bigNumberErrorMsg)
	} else if diceSize <= 1 {
		return fmt.Errorf("%s %d", sizeErrStr, diceSize)
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
