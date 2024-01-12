package diceroller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Max allowed arg slice length
const maxRollArgSliceSize = 5

// Plus symbol character
const plusSymbol = "+"

// Minus symbol character
const minusSymbol = "-"

// Parses a RollArg array. Returns a DiceRoll array for valid RollArgs
// and an error array for invalid ones.
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
// Expected format is XdY[+|-Z].
//
// Valid RollArg examples: "4d4+1", "10d10", "1d6-1", "1D8".
func parseRollArg(rollArg string) (*DiceRoll, error) {
	// Validate RollArg format
	if !regexp.MustCompile("^[[:digit:]]+[d|D][[:digit:]]+([+|-][[:digit:]]+)?$").MatchString(rollArg) {
		// Invalid RollArg, return error
		return nil, fmt.Errorf("invalid RollArg: %s", rollArg)
	}

	// Parse the RollArg
	if rollArg, modifier, argErr := evaluateModifier(rollArg); argErr != nil {
		// Invalid modifier
		return nil, argErr
	} else if diceAmmount, diceSize, argErr := evaluateDiceSizeAndAmmount(rollArg); argErr != nil {
		// Invalid dice values
		return nil, argErr
	} else {
		// Parsed RollArgs, their values gets validated in DiceRoll constructor
		return NewDiceRoll(diceAmmount, diceSize, modifier)
	}
}

// Evaluates rollArg modifier if present.
//
//	-If present, returns parseModifier function call results.
//	-If not present, returns rollArg, a zero and no error.
func evaluateModifier(rollArg string) (string, int, error) {
	// Detect modifier
	if strings.ContainsAny(rollArg, plusSymbol) {
		return parseModifier(rollArg, plusSymbol)

	} else if strings.ContainsAny(rollArg, minusSymbol) {
		return parseModifier(rollArg, minusSymbol)
	}

	// Modifier not present in rollArg
	return rollArg, 0, nil
}

// Parses and validates the size of rollArg modifier.
//
//	-If valid, returns rollArg with the modifier part removed, the modifier value and no error.
//	-If invalid an unexpected parsing error happens, returns rollArg, a zero and an error.
func parseModifier(rollArg string, symbol string) (string, int, error) {
	// Extract the modifier from rollArg and validate its size
	rollArgSlices := strings.Split(rollArg, symbol)

	// Parse the modifier
	if modifier, argErr := parseRollArgSlice(rollArgSlices[1]); argErr == nil {
		if strings.EqualFold(symbol, minusSymbol) {
			// Make the modifier value negative for minus modifiers
			modifier = -modifier
		}
		// Modifier is valid and processed, return a rollArg slice without the modifier part
		return rollArgSlices[0], modifier, nil
	} else {
		// Invalid modifier
		return rollArg, 0, argErr
	}
}

// Evaluates rollArg. Returns its values if a valid, zeroes and an error if invalid.
func evaluateDiceSizeAndAmmount(rollArg string) (int, int, error) {
	// rollArg is either xdy or xDy format by now. To lowercase and split on "d"
	rollArgSlices := strings.Split(strings.ToLower(rollArg), "d")
	if diceAmmount, argErr := parseRollArgSlice(rollArgSlices[0]); argErr != nil {
		// Invalid ammount
		return 0, 0, argErr
	} else if diceSize, argErr := parseRollArgSlice(rollArgSlices[1]); argErr != nil {
		// Invalide size
		return 0, 0, argErr
	} else {
		// Valid rollArg
		return diceAmmount, diceSize, nil
	}
}

// Parses a rollArg slice. Returns its value if valid, zero and an error if invalid.
func parseRollArgSlice(rollArgSlice string) (int, error) {
	// Validate rollArgSlice size
	if argErr := validateRollArgSliceSize(rollArgSlice); argErr != nil {
		return 0, argErr
	}

	// Parse rollArgSlice
	if diceValue, diceErr := strconv.Atoi(rollArgSlice); diceErr == nil {
		return diceValue, nil
	} else {
		return 0, fmt.Errorf("error parsing value of %s: \n%s", rollArgSlice, diceErr.Error())
	}
}

// Validates that rollArgSlice is shorter than defined maxRollArgSliceSize.
func validateRollArgSliceSize(rollArgSlice string) error {
	if len(rollArgSlice) > maxRollArgSliceSize {
		return fmt.Errorf("invalid value: %s. %s", rollArgSlice, bigNumberErrorMsg)
	}
	return nil
}
