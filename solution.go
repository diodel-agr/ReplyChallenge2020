package main

import "fmt"

// ConnectedComponent - structure used to store the information regarding a connected component.
type ConnectedComponent struct {
	x, y  int // position of the first tile of the component.
	ccid  int // id of he component.
	count int // number of total positions from a cc.
}

// tileIsCompatible - function used to check wehter the ndoe at position [@i, @j] is compatible with one of the types: @tile0 or @tile1.
// @return: 0 if the node is compatible with @tile0, 1 if the node is compatible with the @tile1 or -1 otherwise.
func tileIsCompatible(office *Office, i, j int, tile0, tile1 byte) int {
	if i < 0 || i >= office.H || j < 0 || j >= office.W {
		return -1
	}
	node := office.layout[i][j]
	if node.nodeType == tile0 {
		return 0
	} else if node.nodeType == tile1 {
		return 1
	}
	return -1
}

// next - this function is used to return the positions of two nodes.
// The types of these nodes correspond to the types of the 2 parameters @tile 0 and @tile1.
func (cc *ConnectedComponent) next(office *Office) func(tile0, tile1 byte) (x0, y0, x1, y1 int) {
	path := stack.New()
	i, j := cc.x, cc.y
	result := func(tile0, tile1 byte) (x0, y0, x1, y1 int) {
		// invalid position.
		if i < 0 || i >= office.H || j < 0 || j >= office.W {
			return -1, -1, -1, -1
		}
		node := office.layout[i][j]
		// not the current cc.
		if node.ccid != cc.ccid {
			return -1, -1, -1, -1
		}
		// search 2 free tiles.
		if tile0 == tile1 {
			// if the current position is free and the desk is compatible with one of the tiles.
			if node.occupant == nil && (tile0 == node.nodeType || tile1 == node.nodeType) {
				// check if this position and one next to this one (up, down, left or right) are compatible.

			}
		} else {

		}

	}
	return result
}

func (office *Office) expandConnectedComponent(cc *ConnectedComponent, i, j int) {
	if i >= 0 && i < office.H && j >= 0 && j < office.W {
		if office.layout[i][j].nodeType != '#' && office.layout[i][j].ccid == 0 {
			// assign ccid.
			office.layout[i][j].ccid = cc.ccid
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
				// found new tile which is not allocated to any connected compnent.
				// create new connected component.
				cc := ConnectedComponent{i, j, ccid, 0}
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
func (office *Office) next() (f func(tile0, tile1 byte) (x0, y0, x1, y1 int)) {
	cc := office.getConnectedComponents()
	fmt.Println(office.toString())
	fmt.Println("Found", len(cc), "connected components")
	ccidx := 0
	i, j := 0, 0 // the current position of a free position.
	// define the function.
	result := func(tile0, tile1 byte) (x0, y0, x1, y1 int) {
		// check to see wether there are available spaces in the current connected component.
		for ccidx < len(cc) && cc[ccidx].count <= 1 {
			ccidx++
			i, j = cc[ccidx].x, cc[ccidx].y
		}
		// chech to see if there are more connected components to explore.
		if ccidx >= len(cc) {
			return -1, -1, -1, -1
		}
		// search the next free position in the current cc.
		for ; i < office.W; i++ {
			for ; j < office.H; j++ {
				if office.layout[i][j].occupant == nil && office.layout[i][j].nodeType != '#' {
					// we have a available space.
				}
			}
		}
		return 0, 0, 0, 0
	}
	return result
}

func findSolution(data *Data) {
	data.office.next()
	return
}
