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
func RollResultsSum(rollExprResults ...rollResult) (sum int) {
	for e := range rollExprResults {
		sum += rollExprResults[e].Sum()
	}
	return
}

func (rollResult rollResult) Sum() int {
	return diceRollResultsSum(rollResult.results...)
}

// Formatted result output.
func (rollResult rollResult) String() string {
	resultStr := "Roll result : " + rollResult.rollExpr.String() + "\n" // add attribs to string
	for i := range rollResult.results {
		resultStr += rollResult.results[i].String()
	}
	resultStr += fmt.Sprintf("Roll results sum: %d \n", diceRollResultsSum(rollResult.results...))
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
