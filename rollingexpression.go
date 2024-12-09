package diceroller

// Represents a sequence of DiceRolls.
type rollingExpression struct {
	diceRolls []DiceRoll
}

// Constructor of rollingExpression.
func newRollingExpression(diceRolls ...DiceRoll) *rollingExpression {
	return &rollingExpression{append(make([]DiceRoll, 0), diceRolls...)}
}
