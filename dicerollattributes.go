package diceroller

type rollAttribute int

const (
	critAttrib         rollAttribute = 2
	spellAttrib        rollAttribute = 3
	advantageAttrib    rollAttribute = 4
	disadvantageAttrib rollAttribute = 5
	dropLowAttrib      rollAttribute = 6
)

const (
	critStr         string = "crit"
	spellStr        string = "spell"
	advantageStr    string = "adv"
	disadvantageStr string = "dis"
	dropLowStr      string = "dropLow"
)

type attributes interface {
	setRollAttrib(rollAttribute)
	hasAttrib(rollAttribute) bool
}

type rollAttributes struct {
	attribs map[rollAttribute]bool
}

func newRollAttributes() *rollAttributes {
	newRollAttributes := new(rollAttributes)
	newRollAttributes.attribs = make(map[rollAttribute]bool)
	return newRollAttributes
}

func (attribs rollAttributes) setRollAttrib(attrib rollAttribute) {
	attribs.attribs[attrib] = true
}

func (dndAttrib rollAttributes) hasAttrib(attrib rollAttribute) bool {
	return dndAttrib.attribs[attrib]
}
