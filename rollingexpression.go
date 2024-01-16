package diceroller

// To allow adding DiceRolls to a rollingExpression
type rollable interface {
	String() string
	Roll() int
}

type attributes interface {
	setRollAttrib(rollAttribute)
	hasAttrib(rollAttribute) bool
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
