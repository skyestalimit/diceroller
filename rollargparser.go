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
//
// A valid RollArg matches either the DiceRoll format or a rollAttribute string.
//
// DiceRoll format: [X]dY[+|-]Z. In other words XdY or dY followed by + or - Z.
// Valid DiceRoll examples: "5d6", "d20", "4d4+1", "10d10", "1d6-1", "1D8".
//
// RollArgs will be converted into DiceRolls and grouped together in rolling
// expressions. A new rolling expression starts when a rollAttribute is
// parsed after a sequence of DiceRolls.
//
// For example:
//   - 5d6 1d6 10d10 would return one rolling expression containing all three
//     RollArgs as DiceRolls.
//   - adv d20+5 half spell 8d8 would return two rolling expressions: "adv 1d20+5"
//     and "half spell 8d8".
//   - 5d6 roll 4d6 roll 3d6 would return three rolling expressions: "5d6", "4d6" and "3d6".
//
// rollAttribute string list:
//   - roll, hit, dmg : separators, starts a new rolling expressions
//
// DnD rollAttribute string list:
//   - crit: Critical, doubles all dice ammount
//   - spell: Spell, DiceRollResults.String() prints the sum and the sum halved for saves
//   - half: Halves the sums, for resistances and such
//   - adv: Advantage, rolls each dice twice and drops the lowest
//   - dis: Disadvantage, rolls each dice twice and drops the highest
//   - drophigh: Drop High, drops the highest result of a DiceRoll
//   - droplow: Drop Low, drops the lowest result of a DiceRoll
func ParseRollArgs(rollArgs ...string) (rollingExpressions []rollingExpression, errors []error) {
	// We're building rollingExpressions along with their rollAttributes
	rollExpr := newRollingExpression()
	attribs := newDnDRollAttributes()
	diceRollSequence := false

	for i := range rollArgs {
		if rollAttrib := checkForRollAttribute(rollArgs[i]); rollAttrib != 0 {
			if diceRollSequence {
				// Start a new rolling expression after a dice roll sequence ends
				rollingExpressions = append(rollingExpressions, *rollExpr)
				rollExpr = newRollingExpression()
				attribs = newDnDRollAttributes()
			}
			attribs.setRollAttrib(rollAttrib)
			diceRollSequence = false
		} else if diceRoll, err := parseRollArg(rollArgs[i]); err == nil {
			diceRoll.Attribs = attribs
			rollExpr.diceRolls = append(rollExpr.diceRolls, *diceRoll)
			diceRollSequence = true
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
		rollAttrib = rollAttributeMapKey(rollAttributeMap, rollArg)
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

	return NewDiceRoll(diceAmmount, diceSize, modifier, plus)
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
