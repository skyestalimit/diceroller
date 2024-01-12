# diceroller

Package diceroller generates either a sum or DiceRollResults out of RollArgs or DiceRolls. It's just that simple so get rolling!

Supports DnD-style rolls.

## Installation

```bash
go get github.com/skyestalimit/diceroller
```

## Usage

### Rolling using RollArgs

The most intuitive way to roll dice is to send RollArgs to the `PerformRollArgsAndSum` function as shown in examples below.

A RollArg is simply a string representing a dice rolling expression, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`. A valid RollArg matches either of these formats: `XdY`, `XdY+Z`, `XdY-Z`.

To perform the dice rolling expression `1d6`:

```go
PerformRollArgsAndSum("1d6")
```

To perform the dice rolling expression `2d6+4 + 1d8`:

```go
PerformRollArgsAndSum("2d6+4", "1d8")
```

You can also call the function `PerformRollArgs` instead to receive a more detailed results, see the *Viewing Results* section.

### Rolling using DiceRolls

A `DiceRoll` is not necessarily a single dice roll, but a single dice rolling expression, such as `2d6`.

In the end, all the rolls made with this package are made using DiceRolls internally. You can build your own or get them from the RollArg parser.

#### Building a DiceRoll

To build a `DiceRoll` for the dice rolling expression `2d6+1`, use the constructor:

```go
NewDiceRoll(2, 6, 1)
```

#### Building DiceRolls using the RollArg parser

Refer to the *Simple rolling using RollArgs* section for more details about RollArgs.

You can use the parser to build a `DiceRoll` from one or multiple RollArgs:

To build a `DiceRoll` for the dice rolling expression `2d6+1`:

```go
ParseRollArgs("2d6+1")
```

#### Rolling DiceRolls

For simple dice rolling expressions, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`, you only need a single DiceRoll.

For more complex dice rolling expressions, such as `2d6+4 + 1d8`, you need multiple DiceRolls and sum them.

```go
diceRoll2d6plus1, _ := NewDiceRoll(2, 6, 4)
diceRoll1d8, _ := NewDiceRoll(1, 8, 0)
PerformRollsAndSum(*diceRoll2d6plus1, *diceRoll1d8)
```

You can also call the function `PerformRolls` instead to receive a more detailed results, see the *Viewing Results* section.

### Viewing Results

For more details about the result of a `DiceRoll`, a `DiceRollResult` array can be returned by calling either `PerformRolls` or `PerformRollArgs`. An error array is also returned for invalid DiceRolls or RollArgs. You will get one `DiceRollResult` per valid DiceRoll or `RollArg`, and one error for invalid ones.

* The result, or sum, is stored in the `Sum` field.
* An array of each individual dice roll result is stored in the `Dice` field.
* `String()` returns a formatted result string.

You can sum the results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
```
