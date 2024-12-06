package diceroller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RollArg regex
const rollArgFormat string = `^([+-])?(\d+)?[dD](\d+)([+-](\d+))?$`

// Attributes regex
const rollAttribsFormat string = `^[a-z]+$`

// Maximum allowed RollArg length
const maxAllowedRollArgLength int = 5

// Parses a RollArg array. Returns a DiceRoll array for valid RollArgs, an
// error array for invalid ones.
func ParseRollArgs(rollArgs ...string) (rollingExpressions []rollingExpression, errors []error) {
	// We're building rollingExpressions along with their rollAttributes
	rollExpr := newRollingExpression()
	attribs := newRollAttributes()

	for i := range rollArgs {
		if rollAttrib := checkForRollAttribute(rollArgs[i]); rollAttrib != 0 {
			// Start a new rolling expression after a dice roll sequence ends
			if len(rollExpr.diceRolls) > 0 {
				rollingExpressions = append(rollingExpressions, *rollExpr)
				rollExpr = newRollingExpression()
				attribs = newRollAttributes()
			}
			// Apply the rollAttribute to diceRolls
			attribs.setRollAttrib(rollAttrib)
		} else if diceRoll, err := parseRollArg(rollArgs[i]); err == nil {
			diceRoll.rollAttribs = attribs
			rollExpr.diceRolls = append(rollExpr.diceRolls, *diceRoll)
		} else {
			errors = append(errors, err)
		}
	}

	rollingExpressions = append(rollingExpressions, *rollExpr)

	return
}

// Checks if the rollArg is a rollAttribute. Returns the rollAttribute value if it matches, otherwise zero.
func checkForRollAttribute(rollArg string) rollAttribute {
	var rollAttrib rollAttribute = 0
	attribRegEx := regexp.MustCompile(rollAttribsFormat)
	if attribRegEx.MatchString(strings.ToLower(rollArg)) {
		rollAttrib = rollAttributeMap[rollArg]
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
	var rollAttributes *rollAttributes = newRollAttributes()

	// Parse minus sign
	if len(matches[1]) > 0 {
		if strings.EqualFold(matches[1], "-") {
			rollAttributes.setRollAttrib(minusAttrib)
		}
	}

	// Parse dice ammount
	if len(matches[2]) > 0 {
		if value, argErr := parseRollArgSlice(matches[2]); argErr == nil {
			diceAmmount = value
		} else {
			return nil, argErr
		}
	} else {
		// dY syntax.
		diceAmmount = 1
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

	return NewDiceRollWithAttribs(diceAmmount, diceSize, modifier, rollAttributes)
}

// Parses a rollArg slice. Returns its value if valid, zero and an error if invalid.
func parseRollArgSlice(rollArgSlice string) (int, error) {
	// Validate rollArgSlice size, max allowed is 5 not including minus symbol
	maxAllowedLength := maxAllowedRollArgLength
	if strings.ContainsAny(rollArgSlice, "-") {
		maxAllowedLength++
	}

	if len(rollArgSlice) > maxAllowedLength {
		return 0, fmt.Errorf("invalid value: %s. %s", rollArgSlice, bigNumberErrorMsg)
	}

	// Parse rollArgSlice
	return strconv.Atoi(rollArgSlice)
}
