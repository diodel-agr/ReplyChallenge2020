package main

import (
	"fmt"
	"strconv"
	"time"
)

func (cc *ConnectedComponent) addPair(n0, n1 *Node) {
	pair := Pair{n0, n1}
	if n0.nodeType == n1.nodeType {
		// both nodes have the same type: manager desk or developer desk.
		if n0.nodeType == nodeDeskDev {
			// developer desk.
			cc.pairD = append(cc.pairD, pair)
		} else {
			// manager desk.
			cc.pairM = append(cc.pairM, pair)
		}
	} else {
		// nodes have different types -> mixed pair.
		cc.pairX = append(cc.pairX, pair)
	}
}

func (cc *ConnectedComponent) addSingle(n0 *Node) {
	cc.single = append(cc.single, n0)
}

func (office *Office) getAvailableNeighbor(i, j int) *Node {
	var neigh *Node = nil
	if i > 0 && office.layout[i-1][j].nodeType != nodeWall && office.layout[i-1][j].available == true { // check up.
		neigh = &office.layout[i-1][j]
	} else if i+1 < office.H && office.layout[i+1][j].nodeType != nodeWall && office.layout[i+1][j].available == true { // check down.
		neigh = &office.layout[i+1][j]
	} else if j > 0 && office.layout[i][j-1].nodeType != nodeWall && office.layout[i][j-1].available == true { // check left.
		neigh = &office.layout[i][j-1]
	} else if j+1 < office.W && office.layout[i][j+1].nodeType != nodeWall && office.layout[i][j+1].available == true { // check right.
		neigh = &office.layout[i][j+1]
	}
	return neigh
}

func (office *Office) getNeighbors(r *Node) []*Node {
	result := []*Node{}
	pos := r.position
	// up.
	if pos.x > 0 {
		node := office.layout[pos.x-1][pos.y]
		if node.nodeType != nodeWall && node.occupant != nil {
			result = append(result, &node)
		}
	}
	// down.
	if pos.x < office.H-1 {
		node := office.layout[pos.x+1][pos.y]
		if node.nodeType != nodeWall && node.occupant != nil {
			result = append(result, &node)
		}
	}
	// left.
	if pos.y > 0 {
		node := office.layout[pos.x][pos.y-1]
		if node.nodeType != nodeWall && node.occupant != nil {
			result = append(result, &node)
		}
	}
	// right.
	if pos.y < office.W-1 {
		node := office.layout[pos.x][pos.y+1]
		if node.nodeType != nodeWall && node.occupant != nil {
			result = append(result, &node)
		}
	}
	return result
}

