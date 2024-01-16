package diceroller

import (
	"fmt"
	"math"
)

// A DiceRoll represents a dice rolling expression, such as 1d6 or 2d8+1.
type DiceRoll struct {
	DiceAmmount int  // Ammount of dice to be rolled
	DiceSize    int  // Size, or number of faces, of the dice to be rolled
	Modifier    int  // Value to be applied to the sum of rolled dices
	Plus        bool // Determines if the result of the roll is to be added or substracted
	attribs     attributes
}

// Max allowed values for DiceRoll to avoid long run times and overflow.
const maxDiceRollValue int = 99999

// Ridiculous error message to send back for ridiculously big values.
const bigNumberErrorMsg = "This is a dice roller, not a Pi calculator"

// DiceRoll constructor, validates values.
func NewDiceRoll(diceAmmount int, diceSize int, modifier int, plus bool) (*DiceRoll, error) {
	return NewDiceRollWithAttribs(diceAmmount, diceSize, modifier, plus, newRollAttributes())
}

// DiceRoll constructor with rollAttributes, validates values.
func NewDiceRollWithAttribs(diceAmmount int, diceSize int, modifier int, plus bool, attribs *rollAttributes) (*DiceRoll, error) {
	diceRoll := DiceRoll{diceAmmount, diceSize, modifier, plus, attribs}
	if diceErr, ok := validateDiceRoll(diceRoll); !ok {
		return nil, fmt.Errorf("invalid DiceRoll %s", diceErr.Error())
	}
	return &diceRoll, nil
}

// DiceRoll constructor without errors.
func newDiceRoll(diceAmmount int, diceSize int, modifier int, plus bool) *DiceRoll {
	diceRoll, _ := NewDiceRoll(diceAmmount, diceSize, modifier, plus)
	return diceRoll
}

// Human readable DiceRoll string, such as "2d8+1".
func (diceRoll DiceRoll) String() string {
	strDiceRoll := ""

	// Add minus symbol if needed
	if !diceRoll.Plus {
		strDiceRoll += "-"
	}

	// XdY format
	strDiceRoll += fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)

	// Add modifier when necessary
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += "+"
		}
		strDiceRoll += fmt.Sprint(diceRoll.Modifier)
	}

	return strDiceRoll
}

// Validates diceRoll values. Returns nil if valid, error if invalid.
func validateDiceRoll(diceRoll DiceRoll) (error, bool) {
	if diceErr, ok := validateDiceAmmout(diceRoll.DiceAmmount); !ok {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error()), false
	}
	if diceErr, ok := validateDiceSize(diceRoll.DiceSize); !ok {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error()), false
	}
	if diceErr, ok := validateDiceModifier(diceRoll.Modifier); !ok {
		return fmt.Errorf("%s: %s", diceRoll.String(), diceErr.Error()), false
	}
	return nil, true
}

// Validates diceAmmount values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceAmmout(diceAmmount int) (error, bool) {
	if diceAmmount > maxDiceRollValue {
		return fmt.Errorf("invalid dice ammout %d. %s", diceAmmount, bigNumberErrorMsg), false
	} else if diceAmmount <= 0 {
		return fmt.Errorf("invalid dice ammout %d", diceAmmount), false
	}
	return nil, true
}

// Validates diceSize values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceSize(diceSize int) (error, bool) {
	if diceSize > maxDiceRollValue {
		return fmt.Errorf("invalid dice size %d. %s", diceSize, bigNumberErrorMsg), false
	} else if diceSize <= 1 {
		return fmt.Errorf("invalid dice size %d", diceSize), false
	}
	return nil, true
}

// Validates modifier values for DiceRoll. Returns nil if valid, an error if invalid.
func validateDiceModifier(modifier int) (error, bool) {
	if int(math.Abs(float64(modifier))) > maxDiceRollValue {
		return fmt.Errorf("invalid dice modifier %d. %s", modifier, bigNumberErrorMsg), false
	}
	return nil, true
}
