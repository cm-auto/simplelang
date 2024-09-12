# What and Why?

A little toy language I implemented to learn how parsing a language works.\
The compile target is go, which makes it simple to implement, automatically as
portable as go and allows the usage of go modules.\
It consists of a [lexer](src/token/token.go), an [ast parser](src/ast/ast.go)
and a [builder](src/builder/builder.go), which generates go code from the ast.\
I haven't implemented a semantic analyzer, because the go compiler basically
does a lot of that work for free when compiling the output code.\
Originally the language was supposed to only be capable of value declarations
and expressions (including call expressions) and then a call expression like
`print("hello world")` would be compiled to something like the following:

```go
package main

import "fmt"

func main(){
    fmt.Println("hello world")
}
```

However I kept getting more and more interesting ideas and so the language grew,
i.e. if expressions:

```
let does_it_work: string = if true {
    "yes"
} else {
    "no"
}
```

Which compile to the following go code:

```go
var does_it_work string
if true {
    does_it_work = "yes"
} else {
    does_it_work = "no"
}
```

Or block expressions:

```
let another_test: string = {
    let nested: string = {
        "nested"
    }
    print("hi")
    nested
}
```

Go code:

```go
var another_test string
{
    var nested string
    {
        nested = "nested"
    }
    fmt.Println("hi")
    another_test = nested
}
```

Since I didn't really design the language and I kept adding features, the code
is quite messy but I achieved what I hoped with this project, which is getting
an understanding of how compilers work.

# Playing with the language

If you want to play around with the language, just edit [in/main.sl](in/main.sl)
and run `go run src/main.go`, this will output the file
[out/main.go](out/main.go), which should be run with `go run out/main.go`.
