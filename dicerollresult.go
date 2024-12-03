package diceroller

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"
)

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	diceRoll      DiceRoll // Performed DiceRoll
	dice          []int    // Individual dice roll result
	sum           int      // Sum of Dice
	advDisDropped []int    // Dropped advantage/disadvantage dice
	highDropped   []int    // Dropped high dice
	lowDropped    []int    // Dropped low dice
}

// DiceRollResult constructor with DiceRoll readable string and rollAttributes.
func newDiceRollResult(diceRoll DiceRoll) *DiceRollResult {
	return &DiceRollResult{diceRoll, []int{}, 0, []int{}, []int{}, []int{}}
}

// Returns the total sum of a DiceRollResult array.
func DiceRollResultsSum(results ...DiceRollResult) (sum int) {
	for i := range results {
		sum += results[i].sum
	}

	return
}

// Detects a critical hit.
func (rollResult DiceRollResult) hasScoredCritHit() bool {
	critHit := false

	if len(rollResult.dice) == 1 && rollResult.dice[0] == 20 &&
		rollResult.diceRoll.diceAmmount == 1 && rollResult.diceRoll.diceSize == 20 {
		critHit = true
	}

	return critHit
}

// Human readable DiceRollResult string.
func (result DiceRollResult) String() string {
	resultStr := " Result of DiceRoll \""
	advDisStr := ""
	spell := false

	// Start with roll attributes, add them as String
	if result.diceRoll.rollAttribs != nil {
		// Sort the attributes
		attribs := maps.Keys(result.diceRoll.rollAttribs.attribs)

		sort.SliceStable(attribs, func(i int, j int) bool {
			return attribs[i] < attribs[j]
		})

		// Add them to the result string
		for i := range attribs {
			rollAttrib := rollAttribute(attribs[i])
			// Collect that affects the results printout later
			switch rollAttrib {
			case advantageAttrib:
				advDisStr = "Adv drop:"
			case disadvantageAttrib:
				advDisStr = "Dis drop:"
			case spellAttrib:
				spell = true
			}
			resultStr += fmt.Sprintf("%s ", rollAttributeMapKey(rollAttributeMap, rollAttrib))
		}
	}

	// DiceRoll string and dice result array
	resultStr += fmt.Sprintf("%s\": \n  Rolls:     %s\n", result.diceRoll, fmt.Sprint(result.dice))

	// Advantage / disadvantage dropped dice array
	if len(result.advDisDropped) > 0 {
		resultStr += fmt.Sprintf("  %s  %s\n", advDisStr, fmt.Sprint(result.advDisDropped))
	}

	// Dropped High dice array
	if len(result.highDropped) > 0 {
		resultStr += fmt.Sprintf("  Drop High: %s\n", fmt.Sprint(result.highDropped))
	}

	// Dropped Low dice array
	if len(result.lowDropped) > 0 {
		resultStr += fmt.Sprintf("  Drop Low:  %s\n", fmt.Sprint(result.lowDropped))
	}

	// The DiceRoll sum
	resultStr += fmt.Sprintf("  Sum:       %d", result.sum)

	// Half ammount for passing a spell saving throw
	if spell {
		resultStr += fmt.Sprintf(" Saved: %d", halve(result.sum))
	}

	resultStr += "\n"

	return resultStr
}
