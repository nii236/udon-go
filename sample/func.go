// from .UdonPie import *  # IGNORE_LINE

// # You can make the variable var_name global
// # by typing "global {var_name}" at the beginning of the code.
// global a
// global b

// # Variable declaration is performed in the following cases.
// # * When assigned for the first time
// # * At the time of variable declaration of argument
// # The scope of all variables is the entire source code.

// # You must always specify the types of function arguments and return values.
// def func1(x_1: Int32, y_1: Int32) -> Int32:
//     return 2 * x_1 + y_1

// def func2(x_2: Int32, y_2: Int32) -> Int32:
//     return x_2 / y_2

// def print_calc() -> Void:
//     a = func1(100, 1000) # 1200
//     b = func2(a, 10) # 120
//     # If the types do not match,
//     # a function call or operation will result in an error.
//     Debug.Log(Object(b)) # output 120

// # Event declarations begin with _.
// def _start():
//     print_calc()

package main

import "udon-go/asm"

func main() {
	a := func1(100, 1000) // 1200
	b := func2(a, 10)     // 120
	asm.Log(string(b))    // output 120
}

func func1(x1 int, y1 int) int {
	return 2*x1 + y1
}
func func2(x2 int, y2 int) int {
	return x2 / y2
}
