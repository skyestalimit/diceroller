package diceroller

import (
	"fmt"
	"sort"
)

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	DiceRollStr string     // String representation of the performed DiceRoll, such as 1d6
	Attribs     attributes // Contains rollAttributes used to generate the rolls
	Dice        []int      // Individual dice roll result
	Dropped     []int      // Dropped dice roll result
	Sum         int        // Sum of Dice
}

// DiceRollResult constructor with DiceRoll readable string.
func newDiceRollResult(diceRollStr string) *DiceRollResult {
	return newDiceRollResultWithAttribs(diceRollStr, newRollAttributes())
}

// DiceRollResult constructor with DiceRoll readable string and rollAttributes.
func newDiceRollResultWithAttribs(diceRollStr string, attribs *rollAttributes) *DiceRollResult {
	return &DiceRollResult{diceRollStr, attribs, []int{}, []int{}, 0}
}

// Human readable DiceRollResult string.
func (result DiceRollResult) String() string {
	resultStr := "Result of DiceRoll \""
	spell := false

	// Start with roll attributes
	rollAttribsMap := result.Attribs.(*rollAttributes)
	if rollAttribsMap != nil {
		// Sort the attributes
		attribs := make([]rollAttribute, 0, len(rollAttribsMap.attribs))
		for rollAttrib := range rollAttribsMap.attribs {
			attribs = append(attribs, rollAttrib)
		}
		sort.SliceStable(attribs, func(i int, j int) bool {
			return attribs[i] < attribs[j]
		})

		// Add them to the result string
		for i := range attribs {
			rollAttrib := rollAttribute(attribs[i])
			switch rollAttrib {
			case spellAttrib:
				spell = true
			}
			resultStr += fmt.Sprintf("%s ", rollAttributeMap[rollAttrib])
		}
	}

	// DiceRoll string and dice result array
	resultStr += fmt.Sprintf("%s\": \nRolls:   %s\n", result.DiceRollStr, fmt.Sprint(result.Dice))

	// Dropped dice array
	if len(result.Dropped) > 0 {
		resultStr += fmt.Sprintf("Dropped: %s\n", fmt.Sprint(result.Dropped))
	}

	// The DiceRoll sum
	resultStr += fmt.Sprintf("Sum: %d", result.Sum)

	// Half ammount for passing a spell saving throw
	if spell {
		resultStr += fmt.Sprintf(" Save: %d", halve(result.Sum))
	}

	resultStr += "\n"

	return resultStr
}

// Returns the total sum of a DiceRollResult array.
func DiceRollResultsSum(results ...DiceRollResult) (sum int) {
	for i := range results {
		sum += results[i].Sum
	}

	// Minimum DiceRoll result is 1, if at least a die was rolled
	if len(results) > 0 && sum < 1 {
		sum = 1
	}
	return
}
