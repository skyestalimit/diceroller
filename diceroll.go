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
}

// Dice Roll Result structure
type DiceRollResult struct {
	Dice []int
	Sum  int
}

// Basic single DiceRoll
func (diceRoll DiceRoll) performRoll() *DiceRollResult {
	fmt.Println("Rolling for", diceRoll.asString())

	diceRollResult := new(DiceRollResult)
	// Generate rolls DiceAmmount times for DiceSize dices
	for i := 0; i < diceRoll.DiceAmmount; i++ {
		diceGen := getFreshRandomGenerator()
		diceRollResult.Dice = append(diceRollResult.Dice, diceGen.Intn(diceRoll.DiceSize)+1)
		diceRollResult.Sum += diceRollResult.Dice[i]
	}

	// Apply modifier
	diceRollResult.Sum += diceRoll.Modifier
	// Minimum roll result if always 1
	if diceRollResult.Sum <= 0 {
		diceRollResult.Sum = 1
	}

	fmt.Println("Dice rolls:", diceRollResult.Dice, "Sum:", diceRollResult.Sum)
	return diceRollResult
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
func PerformRolls(diceRolls []DiceRoll) map[*DiceRoll]*DiceRollResult {
	rollMap := map[*DiceRoll]*DiceRollResult{}
	for i := range diceRolls {
		rollMap[&diceRolls[i]] = diceRolls[i].performRoll()
	}
	return rollMap
}

// Seeds a fresh random generator
func getFreshRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
