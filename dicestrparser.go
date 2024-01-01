package roller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const plus = "+"
const minus = "-"

func ParseRollArgs(rollArgs []string) []DiceRoll {
	// Parse the roll args and return an array of DiceRoll
	diceRolls := make([]DiceRoll, 0)
	for i := range rollArgs {
		diceRoll, err := parseRollArg(rollArgs[i])
		if err != nil {
			fmt.Println(err)
		} else {
			diceRolls = append(diceRolls, diceRoll)
		}
	}
	return diceRolls
}

func parseRollArg(rollArg string) (DiceRoll, error) {
	diceRoll := new(DiceRoll)

	// Validate arg format
	rollArg = strings.ToLower(rollArg)
	regExp := regexp.MustCompile("^[[:digit:]]+d[[:digit:]]+([+|-][[:digit:]]+)?$")
	if !regExp.MatchString(rollArg) {
		return *diceRoll, createInvalidRollArgError(rollArg)
	}

	// Parse a roll argument and return it as a DiceRoll if valid
	var argErr error = nil
	rollArg, diceRoll.Modifier, argErr = evaluateModifier(rollArg)
	if argErr != nil {
		return *diceRoll, argErr
	}
	diceRoll.DiceAmmount, diceRoll.DiceSize, argErr = evaluateDiceSizeAndAmmount(rollArg)
	return *diceRoll, argErr
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
	invalidErr := validateArgSize(modSlices[1])
	if invalidErr != nil {
		return rollArg, 0, invalidErr
	}
	mod, modErr := strconv.Atoi(modSlices[1])
	if modErr != nil {
		modErr = fmt.Errorf("error converting modifier %s: \n%s", modSlices[1], modErr.Error())
		mod = 0
	} else {
		// Modifier is valid and processed, remove that slice from the rollArg
		rollArg = modSlices[0]
		if strings.EqualFold(symbol, minus) {
			mod = -mod
		}
	}
	return rollArg, mod, modErr
}

func evaluateDiceSizeAndAmmount(rollArg string) (int, int, error) {
	argSlices := strings.Split(rollArg, "d")
	ammount, diceErr := parseDiceSlice(argSlices[0])
	if diceErr != nil {
		return ammount, 0, diceErr
	}
	size, diceErr := parseDiceSlice(argSlices[1])
	if diceErr != nil {
		return ammount, size, diceErr
	}
	if ammount <= 0 || size <= 1 {
		diceErr = createInvalidRollArgError(rollArg)
	}
	return ammount, size, diceErr
}

func parseDiceSlice(diceSlice string) (int, error) {
	invalidErr := validateArgSize(diceSlice)
	if invalidErr != nil {
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
		return fmt.Errorf("invalid value: %s. \nThis is a dice roller, not a Pi calculator", arg)
	}
	return nil
}

func createInvalidRollArgError(rollArg string) error {
	return fmt.Errorf("invalid roll arg: %s", rollArg)
}
