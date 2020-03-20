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

// // next - function used to return the next available position for a manager or a developer.
// // @tile0, @tile1: 'm' for manager, 'd' for developer.
// // @return: tuple of x and y coordinates.
// func (office *Office) next() func(tile0, tile1 byte) *Pair {
// 	ccSlice := office.getConnectedComponents()
// 	fmt.Println(office.toString())
// 	fmt.Println("Found", len(ccSlice), "connected components")
// 	// define the function.
// 	result := func(tile0, tile1 byte) *Pair {
// 		fmt.Println("Finding pair for:", tile0, tile1)
// 		// iterate over the connected components until you find a position for the (tile0, tile1) pair.
// 		for ccidx := 0; ccidx < len(ccSlice); ccidx++ {
// 			cc := ccSlice[ccidx]
// 			fmt.Println(cc.ccid, cc.di, cc.mi, cc.xi)
// 			pairSlice := cc.x // initialise with the mixed slice.
// 			pairIndex := &cc.xi
// 			if tile0 == tile1 && tile0 == 'd' { // pair of developers.
// 				fmt.Println("Dev pair.")
// 				pairSlice = cc.d
// 				pairIndex = &cc.di
// 			} else if tile0 == tile1 && tile0 == 'm' { // pair of managers.
// 				fmt.Println("Man pair")
// 				pairSlice = cc.m
// 				pairIndex = &cc.mi
// 			} else {
// 				fmt.Println("Mixed pair")
// 			}
// 			// check if the current slice has enough pairs.
// 			if *pairIndex == len(pairSlice) {
// 				continue
// 			}
// 			// check if the pair has both desks available.
// 			pair := pairSlice[*pairIndex]
// 			for *pairIndex < len(pairSlice) && office.layout[pair.x0][pair.y0].occupant != nil && office.layout[pair.x1][pair.y1].occupant != nil {
// 				*pairIndex = *pairIndex + 1
// 				if *pairIndex < len(pairSlice) {
// 					pair = pairSlice[*pairIndex]
// 				}
// 			}
// 			// check if the current slice has enough pairs.
// 			if *pairIndex == len(pairSlice) {
// 				continue
// 			}
// 			*pairIndex = *pairIndex + 1
// 			fmt.Println("cc:", cc.ccid, "New pairIndex:", *pairIndex, "In cc:", cc.di, cc.mi, cc.xi)
// 			return &pair
// 		}
// 		return nil
// 	} // end of function definition.
// 	return result
// }

func (office *Office) placeReplyer(pair *Pair, r, s *Replyer) int {
	if r.replType != s.replType && s.replType == office.layout[pair.pos0.x][pair.pos0.y].nodeType {
		office.layout[pair.pos0.x][pair.pos0.y].occupant = s
		office.layout[pair.pos1.x][pair.pos1.y].occupant = r
		return 1
	}
	office.layout[pair.pos0.x][pair.pos0.y].occupant = r
	office.layout[pair.pos1.x][pair.pos1.y].occupant = s
	return 0
}

func (office *Office) getNeigh(x, y int) *Pos {
	if x-1 >= 0 { // up.
		neigh := office.layout[x-1][y]
		if neigh.nodeType != nodeWall && neigh.occupant == nil {
			return &Pos{x - 1, y}
		}
	}
	if y-1 >= 0 { // left.
		neigh := office.layout[x][y-1]
		if neigh.nodeType != nodeWall && neigh.occupant == nil {
			return &Pos{x, y - 1}
		}
	}
	if x+1 < office.H { // down.
		neigh := office.layout[x+1][y]
		if neigh.nodeType != nodeWall && neigh.occupant == nil {
			return &Pos{x + 1, y}
		}
	}
	if y+1 < office.W { // right.
		neigh := office.layout[x][y+1]
		if neigh.nodeType != nodeWall && neigh.occupant == nil {
			return &Pos{x, y + 1}
		}
	}
	return nil
}

func (data *Data) getReplyerPair(r, s byte) (*Replyer, *Replyer) {
	// same type (dev-dev or man-man)
	if r == s {
		if r == nodeDeskDev { // dev-dev.
			if data.heapDev.size == 0 {
				return nil, nil
			}
			best := data.heapDev.remove()
			return best.r, best.s
		}
		// man-man.
		if data.heapMan.size == 0 {
			return nil, nil
		}
		best := data.heapMan.remove()
		return best.r, best.s
	}
	// dev-man.
	if data.heapMix.size == 0 {
		return nil, nil
	}
	best := data.heapMix.remove()
	return best.r, best.s
}

func findSolution(data *Data) []string {
	positions := make(map[*Replyer]*Pos)
	// obtain 'next' function.
	//next := data.office.next()
	// compute total potential and create max-heaps.
	data.computeTotalPotential()
	fmt.Println(data.toString())
	// obtain a solution.
	// iterate over each tile.
	available := []Pos{}
	for i := 0; i < data.office.H; i++ {
		for j := 0; j < data.office.W; j++ {
			tile := data.office.layout[i][j]
			if tile.nodeType != nodeWall && tile.occupant == nil {
				// found free tile.
				// find a neighbor desk.
				neigh := data.office.getNeigh(i, j)
				if neigh != nil {
					// found an available desk, find 2 replyers for them.
					neighNodeType := data.office.layout[neigh.x][neigh.y].nodeType
					r, s := data.getReplyerPair(tile.nodeType, neighNodeType)
					fmt.Println(i, j, neigh.x, neigh.y, r.toString(), s.toString())
					if r != nil && s != nil {
						// place these replyers.
						if data.office.placeReplyer(&Pair{Pos{i, j}, Pos{neigh.x, neigh.y}}, r, s) == 0 {
							// order: r, s.
							positions[r] = &Pos{i, j}
							positions[s] = &Pos{neigh.x, neigh.y}
						} else {
							// order: s, r.
							positions[s] = &Pos{i, j}
							positions[r] = &Pos{neigh.x, neigh.y}
						}
					} else {
						// no replyers for this position.
						fmt.Println("No replyers found...")
					}
				} else {
					// nil neighbor, store this position for later processing.
					available = append(available, *neigh)
				}
			}
		}
	}
	// deal with available desks.
	fmt.Println("Available:", available)
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
