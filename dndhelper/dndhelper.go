// Dndhelper takes command line arguments, sends them to a rollArg parser,
// and attempts to perform the the resulting rolls.

package main

import (
	"fmt"
	"os"

	"github.com/skyestalimit/diceroller"
)

func main() {
	// Validate rollArgs have been captured
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments: ", len(os.Args)-1)
		printUsage()
		return
	}

	// Captured rollArgs
	rollArgs := os.Args[1:len(os.Args)]

	// Roll!
	results, errs := diceroller.PerformRollArgs(rollArgs...)

	// Print out parsing errors
	for i := range errs {
		fmt.Println(errs[i])
	}

	// Print results
	for i := range results {
		fmt.Println(results[i].String())
	}

	// Print total sum
	fmt.Println("Total sum:", diceroller.RollResultsSum(results...))
}

func printUsage() {
	fmt.Println("Usage:	dndhelper [rollArg...]")
}
