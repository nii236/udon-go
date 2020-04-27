package sample

import "fmt"

const Test = 1

type Foo struct {
	PublicBar  int
	privateBaz float64
}

func Sample() int {
	value := privateMethod(1, "hello")
	return value
}
func privateMethod(test int, foo string) int {
	fmt.Println(foo + " world")
	return 1
}
