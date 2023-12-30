package roller

import (
"fmt"
"math/rand"
"time"
)

// Dice Roll structure
type DiceRoll struct {
	DiceAmmount int
	DiceSize    int
	Modifier 	int
}

// func NewDiceRoll(diceAmmount, diceSize, modifier) *DiceRoll{
// 	return &diceRoll := DiceRoll {diceAmmount, diceSize, modifier}
// }

// Basic single dice roll

func PerformSingleRoll(diceRoll DiceRoll) int {

	fmt.Println("Rolling DiceRoll: ", diceRoll)
	rollSum := 0
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceGen := getFreshRandomGenerator()
		diceRollResult := diceGen.Intn(diceRoll.DiceSize) + 1
		fmt.Println("Gen results: ", diceRollResult)
		rollSum += diceRollResult
	}
	fmt.Println("Roll sum: ", rollSum)
	fmt.Println("Roll sum with mod: ", rollSum + diceRoll.Modifier)
	return rollSum + diceRoll.Modifier
}

func PerformRolls(diceRolls []DiceRoll) int {
	rollsSum := 0
	for i := 0; i < len(diceRolls); i++ {
		rollsSum += PerformSingleRoll(diceRolls[i])
	}
	return rollsSum
}

 func getFreshRandomGenerator() *rand.Rand {
 	return rand.New(rand.NewSource(time.Now().UnixNano()))
 }