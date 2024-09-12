package main

import (
	"fmt"
	"math"
	long_name_for_math "math"
)

func something() string {
	return "something"
}
func abs(num float64) float64 {
	return long_name_for_math.Abs(num)
}
func print_any(a any) {
	fmt.Println(a)
}
func print_int_pointee(a *int) {
	fmt.Println(*a)
}
func main() {
	var x = 5
	var y float64 = 7
	y = 4.200000
	const prefix = "John says"
	var text = "hello"
	fmt.Println(fmt.Sprintf("%v: %v world. x: %v, y: %v", prefix, text, x, y))
	const pi = 3.140000
	fmt.Printf("%.2f\n", pi)
	var some_val = something()
	fmt.Println(some_val)
	fmt.Println("hello")
	fmt.Println(abs(-5))
	var binary_expression = math.Pow(10, 2) + 1*0
	fmt.Println(binary_expression)
	if 3 > 1 && true {
		fmt.Println("hi")
	} else {
		fmt.Println("else")
	}
	var _ = "comment: since I haven't implemented a semantic"
	var _ = "         analyzer, the type can't be inferred"
	var _ = "         for if expressions"
	var does_it_work string
	if true {
		does_it_work = "yes"
	} else {
		does_it_work = "no"
	}
	fmt.Println(does_it_work)
	var _ = "comment: the same thing applies for block expressions"
	var another_test string
	{
		var nested string
		{
			nested = "nested"
		}
		fmt.Println("hi")
		another_test = nested
	}
	fmt.Println(another_test)
	var what string
	if true {
		what = "true"
	} else {
		what = "false"
	}
	fmt.Println(what)
	var pointee = 3
	var pointer = &pointee
	print_any(pointee)
	print_any(&pointee)
	print_any(pointer)
	print_int_pointee(pointer)
	print_int_pointee(&pointee)
	var count = 10
	var i = 0
	for {
		i = i + 1
		var _ = "comment: no switch available :("
		var ordinal string
		if i == 1 {
			ordinal = "st"
		} else if i == 2 {
			ordinal = "nd"
		} else if i == 3 {
			ordinal = "rd"
		} else {
			ordinal = "th"
		}
		fmt.Println(fmt.Sprintf("index (%v): hi for the %v%v time", i-1, i, ordinal))
		if i == count {
			break
		}
	}
}
