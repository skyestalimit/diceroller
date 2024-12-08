package diceroller

import "fmt"

// Results of performing a rollingExpression.
type rollResult struct {
	rollExpr rollingExpression
	results  []diceRollResult
}

// Constructor of RollingExpressionResult.
func newRollResult(rollExpr rollingExpression) *rollResult {
	return &rollResult{rollExpr, make([]diceRollResult, 0)}
}

// Sums multiple RollingExpressionResult.
func RollResultSum(rollExprResults ...rollResult) (sum int) {
	for e := range rollExprResults {
		sum += diceRollResultsSum(rollExprResults[e].results...)
	}
	return
}

// Formatted result output.
func (rollExprResult rollResult) String() string {
	resultStr := "Rolling expression: " + rollExprResult.rollExpr.String() + "\n" // add attribs to string
	for i := range rollExprResult.results {
		resultStr += rollExprResult.results[i].String()
	}
	resultStr += fmt.Sprintf("Rolling expression sum: %d \n", diceRollResultsSum(rollExprResult.results...))
	return resultStr
}

// Detects a critical hit.
func (rollResult rollResult) detectScoredCritHit() bool {
	critHit := false
	for i := range rollResult.results {
		result := rollResult.results[i]

		if result.hasScoredCritHit() {
			critHit = true
			break
		}
	}
	return critHit
}
