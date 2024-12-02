package diceroller

import "strings"

type rollAttribute int

// rollAttribute values. 0 is invalid.
const (
	rollAttrib rollAttribute = iota + 1
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
	minusAttrib        rollAttribute = iota + 1
)

// Allowed rollAttribute string as RollArg.
const (
	rollStr string = "roll"
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
	minusAttribStr  string = "minus"
)

var rollAttributeMap = map[rollAttribute]string{
	rollAttrib: rollStr,
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
	minusAttrib:        minusAttribStr,
}

type attributes interface {
	setRollAttrib(rollAttribute)
	hasAttrib(rollAttribute) bool
}

type rollAttributes struct {
	attribs map[rollAttribute]bool
}

// Constructor for rollAttributes.
func newRollAttributes(rollAttribs ...rollAttribute) *rollAttributes {
	newRollAttributes := new(rollAttributes)
	newRollAttributes.attribs = make(map[rollAttribute]bool)
	newRollAttributes.setRollAttrib(rollAttribs...)
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
func (dndAttribs *rollAttributes) setRollAttrib(rollAttribs ...rollAttribute) {
	for i := range rollAttribs {
		switch rollAttribs[i] {
		case advantageAttrib:
			delete(dndAttribs.attribs, disadvantageAttrib)
		case disadvantageAttrib:
			delete(dndAttribs.attribs, advantageAttrib)
		}
		dndAttribs.attribs[rollAttribs[i]] = true
	}
}

// Returns true if wanted is set.
func (dndAttrib *rollAttributes) hasAttrib(wanted rollAttribute) bool {
	return dndAttrib.attribs[wanted]
}

func (dndAttrib *rollAttributes) isCrit() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(critAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isAdvantage() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(advantageAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isDisadvantage() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(disadvantageAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isDropHigh() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(dropHighAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isDropLow() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(dropLowAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isHalf() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(halfAttrib)
	}
	return false
}

func (dndAttrib *rollAttributes) isMinus() bool {
	if dndAttrib != nil {
		return dndAttrib.hasAttrib(minusAttrib)
	}
	return false
}
