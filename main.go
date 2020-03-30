package main

import (
	"fmt"
)

func main() {
	// read data from file.
	data := *readFile("input/", "a_solar.txt")
	// get best solution :)
	solution := findSolution(&data)
	fmt.Println("\nThe solution:\n" + solution)
}

// getAvailableNeighbor
// this function may be improved by returning all the available neighbors and
// choose the one which has the highest score regarding to the one at [i, j]
