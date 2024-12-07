package diceroller

// Represents a sequence of DiceRolls.
type rollingExpression struct {
	diceRolls []DiceRoll
}

// Constructor of rollingExpression.
func newRollingExpression() *rollingExpression {
	return &rollingExpression{make([]DiceRoll, 0)}
}

// Rolling Expression as String.
func (rollExpr rollingExpression) String() string {
	rollExprStr := ""
	for i := range rollExpr.diceRolls {
		rollExprStr += " " + rollExpr.diceRolls[i].String()
	}
	return rollExprStr
}
