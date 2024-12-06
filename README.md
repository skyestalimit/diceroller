# diceroller

Easy to use package for any kind of dice rolling. Can perform complex dice rolling expressions in a single function call.

- Provides every individual dice roll result, kept or dropped, for more result viewing fun.
- Production grade library with complete test coverage and fuzz tests. This lib is bug free.
- Supports DnD 5e rolls.

## Installation

```bash
go get github.com/skyestalimit/diceroller
```

## Usage

### Rolling command line style using RollArgs

Pass rollArgs to the following function and get rolling!

```go
PerformRollArgsAndSum("1d6")
PerformRollArgsAndSum("adv", "1d20+5")
```

A valid RollArg matches either the DiceRoll format or a rollAttribute string.

DiceRoll:

- Format: "[X]dY[+|-]Z".
- Examples: "5d6", "d20", "4d4+1", "10d10", "1d6-1", "1D8".

rollAttribute strings:

- roll, hit, dmg : separators, starts a new rolling expressions
- crit: Critical, doubles all dice ammount
- spell: Spell, DiceRollResults.String() prints the sum and the sum halved for saves
- half: Halves the sums, for resistances and such
- adv: Advantage, rolls each dice twice and drops the lowest
- dis: Disadvantage, rolls each dice twice and drops the highest
- drophigh: Drop High, drops the highest result of a DiceRoll
- droplow: Drop Low, drops the lowest result of a DiceRoll

When a rollAttribute is found, a new rolling expression starts. Rolling expressions are a group of DiceRolls and rollAttributes combined to generate the result.

### Rolling library style using DiceRolls

A `DiceRoll` is not necessarily a single dice roll, but a single dice rolling expression, such as `2d6`.

In the end, all the rolls made with this package are made using DiceRolls internally. You can build your own or get them from the RollArg parser.

#### Rolling DiceRolls

You can also build your own DiceRolls and pass them to the following function to roll them:

```go
PerformRollsAndSum(*diceRoll2d6plus1, *diceRoll1d8)
```

### Viewing Results

For more details about the results, `DiceRollResult` or `RollingExpressionResult` slices can be returned instead of a sum by using`PerformRollArgs` or `PerformRolls`. An `error` slice is also returned containing an error for each invalid DiceRolls or RollArgs. Refer to each struct documentation for more details.

You can sum the results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
```

You can sum the results of a `RollingExpressionResult` array by passing it to:

```go
func RollingExpressionResultSum(results ...RollingExpressionResult) (sum int) 
```
