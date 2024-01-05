package roller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Max allow
const maxRollArgSliceSize = 5
const plus = "+"
const minus = "-"

// Parse command line Roll Arguments. Returns arrays of DiceRolls and errors.
func ParseRollArgs(rollArgs []string) (diceRolls []DiceRoll, errors []error) {
	for i := range rollArgs {
		diceRoll, err := ParseRollArg(rollArgs[i])
		if diceRoll != nil {
			diceRolls = append(diceRolls, *diceRoll)
		} else {
			errors = append(errors, err)
		}
	}
	return diceRolls, errors
}

// Parse and validates a roll arg. Returns a DiceRoll if valid, an error if invalid.
func ParseRollArg(rollArg string) (diceRoll *DiceRoll, argErr error) {
	// Validate arg format
	rollArg = strings.ToLower(rollArg)
	regExp := regexp.MustCompile("^[[:digit:]]+d[[:digit:]]+([+|-][[:digit:]]+)?$")
	if !regExp.MatchString(rollArg) {
		// Invalid roll arg, return error
		return nil, createInvalidRollArgError(rollArg)
	}

	// Parse a roll argument and return it as a DiceRoll if valid. Return error if invalid
	diceRoll = new(DiceRoll)
	if rollArg, diceRoll.Modifier, argErr = evaluateModifier(rollArg); argErr != nil {
		return nil, argErr
	}
	if diceRoll.DiceAmmount, diceRoll.DiceSize, argErr = evaluateDiceSizeAndAmmount(rollArg); argErr != nil {
		return nil, argErr
	}
	return diceRoll, nil
}

// Evaluates if a roll arg modifier is present.
//   - If present, returns parseAndValidateModifier function call results.
//   - If not present, returns the original roll arg, a zero value and no error.
func evaluateModifier(rollArg string) (string, int, error) {
	// Detect if a modifier is present in roll arg. If present, parse and validate
	if strings.ContainsAny(rollArg, plus) {
		return parseAndValidateModifier(rollArg, plus)

	} else if strings.ContainsAny(rollArg, minus) {
		return parseAndValidateModifier(rollArg, minus)
	}
	// Modifier not present in roll arg
	return rollArg, 0, nil
}

// Parses and validate a roll arg modifier.
//   - If valid, returns the roll arg string with the modifier part removed, the modifier value and nil value error.
//   - If invalid, returns the original roll arg, a zero and an error.
//   - For unexpected parsing errors, returns the original roll arg, a zero and an error.
func parseAndValidateModifier(rollArg string, symbol string) (string, int, error) {
	// Extract modifier value and validate its size
	modSlices := strings.Split(rollArg, symbol)
	if invalidErr := validateRollArgSliceSize(modSlices[1]); invalidErr != nil {
		// Invalid modifier value
		return rollArg, 0, invalidErr
	}

	// Parse a roll arg modifier and return its value if valid. Return zero and error if invalid
	if modifier, modErr := strconv.Atoi(modSlices[1]); modErr == nil {
		if strings.EqualFold(symbol, minus) {
			// Make the modifier value negative for minus modifier
			modifier = -modifier
		}
		// Modifier is valid and processed, return rollArg without the modifier part
		return modSlices[0], modifier, modErr
	} else {
		// Unexpected parsing error
		return rollArg, 0, fmt.Errorf("error converting modifier %s: \n%s", modSlices[1], modErr.Error())
	}
}

// Validates, parses and return values of a valid roll arg. Returns an error for an invalid roll arg.
func evaluateDiceSizeAndAmmount(rollArg string) (ammount int, size int, argErr error) {
	// Slice the rollArg into slices and validate values
	argSlices := strings.Split(rollArg, "d")
	if ammount, argErr = parseRollArgSlice(argSlices[0]); argErr != nil {
		return ammount, 0, argErr
	} else if argErr = validateDiceAmmout(ammount); argErr != nil {
		return ammount, 0, argErr
	}
	if size, argErr = parseRollArgSlice(argSlices[1]); argErr != nil {
		return ammount, size, argErr
	} else if argErr = validateDiceSize(size); argErr != nil {
		return ammount, size, argErr
	}

	// Returning valid values
	return ammount, size, nil
}

// Validates and parses a roll arg slice. Returns its value if valid, zero and an error if invalid.
func parseRollArgSlice(argSlice string) (int, error) {
	// Validate arg slice
	if invalidErr := validateRollArgSliceSize(argSlice); invalidErr != nil {
		return 0, invalidErr
	}
	// Parse and return arg slice value
	if diceValue, diceErr := strconv.Atoi(argSlice); diceErr != nil {
		return 0, fmt.Errorf("error converting dice value of %s: \n%s", argSlice, diceErr.Error())
	} else {
		return diceValue, diceErr
	}
}

// Validates that arg slice is shorter than defined maxRollArgSliceSize to avoid overflow and long processing time.
func validateRollArgSliceSize(rollArg string) error {
	if len(rollArg) > maxRollArgSliceSize {
		return fmt.Errorf("invalid value: %s. This is a dice roller, not a Pi calculator", rollArg)
	}
	return nil
}

// Returns a formatter error for invalid roll arg.
func createInvalidRollArgError(rollArg string) error {
	return fmt.Errorf("invalid roll arg: %s", rollArg)
}
