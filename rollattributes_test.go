package diceroller

import (
	"math/rand"
	"testing"
)

func TestRollAttributes(t *testing.T) {
	rollAttribs := newDnDRollAttributes()
	for i := range rollAttributeMap {
		rollAttrib := rollAttributeMapKey(rollAttributeMap, rollAttributeMap[i])
		if rollAttribs.hasAttrib(rollAttrib) == true {
			t.Fatalf("hasAttrib %d returned true, wanted false", i)
		}

		rollAttribs.setRollAttrib(rollAttrib)

		if rollAttribs.hasAttrib(rollAttrib) == false {
			t.Fatalf("hasAttrib %d returned false, wanted true", rollAttrib)
		}

		checkForAttribCompatibility(*rollAttribs, t)
	}
}

func TestPerformRollWithRollAttributes(t *testing.T) {
	rollAttribs := newDnDRollAttributes()
	for i := range rollAttributeMap {
		rollAttrib := rollAttributeMapKey(rollAttributeMap, rollAttributeMap[i])

		rollAttribs.setRollAttrib(rollAttrib)

		diceRoll := validDiceRollsValues[rand.Intn(len(validDiceRollsValues))].diceRoll
		diceRoll.Attribs = rollAttribs

		if result, diceErr := performRoll(diceRoll); diceErr != nil {
			t.Fatalf("DiceRoll %s returned error: %s", diceRoll.String(), diceErr.Error())
		} else if result.Sum == 0 {
			t.Fatalf("DiceRoll %s result = %d, wanted > 0", diceRoll.String(), result.Sum)
		}
	}
}

func checkForAttribCompatibility(rollAttribs dndRollAttributes, t *testing.T) {
	for rollAttrib := range rollAttribs.attribs {
		switch rollAttrib {
		case advantageAttrib:
			if rollAttribs.hasAttrib(disadvantageAttrib) {
				t.Fatalf("Advantage attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
			}
		case disadvantageAttrib:
			if rollAttribs.hasAttrib(advantageAttrib) {
				t.Fatalf("Disadvantage attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
			}
		}
	}
}

func FuzzSetRollAttrib(f *testing.F) {
	f.Add(0)
	f.Fuzz(func(t *testing.T, fuzzedRollAttrib int) {
		rollAttribs := newDnDRollAttributes()
		rollAttrib := rollAttribute(fuzzedRollAttrib)
		rollAttribs.setRollAttrib(rollAttrib)

		if !rollAttribs.hasAttrib(rollAttrib) {
			t.Fatalf("Fuzzed attrib %d not set", fuzzedRollAttrib)
		}

		checkForAttribCompatibility(*rollAttribs, t)
	})
}
