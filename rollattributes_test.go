package diceroller

import (
	"testing"
)

func TestRollAttributes(t *testing.T) {
	rollAttribs := newRollAttributes()
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

func checkForAttribCompatibility(rollAttribs rollAttributes, t *testing.T) {
	for rollAttrib := range rollAttribs.attribs {
		switch rollAttrib {
		case advantageAttrib:
			if rollAttribs.hasAttrib(disadvantageAttrib) || rollAttribs.hasAttrib(dropHighAttrib) || rollAttribs.hasAttrib(dropLowAttrib) {
				t.Fatalf("Advantage attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
			}
		case disadvantageAttrib:
			if rollAttribs.hasAttrib(advantageAttrib) || rollAttribs.hasAttrib(dropHighAttrib) || rollAttribs.hasAttrib(dropLowAttrib) {
				t.Fatalf("Disadvantage attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
			}
		case dropHighAttrib:
			if rollAttribs.hasAttrib(advantageAttrib) || rollAttribs.hasAttrib(disadvantageAttrib) || rollAttribs.hasAttrib(dropLowAttrib) {
				t.Fatalf("Drop high attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
			}
		case dropLowAttrib:
			if rollAttribs.hasAttrib(advantageAttrib) || rollAttribs.hasAttrib(disadvantageAttrib) || rollAttribs.hasAttrib(dropHighAttrib) {
				t.Fatalf("Drop low attrib compatibility check failed, %s is set", rollAttributeMap[rollAttrib])
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

		checkForAttribCompatibility(*rollAttribs, t)
	})
}
