package sample

const Test = 1

type Foo struct {
	PublicBar  int
	privateBaz float64
}

func Sample() int {
	world := 0
	if world > 1 {
		world = world + 1
	} else {
		world = world - 1
	}
	value := privateMethod(world, "hello")
	return value
}
func privateMethod(test int, foo string) int {
	test = test + 1
	return test
}
