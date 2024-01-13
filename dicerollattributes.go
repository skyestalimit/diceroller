package diceroller

type rollType int

const (
	normalRoll rollType = 1
	crit       rollType = 2
	spell      rollType = 3
)

type rollAttribs struct {
	rollType  rollType
	advantage bool
}
