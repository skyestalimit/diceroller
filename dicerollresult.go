package diceroller

import "fmt"

// A DiceRollResult contains the results of performing a DiceRoll
type DiceRollResult struct {
	DiceRollStr string // String representation of the performed DiceRoll, such as 1d6
	Dice        []int  // Individual dice roll results
	Sum         int    // Sum of Dice array
}

// DiceRollResult constructor with DiceRoll readable string.
func NewDiceRollResult(diceRollStr string) *DiceRollResult {
	return &DiceRollResult{diceRollStr, []int{}, 0}
}

// Human readable DiceRollResult string.
func (result DiceRollResult) String() string {
	return fmt.Sprintf("Result of DiceRoll %s: %s Sum: %d", result.DiceRollStr, fmt.Sprint(result.Dice), result.Sum)
}

// Returns the total sum of a DiceRollResult array.
func DiceRollResultsSum(results ...DiceRollResult) (sum int) {
	for i := range results {
		sum += results[i].Sum
	}
	return
}
