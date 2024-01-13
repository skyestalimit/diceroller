package diceroller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parses a RollArg array. Returns a DiceRoll array for valid RollArgs
// and an error array for invalid ones.
//
// Expected format is XdY[+|-Z].
//
// Valid RollArg examples: "4d4+1", "10d10", "1d6-1", "1D8".
func ParseRollArgs(rollArgs ...string) (diceRolls []DiceRoll, errors []error) {
	for i := range rollArgs {
		if diceRoll, err := parseRollArg(rollArgs[i]); diceRoll != nil {
			diceRolls = append(diceRolls, *diceRoll)
		} else {
			errors = append(errors, err)
		}
	}
	return
}

// Parses rollArg. Returns a DiceRoll if valid, an error if invalid.
func parseRollArg(rollArg string) (*DiceRoll, error) {
	// Validate rollArg format
	regex := regexp.MustCompile(`^(\d+)+[dD](\d+)([+-](\d+))?$`)
	if !regex.MatchString(rollArg) {
		return nil, fmt.Errorf("invalid RollArg: %s", rollArg)
	}

	// Parse rollArg into slices using the regex matches
	matches := regex.FindStringSubmatch(rollArg)

	var diceAmmount, diceSize, modifier = 0, 0, 0

	// Parse dice ammount
	if value, argErr := parseRollArgSlice(matches[1]); argErr == nil {
		diceAmmount = value
	} else {
		return nil, argErr
	}
	// Parse dice size
	if value, argErr := parseRollArgSlice(matches[2]); argErr == nil {
		diceSize = value
	} else {
		return nil, argErr
	}
	// Parse modifier
	if len(matches[3]) > 0 {
		if value, argErr := parseRollArgSlice(matches[3]); argErr == nil {
			modifier = value
		} else {
			return nil, argErr
		}
	}

	return NewDiceRoll(diceAmmount, diceSize, modifier)
}

// Parses a rollArg slice. Returns its value if valid, zero and an error if invalid.
func parseRollArgSlice(rollArgSlice string) (int, error) {
	// Validate rollArgSlice size, max allowed is 5 not including minus symbol
	maxAllowedLength := 5
	if strings.ContainsAny(rollArgSlice, "-") {
		maxAllowedLength++
	}

	if len(rollArgSlice) > maxAllowedLength {
		return 0, fmt.Errorf("invalid value: %s. %s", rollArgSlice, bigNumberErrorMsg)
	}

	// Parse rollArgSlice
	return strconv.Atoi(rollArgSlice)
}
