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

A RollArg is simply a string representing a dice rolling expression, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`.

- A valid RollArg matches either of these formats: `XdY`, `XdY+Z`, `XdY-Z`.
- It can be a roll attribute such as `crit`, `adv`, `half`, `droplow`.
- It can be a separator, to signal the start of a new dice rolling expression, such as `roll`.
- See the `ParseRollArgs` function documentation for more details.

To perform the dice rolling expression `1d6`:

```go
PerformRollArgsAndSum("1d6")
```

To perform the dice rolling expression `2d6+4 + 1d8`:

```go
PerformRollArgsAndSum("2d6+4", "1d8")
```

You can also call the function `PerformRollArgs` instead to receive more detailed results, see the **Viewing Results** section.

#### Rolling multiple dice rolling expressions using RollArgs

You can build multiple dice rolling expressions, which can be summed separately, using RollArgs, roll attributes or separators.

Performing these RollArgs `adv`, `d20`, `roll`, `1d4`, `half`, `2d6+1`, `roll`, `2d6` at once would result in the rolling expressions `adv d20`, `1d4`, `half 2d6+1` and `2d6` being rolled separately, with their own roll attributes applied and separate sums. Of course, they all can be summed together as well.

### Rolling library style using DiceRolls

A `DiceRoll` is not necessarily a single dice roll, but a single dice rolling expression, such as `2d6`.

In the end, all the rolls made with this package are made using DiceRolls internally. You can build your own or get them from the RollArg parser.

#### Building a DiceRoll

To build a `DiceRoll`, use the constructor to validate your values. For the dice rolling expression `2d6+1`:

```go
NewDiceRoll(2, 6, 1, true)
```

You can set its attributes later or use `NewDiceRollWithAttribs` instead.

#### Building DiceRolls using the RollArg parser

Refer to the `ParseRollArgs` documentation for more details about RollArgs.

You can use the parser to build a `DiceRoll` from one or multiple RollArgs. They will be return grouped in one or many rollingExpression.

To build a `DiceRoll` for the dice rolling expression `2d6+1`:

```go
ParseRollArgs("2d6+1")
```

#### Rolling DiceRolls

For simple dice rolling expressions, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`, you only need a single DiceRoll.

For more complex dice rolling expressions, such as `2d6+4 + 1d8`, you need multiple DiceRolls and sum them.

```go
diceRoll2d6plus1, _ := NewDiceRoll(2, 6, 4, true)
diceRoll1d8, _ := NewDiceRoll(1, 8, 0, true)
PerformRollsAndSum(*diceRoll2d6plus1, *diceRoll1d8)
```

You can also call the function `PerformRolls` instead to receive a more detailed results, see the **Viewing Results** section.

### Viewing Results

For more details about the results, `DiceRollResult` or `RollingExpressionResult` slices can be returned instead of a sum. An `error` slice is also returned containing an error for each invalid DiceRolls or RollArgs. Refer to each struct documentation for more details.

You can sum the results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
```

You can sum the results of a `RollingExpressionResult` array by passing it to:

```go
func RollingExpressionResultSum(results ...RollingExpressionResult) (sum int) 
```
