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

// Basic single DiceRoll
func (diceRoll DiceRoll) performSingleRoll() int {
	fmt.Println("Rolling for", printoutDiceRoll(diceRoll))
	rollSum := 0
	rollStr := ""
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceGen := getFreshRandomGenerator()
		diceRollResult := diceGen.Intn(diceRoll.DiceSize) + 1
		rollStr += fmt.Sprintf("%d ", diceRollResult)
		rollSum += diceRollResult
	}
	fmt.Println("Dice rolls:", rollStr)
	fmt.Println(fmt.Sprintf("Roll sum with mod for %s: %d", printoutDiceRoll(diceRoll), rollSum + diceRoll.Modifier))
	return rollSum + diceRoll.Modifier
}

// Perform an array of DiceRoll
func PerformRolls(diceRolls []DiceRoll) int {
	rollsSum := 0
	for i := 0; i < len(diceRolls); i++ {
		rollsSum += diceRolls[i].performSingleRoll()
	}
	return rollsSum
}

 func getFreshRandomGenerator() *rand.Rand {
 	return rand.New(rand.NewSource(time.Now().UnixNano()))
 }

 func printoutDiceRoll(diceRoll DiceRoll) string {
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += "+"
		} 	
		strDiceRoll += fmt.Sprintf("%d", diceRoll.Modifier)	
	}
	return strDiceRoll
 }