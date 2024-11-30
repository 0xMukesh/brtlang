# brtlang

a toy programming language whose code is _brainrot_

## usage

the package is not yet published to https://pkg.go.dev so if you want to try it out, you have to build from source

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
brtlang run test.brt
```

## examples

check out [`examples`](./examples/) folder for examples
