package diceroller

import (
	"fmt"
	"sort"
)

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	DiceRollStr   string     // String representation of the performed DiceRoll, such as 1d6
	Attribs       attributes // Contains rollAttributes used to generate the rolls
	Dice          []int      // Individual dice roll result
	Sum           int        // Sum of Dice
	AdvDisDropped []int      // Dropped advantage/disadvantage dice
	HighDropped   []int      // Dropped high dice
	LowDropped    []int      // Dropped low dice
}

// DiceRollResult constructor with DiceRoll readable string.
func newDiceRollResult(diceRollStr string) *DiceRollResult {
	return newDiceRollResultWithAttribs(diceRollStr, newRollAttributes())
}

// DiceRollResult constructor with DiceRoll readable string and rollAttributes.
func newDiceRollResultWithAttribs(diceRollStr string, attribs *rollAttributes) *DiceRollResult {
	return &DiceRollResult{diceRollStr, attribs, []int{}, 0, []int{}, []int{}, []int{}}
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

// Human readable DiceRollResult string.
func (result DiceRollResult) String() string {
	resultStr := "Result of DiceRoll \""
	var advDis rollAttribute = 0
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
			case advantageAttrib:
				advDis = advantageAttrib
			case disadvantageAttrib:
				advDis = disadvantageAttrib
			case spellAttrib:
				spell = true
			}
			resultStr += fmt.Sprintf("%s ", rollAttributeMap[rollAttrib])
		}
	}

	// DiceRoll string and dice result array
	resultStr += fmt.Sprintf("%s\": \n Rolls:     %s\n", result.DiceRollStr, fmt.Sprint(result.Dice))

	// Advantage / disadvantage dropped dice array
	if len(result.AdvDisDropped) > 0 {
		advDisStr := ""
		switch advDis {
		case advantageAttrib:
			advDisStr = "Adv drop:"
		case disadvantageAttrib:
			advDisStr = "Dis drop:"
		}

		resultStr += fmt.Sprintf(" %s  %s\n", advDisStr, fmt.Sprint(result.AdvDisDropped))
	}

	// Dropped High dice array
	if len(result.HighDropped) > 0 {
		resultStr += fmt.Sprintf(" Drop High: %s\n", fmt.Sprint(result.HighDropped))
	}

	// Dropped Low dice array
	if len(result.LowDropped) > 0 {
		resultStr += fmt.Sprintf(" Drop Low:  %s\n", fmt.Sprint(result.LowDropped))
	}

	// The DiceRoll sum
	resultStr += fmt.Sprintf(" Sum:       %d", result.Sum)

	// Half ammount for passing a spell saving throw
	if spell {
		resultStr += fmt.Sprintf(" Saved: %d", halve(result.Sum))
	}

	resultStr += "\n"

	return resultStr
}
