package diceroller

type rollable interface {
}

type rollingExpression struct {
	attribs   attributes
	diceRolls []rollable
}

func newRollingExpression() rollingExpression {
	rollExpr := new(rollingExpression)
	rollExpr.attribs = newRollAttributes()
	rollExpr.diceRolls = make([]rollable, 0)
	return *rollExpr
}
