package main

import (
	"fmt"
)

func main() {
	// read data from file.
	data := *readFile("input/", "a_solar.txt")
	fmt.Println(data.toString())
	// get best solution :)
	solution := findSolution(&data)
	fmt.Println(solution)
}
