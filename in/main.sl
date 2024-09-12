package main

import long_name_for_math "math"


fn something() string {
    return "something"
}

fn abs(num float64) float64{
    return long_name_for_math.Abs(num)
}

fn print_any(a any){
    print(a)
}

fn print_int_pointee(a *int){
    print(*a)
}

fn main(){
    let x = 5
    let y: float64 = 7
    y = 4.2
    const prefix = "John says"
    let text = "hello"
    print($"{prefix}: {text} world. x: {x}, y: {y}")
    const pi = 3.14
    printf("%.2f\n", pi)
    let some_val = something()
    print(some_val)
    fmt.Println("hello")
    print(abs(-5))
    let binary_expression = 10**2 + 1 * 0
    print(binary_expression)
    if 3 > 1 && true {
        print("hi")
    } else {
        print("else")
    }

    let _ = "comment: since I haven't implemented a semantic"
    let _ = "         analyzer, the type can't be inferred"
    let _ = "         for if expressions"
    let does_it_work: string = if true {
        "yes"
    } else {
        "no"
    }
    print(does_it_work)

    let _ = "comment: the same thing applies for block expressions"
    let another_test: string = {
        let nested: string = {
            "nested"
        }
        print("hi")
        nested
    }
    print(another_test)

    let what: string = if true {"true"} else {"false"}
    print(what)

    let pointee = 3
    let pointer = &pointee
    print_any(pointee)
    print_any(&pointee)
    print_any(pointer)

    print_int_pointee(pointer)
    print_int_pointee(&pointee)

    let count = 10
    let i = 0
    loop {
        i = i + 1
        let _ = "comment: no switch available :("
        let ordinal: string = if i == 1 {
            "st"
        } else if i == 2 {
            "nd"
        } else if i == 3 {
            "rd"
        } else {
            "th"
        }
        print($"index ({i - 1}): hi for the {i}{ordinal} time")
        if i == count {
            break
        }
    }
}