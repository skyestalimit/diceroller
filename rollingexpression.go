package diceroller

import "fmt"

// Represents a sequence of DiceRolls.
type rollingExpression struct {
	diceRolls []DiceRoll
}

// Results of performing a rollingExpression.
type rollingExpressionResult struct {
	rollExpr rollingExpression
	Results  []DiceRollResult
}

// Constructor of rollingExpression.
func newRollingExpression() *rollingExpression {
	return &rollingExpression{make([]DiceRoll, 0)}
}

// Constructor of RollingExpressionResult.
func newRollingExpressionResult(rollExpr rollingExpression) *rollingExpressionResult {
	return &rollingExpressionResult{rollExpr, make([]DiceRollResult, 0)}
}

// Sums multiple RollingExpressionResult.
func RollingExpressionResultSum(rollExprResults ...rollingExpressionResult) (sum int) {
	for e := range rollExprResults {
		sum += DiceRollResultsSum(rollExprResults[e].Results...)
	}
	return
}

// Rolling Expression as String.
func (rollExpr rollingExpression) String() string {
	rollExprStr := ""
	for i := range rollExpr.diceRolls {
		rollExprStr += " " + rollExpr.diceRolls[i].String()
	}
	return rollExprStr
}

// Formatted result output.
func (rollExprResult rollingExpressionResult) String() string {
	resultStr := "Rolling expression: " + rollExprResult.rollExpr.String() + "\n" // add attribs to string
	for i := range rollExprResult.Results {
		resultStr += rollExprResult.Results[i].String()
	}
	resultStr += fmt.Sprintf("Rolling expression sum: %d \n", DiceRollResultsSum(rollExprResult.Results...))
	return resultStr
}

// Detects a critical hit.
func (rollExprResult rollingExpressionResult) detectScoredCritHit() bool {
	critHit := false
	for i := range rollExprResult.Results {
		result := rollExprResult.Results[i]

		if result.hasScoredCritHit() {
			critHit = true
			break
		}
	}
	return critHit
}
