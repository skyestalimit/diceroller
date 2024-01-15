package diceroller

import "strings"

type rollAttribute int

const (
	critAttrib         rollAttribute = iota + 1
	spellAttrib        rollAttribute = iota + 1
	halfAttrib         rollAttribute = iota + 1
	advantageAttrib    rollAttribute = iota + 1
	disadvantageAttrib rollAttribute = iota + 1
	dropHighAttrib     rollAttribute = iota + 1
	dropLowAttrib      rollAttribute = iota + 1
)

const (
	critStr         string = "crit"
	spellStr        string = "spell"
	halfStr         string = "half"
	advantageStr    string = "adv"
	disadvantageStr string = "dis"
	dropHighStr     string = "drophigh"
	dropLowStr      string = "droplow"
)

var rollAttributeMap = map[rollAttribute]string{
	critAttrib:         critStr,
	spellAttrib:        spellStr,
	halfAttrib:         halfStr,
	advantageAttrib:    advantageStr,
	disadvantageAttrib: disadvantageStr,
	dropHighAttrib:     dropHighStr,
	dropLowAttrib:      dropLowStr,
}

func rollAttributeMapKey(attribMap map[rollAttribute]string, wanted string) (rollAttribute, bool) {
	for attrib, attribStr := range attribMap {
		if strings.EqualFold(attribStr, wanted) {
			return attrib, true
		}
	}
	return 0, false
}

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

func (rollAttribs rollAttributes) setRollAttrib(attrib rollAttribute) {
	// Prevent attrib incompatibilities.
	switch attrib {
	case advantageAttrib:
		delete(rollAttribs.attribs, disadvantageAttrib)
		delete(rollAttribs.attribs, dropHighAttrib)
		delete(rollAttribs.attribs, dropLowAttrib)
	case disadvantageAttrib:
		delete(rollAttribs.attribs, advantageAttrib)
		delete(rollAttribs.attribs, dropHighAttrib)
		delete(rollAttribs.attribs, dropLowAttrib)
	case dropLowAttrib:
		delete(rollAttribs.attribs, advantageAttrib)
		delete(rollAttribs.attribs, disadvantageAttrib)
		delete(rollAttribs.attribs, dropHighAttrib)
	case dropHighAttrib:
		delete(rollAttribs.attribs, advantageAttrib)
		delete(rollAttribs.attribs, disadvantageAttrib)
		delete(rollAttribs.attribs, dropLowAttrib)
	}
	rollAttribs.attribs[attrib] = true
}

func (dndAttrib rollAttributes) hasAttrib(attrib rollAttribute) bool {
	return dndAttrib.attribs[attrib]
}
