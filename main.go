package main

import (
	"fmt"
)

func main() {
	// read data from file.
	data := *readFile("input/", "a_solar.txt")
	// get best solution :)
	solution := findSolution(&data)
	for _, s := range solution {
		fmt.Println(s)
	}
}
