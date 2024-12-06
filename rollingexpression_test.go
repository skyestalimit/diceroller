package diceroller

import (
	"math/rand"
	"strings"
	"testing"

	"golang.org/x/exp/maps"
)

func TestRollingExpressionWithValidValues(t *testing.T) {
	rollExpr := newRollingExpression()
	rollAttribs := newRollAttributes()

	for i := range validRollArgsAttribs {
		if rollAttrib := checkForRollAttribute(validRollArgsAttribs[i]); rollAttrib > 0 {
			rollAttribs.setRollAttrib(rollAttrib)
		} else {
			t.Fatalf("Valid roll attrib %s has no matching rollAttributes value", validRollArgsAttribs[i])
		}

		for i := range validDiceRollsValues {
			diceRoll := validDiceRollsValues[i].diceRoll
			diceRoll.rollAttribs.setRollAttrib(maps.Keys(rollAttribs.attribs)...)
			rollExpr.diceRolls = append(rollExpr.diceRolls, diceRoll)
		}

		results, diceErrs := performRollingExpressions(*rollExpr)
		resultStr := results[0].String()

		for i := range validDiceRollsValues {
			if !strings.Contains(resultStr, validDiceRollsValues[i].diceRoll.String()) {
				t.Fatalf("Roll result = %s, no match for %s", resultStr, validDiceRollsValues[i].diceRoll.String())
			}
		}

		if diceErrs != nil {
			strErr := ""
			for i := range diceErrs {
				strErr += diceErrs[i].Error() + "\n"
			}
			t.Fatalf("Rolling Expression returned errors: %s", strErr)
		}

		if sum := RollingExpressionResultSum(results...); sum < 1 {
			t.Fatalf("Rolling Expression results sum %d, wanted > 0", sum)

		}
	}
}

func TestRollingExpressionWithInvalidValues(t *testing.T) {
	rollExpr := newRollingExpression()

	for i := range invalidRollArgsAttribs {
		if checkForRollAttribute(invalidRollArgsAttribs[i]) > 0 {
			t.Fatalf("Invalid roll attrib %s has matching rollAttributes value", invalidRollArgsAttribs[i])
		}
	}

	for i := range invalidDiceRollsValues {
		rollExpr.diceRolls = append(rollExpr.diceRolls, invalidDiceRollsValues[i].diceRoll)
	}

	results, diceErrs := performRollingExpressions(*rollExpr)

	if diceErrs == nil {
		t.Fatalf("Invalid Rolling Expression did not return errors")
	}

	if sum := RollingExpressionResultSum(results...); sum > 0 {
		t.Fatalf("Invalid Rolling Expression results %d, wanted <= 0", sum)

	}
}

func FuzzRollingExpression(f *testing.F) {
	f.Add(10, 10)
	f.Fuzz(func(t *testing.T, diceRollAmmount int, attribAmmount int) {
		rollExpr := newRollingExpression()
		rollAttribs := newRollAttributes()

		for i := 0; i < diceRollAmmount; i++ {
			rollExpr.diceRolls = append(rollExpr.diceRolls, *newDiceRoll(rand.Intn(99999)+1, rand.Intn(99999)+1, rand.Intn(99999)+1))
		}

		for i := 0; i < attribAmmount; i++ {
			rollAttribs.setRollAttrib(rollAttribute(i))

			if !rollAttribs.hasAttrib(rollAttribute(i)) {
				t.Fatalf("Random attrib %d is false, wanted true", i)
			}

			checkForAttribCompatibility(*rollAttribs, t)
		}

		_, diceErrs := performRollingExpressions(*rollExpr)

		if diceErrs != nil {
			strErr := ""
			for i := range diceErrs {
				strErr += diceErrs[i].Error() + "\n"
			}
			t.Fatalf("Rolling Expression returned errors: %s", strErr)
		}
	})
}
