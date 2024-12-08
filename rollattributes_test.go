package diceroller

import (
	"math/rand"
	"testing"
)

func TestRollAttributes(t *testing.T) {
	rollAttribs := newRollAttributes()
	for i := range rollAttributeMap {
		rollAttrib := rollAttributeMap[i]
		rollAttribs.setRollAttrib(rollAttrib)

		if rollAttribs.hasAttrib(rollAttrib) == false {
			t.Fatalf("hasAttrib %d returned false, wanted true", rollAttrib)
		}

		checkForAttribCompatibility(rollAttribs, t)
	}
}

func TestPerformRollWithRollAttributes(t *testing.T) {
	rollAttribs := newRollAttributes()
	for i := range rollAttributeMap {
		rollAttrib := rollAttributeMap[i]

		rollAttribs.setRollAttrib(rollAttrib)

		diceRoll := validDiceRollsValues[rand.Intn(len(validDiceRollsValues))].diceRoll
		diceRoll.rollAttribs = rollAttribs

		if result, diceErr := validateAndperformRoll(diceRoll); diceErr != nil {
			t.Fatalf("DiceRoll %s returned error: %s", diceRoll.String(), diceErr.Error())
		} else if result.sum == 0 {
			t.Fatalf("DiceRoll %s result = %d, wanted > 0", diceRoll.String(), result.sum)
		}
	}
}

func checkForAttribCompatibility(rollAttribs *rollAttributes, t *testing.T) {
	for rollAttrib := range rollAttribs.attribs {
		switch rollAttrib {
		case advantageAttrib:
			if rollAttribs.hasAttrib(disadvantageAttrib) {
				t.Fatalf("Advantage attrib compatibility check failed, %s is set", rollAttributeMapKey(rollAttributeMap, rollAttrib))
			}
		case disadvantageAttrib:
			if rollAttribs.hasAttrib(advantageAttrib) {
				t.Fatalf("Disadvantage attrib compatibility check failed, %s is set", rollAttributeMapKey(rollAttributeMap, rollAttrib))
			}
		}
	}
}

func FuzzSetRollAttrib(f *testing.F) {
	f.Add(0)
	f.Fuzz(func(t *testing.T, fuzzedRollAttrib int) {
		rollAttribs := newRollAttributes()
		rollAttrib := rollAttribute(fuzzedRollAttrib)
		rollAttribs.setRollAttrib(rollAttrib)

		if !rollAttribs.hasAttrib(rollAttrib) {
			t.Fatalf("Fuzzed attrib %d not set", fuzzedRollAttrib)
		}

		checkForAttribCompatibility(rollAttribs, t)
	})
}
