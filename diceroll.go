package roller

import (
	"fmt"
	"math"
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
	DiceRollStr string
	Dice        []int
	Sum         int
}

const MaxDiceRollValue int = 99999

// DiceRoll constructor
func NewDiceRoll(ammount int, size int, modifier int) (DiceRoll, error) {
	diceRoll := DiceRoll{ammount, size, modifier}
	if diceErr := validateDiceRoll(diceRoll); diceErr == nil {
		return diceRoll, nil
	} else {
		return DiceRoll{}, fmt.Errorf("invalid DiceRoll: %s", diceErr.Error())
	}
}

// DiceRollResult constructor
func NewDiceRollResult(diceRollStr string) DiceRollResult {
	return DiceRollResult{diceRollStr, []int{}, 0}
}

// Perform an array of DiceRoll
func PerformRolls(diceRolls []DiceRoll) (map[*DiceRoll]*DiceRollResult, []error) {
	diceRollMap := map[*DiceRoll]*DiceRollResult{}
	diceErrs := make([]error, 0)
	for i := range diceRolls {
		if result, diceErr := diceRolls[i].PerformRoll(); diceErr == nil {
			diceRollMap[&diceRolls[i]] = result
		} else {
			diceErrs = append(diceErrs, diceErr)
		}
	}
	return diceRollMap, diceErrs
}

// Performs a DiceRoll
// Returns DiceRollResult
func (diceRoll DiceRoll) PerformRoll() (*DiceRollResult, error) {
	// Validate DiceRoll values
	if diceErr := validateDiceRoll(diceRoll); diceErr != nil {
		return nil, diceErr
	}
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
	return diceRollResult, nil
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

func validateDiceRoll(diceRoll DiceRoll) error {
	if diceErr := validateDiceAmmout(diceRoll.DiceAmmount); diceErr != nil {
		return fmt.Errorf("dice roll %s, %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceSize(diceRoll.DiceSize); diceErr != nil {
		return fmt.Errorf("dice roll %s, %s", diceRoll.String(), diceErr.Error())
	}
	if diceErr := validateDiceModifier(diceRoll.Modifier); diceErr != nil {
		return fmt.Errorf("dice roll %s, %s", diceRoll.String(), diceErr.Error())
	}
	return nil
}

func validateDiceAmmout(ammount int) error {
	if ammount > MaxDiceRollValue || ammount <= 0 {
		return fmt.Errorf("invalid dice ammout: %d", ammount)
	}
	return nil
}

func validateDiceSize(size int) error {
	if size > MaxDiceRollValue || size <= 1 {
		return fmt.Errorf("invalid dice size: %d", size)
	}
	return nil
}

func validateDiceModifier(modifier int) error {
	if int(math.Abs(float64(modifier))) > MaxDiceRollValue {
		return fmt.Errorf("invalid dice modifier: %d", modifier)
	}
	return nil
}

// Human readable DiceRollResult String
func (result DiceRollResult) String() string {
	return fmt.Sprintf("Dice rolls: %s Sum: %d", fmt.Sprint(result.Dice), result.Sum)
}

// Seeds a fresh random generator
func getFreshRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
