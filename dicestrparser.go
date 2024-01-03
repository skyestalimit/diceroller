package roller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const plus = "+"
const minus = "-"

// Parse command line Roll Arguments
// Returns a map of DiceRoll and error
func ParseRollArgs(rollArgs []string) map[*DiceRoll]error {
	diceRollMap := map[*DiceRoll]error{}
	for i := range rollArgs {
		diceRoll, err := ParseRollArg(rollArgs[i])
		diceRollMap[diceRoll] = err
	}
	return diceRollMap
}

// Parse a single Roll Argument
// Returns a DiceRoll
func ParseRollArg(rollArg string) (diceRoll *DiceRoll, argErr error) {
	// Validate arg format
	rollArg = strings.ToLower(rollArg)
	regExp := regexp.MustCompile("^[[:digit:]]+d[[:digit:]]+([+|-][[:digit:]]+)?$")
	if !regExp.MatchString(rollArg) {
		return nil, createInvalidRollArgError(rollArg)
	}

	// Parse a roll argument and return it as a DiceRoll if valid
	diceRoll = new(DiceRoll)
	if rollArg, diceRoll.Modifier, argErr = evaluateModifier(rollArg); argErr != nil {
		return diceRoll, argErr
	}
	diceRoll.DiceAmmount, diceRoll.DiceSize, argErr = evaluateDiceSizeAndAmmount(rollArg)
	return diceRoll, argErr
}

func evaluateModifier(rollArg string) (string, int, error) {
	mod := 0
	var modErr error = nil
	if strings.ContainsAny(rollArg, plus) {
		rollArg, mod, modErr = parseModifier(rollArg, plus)

	} else if strings.ContainsAny(rollArg, minus) {
		rollArg, mod, modErr = parseModifier(rollArg, minus)
	}
	return rollArg, mod, modErr
}

func parseModifier(rollArg string, symbol string) (string, int, error) {
	modSlices := strings.Split(rollArg, symbol)
	if invalidErr := validateArgSize(modSlices[1]); invalidErr != nil {
		return rollArg, 0, invalidErr
	}
	mod, modErr := strconv.Atoi(modSlices[1])
	if modErr == nil {
		// Modifier is valid and processed, remove that slice from the rollArg
		rollArg = modSlices[0]
		if strings.EqualFold(symbol, minus) {
			mod = -mod
		}
	} else {
		modErr = fmt.Errorf("error converting modifier %s: \n%s", modSlices[1], modErr.Error())
		mod = 0
	}
	return rollArg, mod, modErr
}

func evaluateDiceSizeAndAmmount(rollArg string) (ammount int, size int, argErr error) {
	argSlices := strings.Split(rollArg, "d")
	if ammount, argErr = parseDiceSlice(argSlices[0]); argErr != nil {
		return ammount, 0, argErr
	}
	if size, argErr = parseDiceSlice(argSlices[1]); argErr != nil {
		return ammount, size, argErr
	}
	if ammount <= 0 || size <= 1 {
		argErr = createInvalidRollArgError(rollArg)
	}
	return ammount, size, argErr
}

func parseDiceSlice(diceSlice string) (int, error) {
	if invalidErr := validateArgSize(diceSlice); invalidErr != nil {
		return 0, invalidErr
	}
	dice, diceErr := strconv.Atoi(diceSlice)
	if diceErr != nil {
		diceErr = fmt.Errorf("error converting dice value of %s: \n%s", diceSlice, diceErr.Error())
		dice = 0
	}
	return dice, diceErr
}

func validateArgSize(arg string) error {
	if len(arg) > 5 {
		return fmt.Errorf("invalid value: %s. This is a dice roller, not a Pi calculator", arg)
	}
	return nil
}

func createInvalidRollArgError(rollArg string) error {
	return fmt.Errorf("invalid roll arg: %s", rollArg)
}
