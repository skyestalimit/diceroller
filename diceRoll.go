package roller

import (
	"fmt"
	"strconv"
	"strings"
)

// Dice Roll structure
type DiceRoll struct {
	diceSize    int
	diceAmmount int
}

func ParseSingleRollStr(rollStr string) DiceRoll {
	// Parse the dice roll and return the ammount and size of dice
	rollDef := strings.Split(rollStr, "d")
	dice, diceErr := strconv.Atoi(rollDef[0])
	if diceErr != nil {
		fmt.Println(fmt.Sprintf("Error converting dice value of %s: %s", rollDef, diceErr.Error()))
		dice = 0
	}
	ammount, ammountErr := strconv.Atoi(rollDef[1])
	if ammountErr != nil {
		fmt.Println(fmt.Sprintf("Error converting ammount value of %s: %s", rollDef, ammountErr.Error()))
		ammount = 0
	}
	return DiceRoll{dice, ammount}
}
