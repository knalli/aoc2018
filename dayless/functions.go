package dayless

import "strconv"

func Greet(name string) {
	println("Hello " + name)
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
