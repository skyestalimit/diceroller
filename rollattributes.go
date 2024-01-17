package diceroller

import "strings"

type rollAttribute int

// rollAttribute values. 0 is invalid.
const (
	// DnD rollAttribute
	hitAttrib          rollAttribute = iota + 1
	dmgAttrib          rollAttribute = iota + 1
	critAttrib         rollAttribute = iota + 1
	spellAttrib        rollAttribute = iota + 1
	halfAttrib         rollAttribute = iota + 1
	advantageAttrib    rollAttribute = iota + 1
	disadvantageAttrib rollAttribute = iota + 1
	dropHighAttrib     rollAttribute = iota + 1
	dropLowAttrib      rollAttribute = iota + 1
)

// Allowed rollAttribute string as RollArg.
const (
	// DnD rollAttribute strings
	hitStr          string = "hit"
	dmgStr          string = "dmg"
	critStr         string = "crit"
	spellStr        string = "spell"
	halfStr         string = "half"
	advantageStr    string = "adv"
	disadvantageStr string = "dis"
	dropHighStr     string = "drophigh"
	dropLowStr      string = "droplow"
)

var rollAttributeMap = map[rollAttribute]string{
	// DnD rollAttribute map
	hitAttrib:          hitStr,
	dmgAttrib:          dmgStr,
	critAttrib:         critStr,
	spellAttrib:        spellStr,
	halfAttrib:         halfStr,
	advantageAttrib:    advantageStr,
	disadvantageAttrib: disadvantageStr,
	dropHighAttrib:     dropHighStr,
	dropLowAttrib:      dropLowStr,
}

type attributes interface {
	setRollAttrib(rollAttribute)
	hasAttrib(rollAttribute) bool
}

type rollAttributes struct {
	attribs map[rollAttribute]bool
}

// Constructor for rollAttributes.
func newRollAttributes() *rollAttributes {
	newRollAttributes := new(rollAttributes)
	newRollAttributes.attribs = make(map[rollAttribute]bool)
	return newRollAttributes
}

// To retrieve the roleAttribute matching wanted roleAttribute string.
func rollAttributeMapKey(attribMap map[rollAttribute]string, wanted string) rollAttribute {
	for attrib, attribStr := range attribMap {
		if strings.EqualFold(attribStr, wanted) {
			return attrib
		}
	}
	return 0
}

// Sets attrib to true and prevents rollAttribute incompatibilities.
func (rollAttribs rollAttributes) setRollAttrib(attrib rollAttribute) {
	switch attrib {
	case advantageAttrib:
		delete(rollAttribs.attribs, disadvantageAttrib)
	case disadvantageAttrib:
		delete(rollAttribs.attribs, advantageAttrib)
	}
	rollAttribs.attribs[attrib] = true
}

// Returns true if wanted is set.
func (attrib rollAttributes) hasAttrib(wanted rollAttribute) bool {
	return attrib.attribs[wanted]
}
