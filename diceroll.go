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
	Modifier    int
	result      diceRollResult
}

// Dice Roll Result structure
type diceRollResult struct {
	dice []int
	sum  int
}

// Basic single DiceRoll
func (diceRoll DiceRoll) performRoll() int {
	fmt.Println("Rolling for", diceRoll.asString())

	// Generate rolls DiceAmmount times for DiceSize dices
	diceRoll.result.dice = make([]int, 0)
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceGen := getFreshRandomGenerator()
		diceRoll.result.dice = append(diceRoll.result.dice, diceGen.Intn(diceRoll.DiceSize)+1)
		diceRoll.result.sum += diceRoll.result.dice[i]
	}

	// Apply modifier
	diceRoll.result.sum += diceRoll.Modifier
	// Minimum roll result if always 1
	if diceRoll.result.sum <= 0 {
		diceRoll.result.sum = 1
	}

	fmt.Println("Dice rolls:", diceRoll.result.dice, "Sum:", diceRoll.result.sum)
	return diceRoll.result.sum
}

// Human readable DiceRoll String
func (diceRoll DiceRoll) asString() string {
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += "+"
		}
		strDiceRoll += fmt.Sprintf("%d", diceRoll.Modifier)
	}
	return strDiceRoll
}

// Perform an array of DiceRoll
func PerformRolls(diceRolls []DiceRoll) int {
	rollsSum := 0
	for i := range diceRolls {
		rollsSum += diceRolls[i].performRoll()
	}
	return rollsSum
}

// Seeds a fresh random generator
func getFreshRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
