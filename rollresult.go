package diceroller

import "fmt"

// Results of performing a rollingExpression.
type rollResult struct {
	results []diceRollResult
}

// Constructor of rollResult.
func newRollResult() *rollResult {
	return &rollResult{make([]diceRollResult, 0)}
}

// Sums multiple rollResult.
func RollResultsSum(rollResults ...rollResult) (sum int) {
	for e := range rollResults {
		sum += rollResults[e].Sum()
	}
	return
}

func (rollResult rollResult) Sum() int {
	return diceRollResultsSum(rollResult.results...)
}

// Formatted result output.
func (rollResult rollResult) String() string {
	resultStr := "Roll result : \n" // add attribs to string
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
