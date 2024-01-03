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

// Perform an array of DiceRoll
func PerformRolls(diceRolls []DiceRoll) map[*DiceRoll]*DiceRollResult {
	diceRollMap := map[*DiceRoll]*DiceRollResult{}
	for i := range diceRolls {
		diceRollMap[&diceRolls[i]] = diceRolls[i].PerformRoll()
	}
	return diceRollMap
}

// Performs a DiceRoll
// Returns DiceRollResult
func (diceRoll DiceRoll) PerformRoll() *DiceRollResult {
	diceRollResult := new(DiceRollResult)

	// Generate rolls DiceAmmount times for DiceSize dices
	fmt.Println("Rolling for " + diceRoll.String())
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

	fmt.Println(diceRollResult.String())
	return diceRollResult
}

// Human readable DiceRoll String
func (diceRoll DiceRoll) String() string {
	strDiceRoll := fmt.Sprintf("%dd%d", diceRoll.DiceAmmount, diceRoll.DiceSize)
	if diceRoll.Modifier != 0 {
		if diceRoll.Modifier > 0 {
			strDiceRoll += "+"
		}
		strDiceRoll += fmt.Sprintf("%d", diceRoll.Modifier)
	}
	return strDiceRoll
}

// Human readable DiceRollResult String
func (result DiceRollResult) String() string {
	return fmt.Sprintf("Dice rolls: %s Sum: %d", fmt.Sprint(result.Dice), result.Sum)
}

// Seeds a fresh random generator
func getFreshRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
