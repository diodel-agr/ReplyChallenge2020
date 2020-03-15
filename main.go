package main

import (
	"fmt"
	"strconv"
)

func main() {
	// read data from file.
	data := *readFile("input/", "a_solar.txt")
	fmt.Println(data.toString())
	// compute total potential and create max heap.
	maxHeap := *data.computeTotalPotential()
	fmt.Println("The Max Heap is ")
	for i := 0; i < maxHeap.size; i++ {
		fmt.Print(strconv.Itoa(maxHeap.remove().value) + " ")
	}
	fmt.Println()
	// get best solution :)

}
