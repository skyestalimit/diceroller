package diceroller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RollArg regex
const rollArgFormat string = `^([+-])?(\d+)+[dD](\d+)([+-](\d+))?$`

// Attributes regex
const rollAttribsFormat string = `^[a-z]+$`

// Parses a RollArg array. Returns a DiceRoll array for valid RollArgs
// and an error array for invalid ones.
//
// Expected format is XdY[+|-Z].
//
// Valid RollArg examples: "4d4+1", "10d10", "1d6-1", "1D8".
func ParseRollArgs(rollArgs ...string) (rollExpr rollingExpression, errors []error) {
	rollExpr = newRollingExpression()
	for i := range rollArgs {
		if rollAttrib := checkForRollAttribute(rollArgs[i]); rollAttrib > 0 {
			rollExpr.attribs.setRollAttrib(rollAttrib)
		} else if diceRoll, err := parseRollArg(rollArgs[i]); diceRoll != nil {
			rollExpr.diceRolls = append(rollExpr.diceRolls, *diceRoll)
		} else {
			errors = append(errors, err)
		}
	}
	return
}

func checkForRollAttribute(rollArg string) rollAttribute {
	var rollAttrib rollAttribute = 0
	attribRegEx := regexp.MustCompile(rollAttribsFormat)
	if attribRegEx.MatchString(strings.ToLower(rollArg)) {
		switch rollArg {
		case critStr:
			rollAttrib = critAttrib
		case spellStr:
			rollAttrib = spellAttrib
		case advantageStr:
			rollAttrib = advantageAttrib
		case disadvantageStr:
			rollAttrib = disadvantageAttrib
		case dropLowStr:
			rollAttrib = dropLowAttrib
		}
	}
	return rollAttrib
}

// Parses rollArg. Returns a DiceRoll if valid, an error if invalid.
func parseRollArg(rollArg string) (*DiceRoll, error) {
	// Validate rollArg format
	rollArgregex := regexp.MustCompile(rollArgFormat)
	if !rollArgregex.MatchString(rollArg) {
		return nil, fmt.Errorf("invalid RollArg: %s", rollArg)
	}

	// Parse rollArg into slices using the regex matches
	matches := rollArgregex.FindStringSubmatch(rollArg)

	var diceAmmount, diceSize, modifier = 0, 0, 0
	plus := true
	// Parse minus sign
	if len(matches[1]) > 0 {
		if strings.EqualFold(matches[1], "-") {
			plus = false
		}
	}
	// Parse dice ammount
	if value, argErr := parseRollArgSlice(matches[2]); argErr == nil {
		diceAmmount = value
	} else {
		return nil, argErr
	}
	// Parse dice size
	if value, argErr := parseRollArgSlice(matches[3]); argErr == nil {
		diceSize = value
	} else {
		return nil, argErr
	}
	// Parse modifier
	if len(matches[4]) > 0 {
		if value, argErr := parseRollArgSlice(matches[4]); argErr == nil {
			modifier = value
		} else {
			return nil, argErr
		}
	}

	return NewDiceRoll(diceAmmount, diceSize, modifier, plus)
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
