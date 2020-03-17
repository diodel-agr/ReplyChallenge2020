package main

import "fmt"

// Pair - structure used to store a pair of desks.
// The pairs can be either developer-developer, manager-manager or mixed.
type Pair struct {
	x0, y0, x1, y1 int
}

// ConnectedComponent - structure used to store the information regarding a connected component.
type ConnectedComponent struct {
	ccid       int // id of he component.
	d, m, x    []Pair
	di, mi, xi int
}

// NewCC - function used to create a new ConnectedComponent variable.
func NewCC(id int) ConnectedComponent {
	return ConnectedComponent{id, []Pair{}, []Pair{}, []Pair{}, 0, 0, 0}
}

func (office *Office) checkRightDownPair(cc *ConnectedComponent, i, j int) {
	this := office.layout[i][j]
	// check right.
	if j+1 < office.W {
		right := office.layout[i][j+1]
		if right.nodeType != '#' {
			// assign the pair (this, right) to the corresponding slice.
			if this.nodeType == right.nodeType && this.nodeType == '_' { // developer desks.
				cc.d = append(cc.d, Pair{i, j, i, j + 1})
			} else if this.nodeType == right.nodeType && this.nodeType == 'M' { // manager desks.
				cc.m = append(cc.m, Pair{i, j, i, j + 1})
			} else { // mixed desks.
				cc.x = append(cc.x, Pair{i, j, i, j + 1})
			}
		}
	}
	// check down.
	if i+1 < office.H {
		down := office.layout[i+1][j]
		if down.nodeType != '#' {
			// assign the pair (this, down) to the corresponding slice.
			if this.nodeType == down.nodeType && this.nodeType == '_' { // developer desks.
				cc.d = append(cc.d, Pair{i, j, i + 1, j})
			} else if this.nodeType == down.nodeType && this.nodeType == 'M' { // manager desks.
				cc.m = append(cc.m, Pair{i, j, i + 1, j})
			} else { // mixed desks.
				cc.x = append(cc.x, Pair{i, j, i + 1, j})
			}
		}
	}
}

func (office *Office) expandConnectedComponent(cc *ConnectedComponent, i, j int) {
	// check if the position is valid.
	if i >= 0 && i < office.H && j >= 0 && j < office.W {
		if office.layout[i][j].nodeType != '#' && office.layout[i][j].ccid == 0 {
			// assign ccid.
			office.layout[i][j].ccid = cc.ccid
			// check for nearby pairs (only right and down).
			office.checkRightDownPair(cc, i, j)
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
				cc := NewCC(ccid)
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
		// iterate over the connected components until you find a position for the (tile0, tile1) pair.
		for ccidx := 0; ccidx < len(ccSlice); ccidx++ {
			cc := ccSlice[ccidx]
			pairSlice := cc.x // initialise with the mixed slice.
			pairIndex := &cc.xi
			if tile0 == tile1 && tile0 == '_' { // pair of developers.
				pairSlice = cc.d
				pairIndex = &cc.di
			} else if tile0 == tile1 && tile0 == 'M' { // pair of managers.
				pairSlice = cc.m
				pairIndex = &cc.mi
			}
			// check if the current slice has enough pairs.
			if *pairIndex == len(pairSlice) {
				continue
			}
			// check if the pair has both desks available.
			pair := pairSlice[*pairIndex]
			for *pairIndex < len(pairSlice) && office.layout[pair.x0][pair.y0].occupant != nil && office.layout[pair.x1][pair.y1].occupant != nil {
				*pairIndex = *pairIndex + 1
				pair = pairSlice[*pairIndex]
			}
			// check if the current slice has enough pairs.
			if *pairIndex == len(pairSlice) {
				continue
			}
			*pairIndex = *pairIndex + 1
			return &pair
		}
		return nil
	} // end of function definition.
	return result
}

func (office *Office) placeReplyer(pair *Pair, r, s *Replyer) int {
	if r.replType != s.replType && s.replType == office.layout[pair.x0][pair.y0].nodeType {
		office.layout[pair.x0][pair.y0].occupant = s
		office.layout[pair.y0][pair.y1].occupant = r
		return 1
	}
	office.layout[pair.x0][pair.y0].occupant = r
	office.layout[pair.y0][pair.y1].occupant = s
	return 0
}

type Pos struct {
	x, y int
}

func findSolution(data *Data) []string {
	positions := make(map[*Replyer]*Pos)
	// obtain 'next' function.
	next := data.office.next()
	// compute total potential and create max heap.
	maxHeap := *data.computeTotalPotential()
	// obtain a solution.
	// iterate over each tile.
	for maxHeap.size != 0 {
		// get the best pair of replyers.
		best := maxHeap.remove()
		// check if either of them is already placed.
		if positions[best.r] != nil && positions[best.s] != nil {
			pair := next(best.r.replType, best.s.replType)
			if pair != nil {
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
		} // -> if either of them is placed...good luck!
	}
	// create the result, by iterating the developers and managers slice and checking the positions into the positoions map.
	result := make([]string, len(data.devs)+len(data.mans))
	// ihi
	return result
}
