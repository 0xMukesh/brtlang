# brtlang

brtlang is a toy, meme-inspired programming language which replaces the conventional programming keywords with meme references.

## usage

the package is not yet published to https://pkg.go.dev so if you want to try it out, you have to build it from source

```
git clone https://github.com/0xmukesh/brtlang
cd brtlang
go build -o brtlang cmd/interpreter/main.go
```

create a new file ending with `.brt` (`test.brt`)

```
yap("hello, world");
```

execute the code via the following command

```
./brtlang run test.brt
```

## language reference

## keywords

| brtlang | golang equivalent |
| ------- | ----------------- |
| rizz    | var               |
| edging  | if                |
| mid     | else if           |
| amogus  | else              |
| bet     | true              |
| cap     | false             |
| nada    | nil               |
| vibin   | while             |
| chillin | for               |
| skibidi | func              |
| bussin  | return            |

## built-in functions

1. `yap(string)` - equivalent to `fmt.Println`

## operators

1. `+` - addition
2. `++` - increment
3. `-` - subtraction
4. `--` - decrement
5. `*` - multiplication
6. `/` - division
7. `%` - modulo
8. `=` - assignment
9. `<` - less than
10. `<=` - less than equal to
11. `>` - greater than
12. `>=` - greater than equal to
13. `==` - equal
14. `&&` - and
15. `||` - or

## examples

check out [`examples`](./examples/) folder for examples
