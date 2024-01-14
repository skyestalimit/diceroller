package diceroller

import "fmt"

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	DiceRollStr string // String representation of the performed DiceRoll, such as 1d6
	attribs     attributes
	Dice        []int // Individual dice roll results
	Dropped     []int
	Sum         int // Sum of Dice array
}

// DiceRollResult constructor with DiceRoll readable string.
func newDiceRollResult(diceRollStr string) *DiceRollResult {
	return &DiceRollResult{diceRollStr, nil, []int{}, []int{}, 0}
}

// Human readable DiceRollResult string.
func (result DiceRollResult) String() string {
	resultStr := "Result of DiceRoll "
	spell := false
	rollAttribs := result.attribs.(*rollAttributes)
	if rollAttribs != nil {
		if rollAttribs.hasAttrib(critAttrib) {
			resultStr += fmt.Sprintf("%s ", critStr)
		}
		if rollAttribs.hasAttrib(spellAttrib) {
			resultStr += fmt.Sprintf("%s ", spellStr)
			spell = true
		}
		if rollAttribs.hasAttrib(advantageAttrib) {
			resultStr += fmt.Sprintf("%s ", advantageStr)
		}
		if rollAttribs.hasAttrib(disadvantageAttrib) {
			resultStr += fmt.Sprintf("%s ", disadvantageStr)
		}
	}

	resultStr += fmt.Sprintf("%s: %s ", result.DiceRollStr, fmt.Sprint(result.Dice))

	if len(result.Dropped) > 0 {
		resultStr += fmt.Sprintf("Dropped: %s ", fmt.Sprint(result.Dropped))
	}

	resultStr += fmt.Sprintf("Sum: %d", result.Sum)

	if spell {
		half := result.Sum / 2
		if half < 1 {
			half = 1
		}
		resultStr += fmt.Sprintf(" Half: %d", half)
	}

	return resultStr
}

// Returns the total sum of a DiceRollResult array.
func DiceRollResultsSum(results ...DiceRollResult) (sum int) {
	for i := range results {
		sum += results[i].Sum
	}
	if len(results) > 0 && sum < 1 {
		sum = 1
	}
	return
}
