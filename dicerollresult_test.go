package diceroller

import (
	"math/rand"
	"testing"
)

func FuzzDiceRollResultSum(f *testing.F) {
	f.Add(10)
	f.Fuzz(func(t *testing.T, rolls int) {
		results := make([]diceRollResult, 0)
		for i := 0; i < rolls; i++ {
			diceRoll := newDiceRoll(rand.Intn(99999)+1, rand.Intn(99999)+1, rand.Intn(99999)+1)
			result, _ := validateAndperformRoll(*diceRoll)
			results = append(results, *result)
		}

		diceRollResultsSum(results...)
	})
}
