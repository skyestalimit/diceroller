package main

import (
	"math/rand"
	"os"
	"testing"

	"github.com/skyestalimit/diceroller"
)

func TestPerformValidRollArgs(t *testing.T) {
	os.Args = diceroller.ValidRollArgs
	main()
}

func TestPerformInvalidRollArgs(t *testing.T) {
	os.Args = invalidRollArgs
	main()
}

func FuzzPerformRollArgs(f *testing.F) {
	f.Add("8d4-1")
	f.Fuzz(func(t *testing.T, fuzzedRollArg string) {
		os.Args = append(os.Args, fuzzedRollArg)
		main()
	})
}

func FuzzParseManyRollArgs(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, rollArgAmmount int) {
		rollArgsArray := make([]string, 0)
		for i := 0; i < rollArgAmmount; i++ {
			rollArg := ""
			argType := rand.Intn(4) + 1

			switch argType {
			case 1:
				rollArg = validRollArgs[rand.Intn(len(validRollArgs))]
			case 2:
				rollArg = invalidRollArgs[rand.Intn(len(invalidRollArgs))]
			case 3:
				rollArg = validRollArgsAttribs[rand.Intn(len(validRollArgsAttribs))]
			case 4:
				rollArg = invalidRollArgsAttribs[rand.Intn(len(invalidRollArgsAttribs))]
			}

			rollArgsArray = append(rollArgsArray, rollArg)
		}
		os.Args = rollArgsArray
		main()
	})
}
