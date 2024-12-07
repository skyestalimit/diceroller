package diceroller

import (
	"math/rand"
	"testing"
)

func FuzzRollResultSum(f *testing.F) {
	f.Add(10, 4)
	f.Fuzz(func(t *testing.T, rolls int, exprs int) {
		rollExprs := make([]rollingExpression, 0)
		for e := 0; e < exprs; e++ {
			rollExpr := newRollingExpression()
			for i := 0; i < rolls; i++ {
				diceRoll := newDiceRoll(rand.Intn(99999)+1, rand.Intn(99999)+1, rand.Intn(99999)+1)
				rollExpr.diceRolls = append(rollExpr.diceRolls, *diceRoll)
			}
			rollExprs = append(rollExprs, *rollExpr)
		}
		rollResults, _ := performRollingExpressions()
		RollResultSum(rollResults...)
	})
}
