package diceroller

import (
	"fmt"
	"math"
	"sort"
)

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	DiceRollStr string // String representation of the performed DiceRoll, such as 1d6
	Attribs     attributes
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
	resultStr := "Result of DiceRoll \""
	half := false

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
			case halfAttrib:
				half = true
			}
			resultStr += fmt.Sprintf("%s ", rollAttributeMap[rollAttrib])
		}
	}

	resultStr += fmt.Sprintf("%s\": %s ", result.DiceRollStr, fmt.Sprint(result.Dice))

	if len(result.Dropped) > 0 {
		resultStr += fmt.Sprintf("Dropped: %s ", fmt.Sprint(result.Dropped))
	}

	resultStr += fmt.Sprintf("Sum: %d", result.Sum)

	if half {
		resultStr += fmt.Sprintf(" Half: %d", halve(result.Sum))
	}

	return resultStr
}

// Returns the total sum of a DiceRollResult array.
func DiceRollResultsSum(results ...DiceRollResult) (sum int) {
	for i := range results {
		toAdd := results[i].Sum
		if results[i].Attribs.hasAttrib(halfAttrib) {
			toAdd = halve(toAdd)
		}
		sum += toAdd
	}
	if len(results) > 0 && sum < 1 {
		sum = 1
	}
	return
}

func halve(result int) int {
	halved := 0
	minus := false
	if result < 0 {
		minus = true
	}
	halved = int(math.Abs(float64(result))) / 2
	if halved < 1 {
		halved = 1
	}
	if minus {
		halved = -halved
	}
	return halved
}
