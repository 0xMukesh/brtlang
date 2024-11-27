### expressions

1. literal - string literals, number literals, booleans and `nil`
2. grouping - `(`expression`)`
3. unary - `!`expression or `-`expression
4. binary - expression operator expression (operator = `+`, `-`, `*`, `\`, `==`, `!=`, `>`, `<`, `>=`, `<=`)

### order of precedence

```
expression → equality ;
equality → comparison ( ( "!=" | "==" ) comparison )* ;
comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term → factor ( ( "-" | "+" ) factor )* ;
factor → unary ( ( "/" | "*" ) unary )* ;
unary → ( "!" | "-" ) unary
| primary ;
primary → NUMBER | STRING | "true" | "false" | "nil"
| "(" expression ")" ;
```

**note**: each rule contains logic for itself and for rules above it

rules in order of precedence (lowest to highest)

1. equality (`==`, `!=`)
2. comparison (`>`, `<`, `>=`, `<=`)
3. term (`+`, `-`)
4. factor (`*`, `/`)
5. unary (`!`, `-`)
6. primary (string literals, number literals, booleans and `nil`)
