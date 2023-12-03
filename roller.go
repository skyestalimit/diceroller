package roller

// Basic single dice roll

func RollSingleDice(rollStr string) int {
	// parse the rollStr and perform a roll
	return performRoll(ParseSingleRollStr(rollStr))
}

func performRoll(roll DiceRoll) int {

	return roll.diceSize * roll.diceAmmount
}
