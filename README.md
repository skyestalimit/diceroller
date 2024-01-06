# diceroller

Easy to use dice rolling library for any type of dice rolling. 

Supports DnD-style rolls.

## Installation

```
go get github.com/skyestalimit/diceroller
```

## Usage

For simple dice rolling, such as `1d6`, `1d8-1`, `4d4+4` or `10d10`, see the **Single Dice Rolls** section.

For more complex dice rolling, such as `2d6+4 + 1d8`, see the **Multiple Dice Rolls** section.

### Single Rolls

To perform a single dice roll, you simply call PerformRoll() on a DiceRoll and get the result:

```go
func (diceRoll DiceRoll) PerformRoll() (*DiceRollResult, error)
```

### Multiple Dice Rolls

You can perform multiple DiceRolls by passing a DiceRoll array to:

```go
func PerformRolls(diceRolls []DiceRoll) (results []DiceRollResult, diceErrs []error)
```

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

#### Using Roll Arg Parser to build DiceRolls

You can use the parser to build a single `DiceRoll` from a roll arg string:

```go
func ParseRollArg(rollArg string) (*DiceRoll, error)
```

For example, to build a `DiceRoll` for the dice rolling expression `2d6+1` you would pass this roll arg string to the above function: `"2d6+1"`

You can also use the parser to build a `DiceRoll` array from a roll arg string array:

```go
func ParseRollArgs(rollArgs []string) (diceRolls []DiceRoll, errors []error)
```

For example, to build an array of `DiceRoll` for the dice rolling expression `2d6+4 + 1d8`, you would pass a roll arg string array such as `[]string{"2d6+4","1d8"}` to the above function.

The results can then be summed, see the **Viewing Results** section.

##### Valid Roll Args

Roll args must follow either of these formats: `XdY`, `XdY+Z`, `XdY-Z`.
It will be otherwise declared invalid, will not create a `DiceRoll` and will return an error.

Examples of valid roll args:

```go
"1d6", "4d4+4", "1D8-1", "10d10", "1d00"
```

### Viewing Results

A `DiceRoll` result is returned in a `DiceRollResult`. 
* The result, or sum, is stored in the `Sum` field.
* The individual dice rolls are stored in the `Dice` field. 
* `String()` returns a formatted result string.

You can sum results of a `DiceRollResult` array by passing it to:

```go
func DiceRollResultsSum(results []DiceRollResult) (sum int)
``` 

## Roadmap

The following versions and features are currently planned:

* 1.0.0
** Code coverage support
** Include fuzzing tests

* 1.1.0
** Commands support: `crit`, `spell`
** Introducing Roll attributes to support special rolling features

* 1.2.0
** Command support: `charGen`