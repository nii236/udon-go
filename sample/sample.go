package sample

const Test = 1

func Sample() int {
	value := privateMethod()
	return value
}
func privateMethod() int {
	return 1
}
