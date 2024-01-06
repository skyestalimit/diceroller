package diceroller

import "fmt"

// Dice Roll Result structure
type DiceRollResult struct {
	DiceRollStr string
	Dice        []int
	Sum         int
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
func DiceRollResultsSum(results []DiceRollResult) (sum int) {
	for i := range results {
		sum += results[i].Sum
	}
	return sum
}
