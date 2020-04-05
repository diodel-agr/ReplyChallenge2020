package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fnamelist := []string{"a_solar.txt", "b_dream.txt", "c_soup.txt", "d_maelstrom.txt", "e_igloos.txt", "f_glitch.txt"}
	pathInput := "input/"
	pathOutput := "output/"
	name := fnamelist[1]
	fmt.Println("Processing file:", name)
	// read data from file.
	start := time.Now()
	data := *readFile(pathInput, name)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("File read ready.", elapsed)
	// get best solution :)
	start = time.Now()
	solution := findSolution(&data)
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println("Solution computed in", elapsed)
	// write solution.
	file, err := os.Create(pathOutput + name)
	if err != nil {
		panic(err)
	}
	file.WriteString(solution)
	file.Close()
	fmt.Println("File saved.")
}

// getAvailableNeighbor
// this function may be improved by returning all the available neighbors and
// choose the one which has the highest score regarding to the one at [i, j]

// placeReplyer - when the (r, s) pair has to be placed, id only one of them
// has allready been placed, find wether there is a free place next to it.
