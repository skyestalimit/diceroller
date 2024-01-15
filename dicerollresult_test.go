package diceroller

import "testing"

func FuzzDiceRollResultSum(f *testing.F) {
	f.Add("2d6+1")
	f.Fuzz(func(t *testing.T, diceStr string) {
		newDiceRollResult(diceStr)
		DiceRollResultsSum(*newDiceRollResult(diceStr))
	})
}
