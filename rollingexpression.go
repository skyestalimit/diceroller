package diceroller

import "fmt"

// Represents a sequence of DiceRolls.
type rollingExpression struct {
	diceRolls []DiceRoll
}

// Results of performing a rollingExpression.
type RollingExpressionResult struct {
	Results []DiceRollResult
	Sum     int
}

// Constructor of rollingExpression.
func newRollingExpression() *rollingExpression {
	rollExpr := new(rollingExpression)
	rollExpr.diceRolls = make([]DiceRoll, 0)
	return rollExpr
}

// Constructor of RollingExpressionResult.
func newRollingExpressionResult() *RollingExpressionResult {
	rollExpr := new(RollingExpressionResult)
	rollExpr.Sum = 0
	rollExpr.Results = make([]DiceRollResult, 0)
	return rollExpr
}

// Sums multiple RollingExpressionResult.
func RollingExpressionResultSum(results ...RollingExpressionResult) (sum int) {
	for i := range results {
		sum += results[i].Sum
	}
	return
}

// Formatted result output.
func (rollExpr RollingExpressionResult) String() string {
	resultStr := "Rolling expression\n"
	for i := range rollExpr.Results {
		resultStr += rollExpr.Results[i].String()
	}
	resultStr += fmt.Sprintf("\nRoll sum: %d \n", DiceRollResultsSum(rollExpr.Results...))
	return resultStr
}
