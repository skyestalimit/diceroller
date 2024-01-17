package diceroller

import "fmt"

type rollingExpression struct {
	attribs   attributes
	diceRolls []DiceRoll
}

type rollingExpressionResult struct {
	results []DiceRollResult
}

func newRollingExpression() *rollingExpression {
	rollExpr := new(rollingExpression)
	rollExpr.attribs = newDnDRollAttributes()
	rollExpr.diceRolls = make([]DiceRoll, 0)
	return rollExpr
}

func newRollingExpressionResult() *rollingExpressionResult {
	rollExpr := new(rollingExpressionResult)
	rollExpr.results = make([]DiceRollResult, 0)
	return rollExpr
}

func RollingExpressionResultSum(results ...rollingExpressionResult) (sum int) {
	for i := range results {
		sum += DiceRollResultsSum(results[i].results...)
	}
	return
}

func (rollExpr rollingExpressionResult) String() string {
	resultStr := ""
	for i := range rollExpr.results {
		resultStr += rollExpr.results[i].String()
	}
	resultStr += fmt.Sprintf("\nRoll sum: %d \n", DiceRollResultsSum(rollExpr.results...))
	return resultStr
}