func (office *Office) expandConnectedComponent(cc *ConnectedComponent, i, j int) {
	// check if the position is valid.
	if i >= 0 && i < office.H && j >= 0 && j < office.W {
		node := &office.layout[i][j]
		if node.nodeType != '#' && node.ccid == 0 {
			// assign ccid.
			node.ccid = cc.ccid
			// check if the current tile has been allocated to a pair.
			if node.available == true {
				// allocate to a pair OR add it to the single slice.
				// find a neighbor.
				neigh := office.getAvailableNeighbor(i, j)
				if neigh != nil {
					// add it to one Pair slice of this cc.
					neigh.available = false
					cc.addPair(node, neigh)
				} else {
					// add @node to the single list.
					cc.addSingle(node)
				}
			}
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

func (data *Data) updateScoreNeigh(r, s *Node, score *int) {
	scoremap := *(data.scoreMap)
	// find neigh of r.
	neigh := data.office.getNeighbors(r)
	for _, n := range neigh {
		if n.occupant != s.occupant {
			*score = *score + scoremap[r.occupant][n.occupant]
		}
	}
	// find neigh of s.
	neigh = data.office.getNeighbors(s)
	for _, n := range neigh {
		if n.occupant != r.occupant {
			*score = *score + scoremap[s.occupant][n.occupant]
		}
	}
}

// placeReplyer - function used to iterate over the slices of a pair of a connected component and update the score and the placed map.
func (data *Data) placeReplyer(cc *ConnectedComponent, score *int, placed *map[*Replyer]*Node, pair *[]Pair, heap *maxheap) {
	for _, p := range *pair {
		// obtain the best pair of developers from the heap and place it into the pair.
		ok := 0
		for heap.size > 0 {
			devs := heap.remove()
			// check wether either of the developers have already been placed.
			if (*placed)[devs.r] != nil || (*placed)[devs.s] != nil { // TODO: this condition could be improoved.
				continue
			}
			// place developers.
			if (devs.r.replType == replyerDev && p.node0.nodeType == nodeDeskDev) ||
				(devs.r.replType == replyerMan && p.node0.nodeType == nodeDeskMan) {
				p.node0.occupant = devs.r
				p.node1.occupant = devs.s
				// mark them as placed.
				(*placed)[devs.r] = p.node0
				(*placed)[devs.s] = p.node1
				// fmt.Println("Placed dev", devs.r.replID, "on", p.node0.position.x, p.node0.position.y)
				// fmt.Println("Placed dev", devs.s.replID, "on", p.node1.position.x, p.node1.position.y, "+", devs.value)
			} else {
				p.node0.occupant = devs.s
				p.node1.occupant = devs.r
				// mark them as placed.
				(*placed)[devs.s] = p.node0
				(*placed)[devs.r] = p.node1
				// fmt.Println("Placed dev", devs.s.replID, "on", p.node0.position.x, p.node0.position.y)
				// fmt.Println("Placed dev", devs.r.replID, "on", p.node1.position.x, p.node1.position.y, "+", devs.value)
			}
			// update score.
			*score = *score + devs.value
			data.updateScoreNeigh(p.node0, p.node1, score)
			ok = 1
			break
		}
		if ok != 1 {
			// place this pair into the 'single' slice.
			cc.addSingle(p.node0)
			cc.addSingle(p.node1)
			break
		}
	}
}

// composeResult - function used to create the result string of the solution.
// @data: Data object.
// @placed: map of placed replyers.
// @return: string representing the arrangement of replyers.
func (data *Data) composeResult(placed *map[*Replyer]*Node) string {
	result := ""
	for i := range data.devs {
		if (*placed)[&data.devs[i]] != nil {
			pos := (*placed)[&data.devs[i]].position
			result += strconv.Itoa(pos.y) + " " + strconv.Itoa(pos.x) + "\n"
		} else {
			result += "X\n"
		}
	}
	for i := range data.mans {
		if (*placed)[&data.mans[i]] != nil {
			pos := (*placed)[&data.mans[i]].position
			result += strconv.Itoa(pos.y) + " " + strconv.Itoa(pos.x) + "\n"
		} else {
			result += "X\n"
		}
	}
	return result
}

// findSolution - function used to find the solution of the problem.
// @data: the input of the problem.
// @result: the result.
func findSolution(data *Data) string {
	// compute total potential and create max-heaps.
	start := time.Now()
	data.computeTotalPotential()
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Total potentials computed.", elapsed)
	// obtain connected components.
	start = time.Now()
	conncomp := data.office.getConnectedComponents()
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println("Connected components computed.", elapsed)
	// obtain a solution.
	placed := make(map[*Replyer]*Node) // map of placed replyers.
	score := 0                         // the score of the current arrangement.

	// iterate over each connected component to place the best pairs of replyers.
	for _, cc := range conncomp {
		fmt.Println("Working with cc: ", cc.toString())
		// for each pair of replyers in this cc's slices, find the best replyers.
		// developers pair.
		data.placeReplyer(&cc, &score, &placed, &cc.pairD, data.heapDev)
		// managers pair.
		data.placeReplyer(&cc, &score, &placed, &cc.pairM, data.heapMan)
		// mixed pair.
		data.placeReplyer(&cc, &score, &placed, &cc.pairX, data.heapMix)
	}
	// TODO: iterate over each connected component to fill the 'single' nodes.
	// for _, cc := range conncomp {
	// 	for _, n := range cc.single {
	// 		// based on the neighbors, find the best match.
	//
	// 	}
	// }
	// print score and scoreMap.
	// for _, r := range *(data.scoreMap) {
	// 	for _, s := range r {
	// 		fmt.Print(s, " ")
	// 	}
	// 	fmt.Println()
	// }
	fmt.Println("Score:", score)
	// compose the result and return.
	return data.composeResult(&placed)
}
