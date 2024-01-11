# diceroller

Easy to use dice rolling library for any type of dice rolling.

Provides simple rolling functions as well as flexible parsing and rolling options.

Supports DnD-style rolls.

## Installation

```bash
go get github.com/skyestalimit/diceroller
```

## Usage

### Simple rolling using rollArgs

The most intuitive way to roll dice is to send rollArgs to the `PerformRollArgsAndSum` function as shown in examples below.

A rollArg is simply a string representing a dice rolling expression, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`. A valid rollArg matches either of these formats: `XdY`, `XdY+Z`, `XdY-Z`.

To perform the dice rolling expression `1d6`:

```go
PerformRollArgsAndSum("1d6")
```

To perform the dice rolling expression `2d6+4 + 1d8`:

```go
PerformRollArgsAndSum("2d6+4", "1d8")
```

You can also call the function `func PerformRollArgs(rollArgs ...string) ([]DiceRollResult, []error)` instead to receive a more detailed result in a DiceRollResult struct, as well as possible parsing and validation errors.

### Other rolling options

#### Using int values

```go
func PerformRollAndSum(diceAmmount int, diceSize int, modifier int) int
```

```go
func PerformRoll(diceAmmount int, diceSize int, modifier int) (*DiceRollResult, error)
```

To roll the dice rolling expression `2d6`:

```go
PerformRollAndSum(2, 6, 0)
```

### Using DiceRolls

In the end, all the rolls made with this package are made using DiceRolls internally. You can build your own or get them from the rollArg parser.

#### Building a DiceRoll

A `DiceRoll` is not necessarily a single dice roll, but a single dice rolling expression, such as `2d6`.

To build a `DiceRoll` for the dice rolling expression `2d6+1`, use the constructor:

```go
NewDiceRoll(2, 6, 1)
```

#### Building DiceRolls using the rollArg parser

Refer to the *Simple rolling using rollArgs* section for more details about rollArgs.

You can use the parser to build a `DiceRoll` from a rollArg:

To build a `DiceRoll` for the dice rolling expression `2d6+1`:

```go
ParseRollArg("2d6+1")
```

You can also use the parser to build a `DiceRoll` array from multiple rollArgs:

To build an array of `DiceRoll` for the dice rolling expression `2d6+4 + 1d8`:

```go
ParseRollArgs("2d6+4", "1d8")
```

#### Rolling DiceRolls

For simple dice rolling expressions, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`, you only need a single DiceRoll. Call `PerformRollAndSum()` or `PerformRoll()` on a DiceRoll and get the result.

For more complex dice rolling expressions, such as `2d6+4 + 1d8`, you need multiple DiceRolls.

```go
diceRoll2d6plus1, _ := NewDiceRoll(2, 6, 1)
diceRoll1d8, _ := NewDiceRoll(1, 8, 0)
PerformRollsAndSum(*diceRoll2d6plus1, *diceRoll1d8)
```

Or you can get a DiceRollResult array instead for more detailed results:

```go
PerformRolls(*diceRoll2d6plus1, *diceRoll1d8)
```

The results can then be summed, see the **Viewing Results** section.

### Viewing Results

For more details about the result of a `DiceRoll`, a `DiceRollResult` can be returned.

* The result, or sum, is stored in the `Sum` field.
* An array of each individual dice roll result is stored in the `Dice` field.
* `String()` returns a formatted result string.

You can sum results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
```
