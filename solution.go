package main

import (
	"fmt"
	"strconv"
)

func (office *Office) expandConnectedComponent(cc *ConnectedComponent, i, j int) {
	// check if the position is valid.
	if i >= 0 && i < office.H && j >= 0 && j < office.W {
		if office.layout[i][j].nodeType != '#' && office.layout[i][j].ccid == 0 {
			// assign ccid.
			office.layout[i][j].ccid = cc.ccid
			// increment count.
			cc.count++
			// expand up.
			office.expandConnectedComponent(cc, i-1, j)
			// expand down.
			office.expandConnectedComponent(cc, i+1, j)
			// expand left.
			office.expandConnectedComponent(cc, i, j-1)
			// expand rigth.
			office.expandConnectedComponent(cc, i, j+1)
		}
	}
}

func (office *Office) getConnectedComponents() []ConnectedComponent {
	var result []ConnectedComponent
	ccid := 1
	for i := 0; i < office.H; i++ {
		for j := 0; j < office.W; j++ {
			if office.layout[i][j].nodeType != '#' && office.layout[i][j].ccid == 0 {
				// found new tile which is not allocated to any connected component.
				// create new connected component.
				cc := NewCC(ccid, Pos{i, j})
				ccid++
				// expand cc.
				office.expandConnectedComponent(&cc, i, j)
				result = append(result, cc)
			}
		}
	}
	return result
}

// next - function used to return the next available position for a manager or a developer.
// @tile0, @tile1: 'm' for manager, 'd' for developer.
// @return: tuple of x and y coordinates.
func (office *Office) next() func(tile0, tile1 byte) *Pair {
	ccSlice := office.getConnectedComponents()
	fmt.Println(office.toString())
	fmt.Println("Found", len(ccSlice), "connected components")
	// define the function.
	result := func(tile0, tile1 byte) *Pair {
		fmt.Println("Finding pair for:", tile0, tile1)
		// iterate over the connected components until you find a position for the (tile0, tile1) pair.
		for ccidx := 0; ccidx < len(ccSlice); ccidx++ {
			cc := ccSlice[ccidx]
			fmt.Println(cc.ccid, cc.di, cc.mi, cc.xi)
			pairSlice := cc.x // initialise with the mixed slice.
			pairIndex := &cc.xi
			if tile0 == tile1 && tile0 == 'd' { // pair of developers.
				fmt.Println("Dev pair.")
				pairSlice = cc.d
				pairIndex = &cc.di
			} else if tile0 == tile1 && tile0 == 'm' { // pair of managers.
				fmt.Println("Man pair")
				pairSlice = cc.m
				pairIndex = &cc.mi
			} else {
				fmt.Println("Mixed pair")
			}
			// check if the current slice has enough pairs.
			if *pairIndex == len(pairSlice) {
				continue
			}
			// check if the pair has both desks available.
			pair := pairSlice[*pairIndex]
			for *pairIndex < len(pairSlice) && office.layout[pair.x0][pair.y0].occupant != nil && office.layout[pair.x1][pair.y1].occupant != nil {
				*pairIndex = *pairIndex + 1
				if *pairIndex < len(pairSlice) {
					pair = pairSlice[*pairIndex]
				}
			}
			// check if the current slice has enough pairs.
			if *pairIndex == len(pairSlice) {
				continue
			}
			*pairIndex = *pairIndex + 1
			fmt.Println("cc:", cc.ccid, "New pairIndex:", *pairIndex, "In cc:", cc.di, cc.mi, cc.xi)
			return &pair
		}
		return nil
	} // end of function definition.
	return result
}

func (office *Office) placeReplyer(pair *Pair, r, s *Replyer) int {
	if r.replType != s.replType && s.replType == office.layout[pair.x0][pair.y0].nodeType {
		office.layout[pair.x0][pair.y0].occupant = s
		office.layout[pair.x1][pair.y1].occupant = r
		return 1
	}
	office.layout[pair.x0][pair.y0].occupant = r
	office.layout[pair.x1][pair.y1].occupant = s
	return 0
}

func findSolution(data *Data) []string {
	positions := make(map[*Replyer]*Pos)
	// obtain 'next' function.
	next := data.office.next()
	// compute total potential and create max heap.
	maxHeap := *data.computeTotalPotential()
	// fmt.Println("The Max Heap is ")
	// for i := 0; i < maxHeap.size; i++ {
	// 	fmt.Print(strconv.Itoa(maxHeap.remove().value) + " ")
	// }
	// fmt.Println()
	// obtain a solution.
	// iterate over each tile.
	for maxHeap.size != 0 {
		// get the best pair of replyers.
		best := maxHeap.remove()
		fmt.Println("Best pair: ", best.r.toString(), ":", best.s.toString())
		// check if either of them is already placed.
		if positions[best.r] == nil && positions[best.s] == nil {
			pair := next(best.r.replType, best.s.replType)
			if pair != nil {
				fmt.Println("Pair: [", pair.x0, ":", pair.y0, "] [", pair.x1, ":", pair.y1, "]")
				// place the 2 replyers to the correspoding place.
				order := data.office.placeReplyer(pair, best.r, best.s)
				if order == 0 {
					positions[best.s] = &Pos{pair.x0, pair.y0}
					positions[best.r] = &Pos{pair.x1, pair.y1}
				} else {
					positions[best.r] = &Pos{pair.x0, pair.y0}
					positions[best.s] = &Pos{pair.x1, pair.y1}
				}
			} else {
				fmt.Println("Pair not found for ", best.r.toString(), ":", best.s.toString())
			}
		} else { // -> if either of them is placed...good luck!
			fmt.Println("One of them is already placed.")
		}
	}
	// create the result, by iterating the developers and managers slice and checking the positions into the positoions map.
	result := []string{}
	// developers.
	for i := 0; i < len(data.devs); i++ {
		pos := positions[&data.devs[i]]
		if pos != nil {
			result = append(result, strconv.Itoa(pos.x)+" "+strconv.Itoa(pos.y))
		} else {
			result = append(result, "X")
		}
	}
	// managers.
	for i := 0; i < len(data.mans); i++ {
		pos := positions[&data.mans[i]]
		if pos != nil {
			result = append(result, strconv.Itoa(pos.x)+" "+strconv.Itoa(pos.y))
		} else {
			result = append(result, "X")
		}
	}
	return result
}
