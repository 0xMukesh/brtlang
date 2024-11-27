### expressions

1. literal - string literals, number literals, booleans and `nil`
2. grouping - `(`expression`)`
3. unary - `!`expression or `-`expression
4. binary - expression operator expression (operators includes `+`, `-`, `*`, `\`, `==`, `!=`, `>`, `<`, `>=`, `<=`)

### order of precedence

```
expression → equality ;
equality → comparison ( ( "!=" | "==" ) comparison )* ;
comparison → additive ( ( ">" | ">=" | "<" | "<=" ) additive )* ;
additive → multiplicative ( ( "-" | "+" ) multiplicative )* ;
multiplicative → unary ( ( "/" | "*" ) unary )* ;
unary → ( "!" | "-" ) unary
| primary ;
primary → NUMBER | STRING | "true" | "false" | "nil"
| "(" expression ")" ;
```

**note**: each rule contains logic for itself and for rules above it

rules in order of precedence (lowest to highest)

1. equality (`==`, `!=`)
2. comparison (`>`, `<`, `>=`, `<=`)
3. additive (`+`, `-`)
4. multiplicative (`*`, `/`)
5. unary (`!`, `-`)
6. primary (string literals, number literals, booleans and `nil`)
