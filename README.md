# diceroller

Easy to use dice rolling library for any type of dice rolling.

Supports DnD-style rolls.

## Installation

```bash
go get github.com/skyestalimit/diceroller
```

## Usage

For simple dice rolling, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`, see the **Single Dice Rolls** section.

For more complex dice rolling, such as `2d6+4 + 1d8`, see the **Multiple Dice Rolls** section.

### Single rolls

To perform a single dice roll, you simply call PerformRoll() on a DiceRoll and get the sum:

```go
func (diceRoll DiceRoll) PerformRollAndSum() int
```

Or you can get a DiceRollResult instead for a more detailed result:

```go
func (diceRoll DiceRoll) PerformRoll() (*DiceRollResult, error)
```

### Multiple dice rolls

You can perform multiple DiceRolls and get the sum by passing DiceRolls to:

```go
func PerformRollsAndSum(diceRolls ...DiceRoll) int
```

Or you can get a DiceRollResult array instead for more detailed results:

```go
func PerformRolls(diceRolls ...DiceRoll) (results []DiceRollResult, diceErrs []error)
```

The results can then be summed, see the **Viewing Results** section.

### DiceRoll definition

A `DiceRoll` is not necessarily a single dice roll, but a single dice rolling expression, such as `2d6`.

#### Building a DiceRoll

You can build your own `DiceRoll` by using the constructor:

```go
func NewDiceRoll(ammount int, size int, modifier int) (*DiceRoll, error)
```

For example, to build a `DiceRoll` for the dice rolling expression `2d6+1`, you would pass these values to the constructor:

```go
diceRoll, err := NewDiceRoll(2, 6, 1)
```

### RollArg definition

A rollArg is simply a string representing a dice roll expression, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`.

A valid rollArg matches either of these formats: `XdY`, `XdY+Z`, `XdY-Z`.

Examples:

```go
"1d6", "4d4+4", "1D8-1", "10d10", "1d00"
```

#### Building DiceRolss using the rollArg parser

You can use the parser to build a `DiceRoll` from a rollArg:

```go
func ParseRollArg(rollArg string) (*DiceRoll, error)
```

For example, to build a `DiceRoll` for the dice rolling expression `2d6+1` you would pass this rollArg to the above function: `"2d6+1"`

You can also use the parser to build a `DiceRoll` array from multiple rollArgs:

```go
func ParseRollArgs(rollArgs ...string) (diceRolls []DiceRoll, errors []error)
```

For example, to build an array of `DiceRoll` for the dice rolling expression `2d6+4 + 1d8`, you would pass the rollArgs `"2d6+4"` and `"1d8"`.

Example:

```go
ParseRollArgs("2d6+4", "1d8")
```

Or:

```go
ParseRollArgs([]string{"2d6+4","1d8"})
```

### Viewing Results

For more details about the result of a `DiceRoll`, it can return its result in a `DiceRollResult`.

* The result, or sum, is stored in the `Sum` field.
* An array of the individual dice roll results are stored in the `Dice` field.
* `String()` returns a formatted result string.

You can sum results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
```

## Roadmap

The following versions and features are currently planned:

### 1.0.0

* Code coverage support
* Include fuzzing tests

### 1.1.0

* Commands support: `crit`, `spell`
* Introducing Roll attributes to support special rolling features

### 1.2.0

* Command support: `charGen`
