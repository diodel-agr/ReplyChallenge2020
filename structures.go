package main

import (
	"fmt"
	"strconv"
)

const (
	intBitSize64  int = 64
	intConvBase10 int = 10
)

const (
	replyerDev byte = 'd'
	replyerMan byte = 'm'
)

const (
	nodeWall    byte = '#'
	nodeDeskDev byte = '_'
	nodeDeskMan byte = 'M'
)

// Replyer - structure used to store the information of a replyer (developer or manager)
// @replID: the id of the replyer.
// @replType: m (manager) or d (developer)
// @company: the id of the company.
// @bonus: the replyer bonus.
// @skills: slice storing the skill ids.
type Replyer struct {
	replID   int
	replType byte
	company  int
	bonus    int
	skills   []int
}

func (r Replyer) toString() string {
	result := strconv.Itoa(r.replID) + " " + string(r.replType) + " " +
		strconv.FormatInt(int64(r.company), intConvBase10) + " " +
		strconv.FormatInt(int64(r.bonus), intConvBase10)
	if r.skills != nil {
		result += " ["
		for i := 0; i < len(r.skills); i++ {
			result += " " + strconv.FormatInt(int64(r.skills[i]), intConvBase10)
		}
		result += " ]"
	}
	return result

}

// Pos - structure used to store the position on the map (office floor).
type Pos struct {
	x, y int
}

// Node - structure used to store the type of the floor element: wall, developer desk or manager desk.
// @ccid: the id of the connected component allocated to this node.
// @nodeType: character used to code the type of the element.
// @available: tells if the node has an allocated replyer.
// @postion: the position of this node on the map.
type Node struct {
	ccid      int  // connected component id.
	nodeType  byte // '#' or '_' or 'M'.
	available bool
	position  Pos
	occupant  *Replyer
}

func (n Node) toString() string {
	return string(n.nodeType)
	//return "(" + strconv.Itoa(n.ccid) + ")" + string(n.nodeType)
}

// Office - structure used to store the information regarding the floor.
// @W: floor width.
// @H: floor heigth.
// @vacant: the number of free desks.
// layout: the map of the floor.
type Office struct {
	W, H   int
	layout [][]Node
}

func (o Office) toString() string {
	result := "Office size: " + strconv.FormatInt(int64(o.H), intConvBase10) +
		"x" + strconv.FormatInt(int64(o.W), intConvBase10) + "\n"
	for i := 0; i < o.H; i++ {
		for j := 0; j < o.W; j++ {
			result = result + o.layout[i][j].toString() + " "
		}
		result = result + "\n"
	}
	return result
}

// Data - structure used to store all the input data.
// @office: the office.
// @devs: slice of developers.
// @mans: slice of managers.
// @companies: mapping between company name and company id.
// @skills: mapping between skill name and skill id.
type Data struct {
	office    Office
	devs      []Replyer
	mans      []Replyer
	companies map[string]int // the list of all companies.
	skills    map[string]int // the list of all strings.
	heapDev   *maxheap       // max-heap of developers pair.
	heapMan   *maxheap       // max-heap of managers pair.
	heapMix   *maxheap       // max-heap of manager-developer pair.
	scoreMap  *map[*Replyer]map[*Replyer]int
}

func (d Data) toString() string {
	// office.
	result := "Map:\n" + d.office.toString() + "\nDevelopers:\n"
	// developers.
	size := len(d.devs)
	for i := 0; i < size; i++ {
		result += d.devs[i].toString() + "\n"
	}
	// managers.
	result += "\nManagers:\n"
	size = len(d.mans)
	for i := 0; i < size; i++ {
		result += d.mans[i].toString() + "\n"
	}
	// companies.
	result += "\nCompanies: " + fmt.Sprintf("%v", d.companies) + "\n"
	// skills.
	result += "\nSkills: " + fmt.Sprintf("%v", d.skills) + "\n"
	// heapDev.
	//result += "HeapDev: " + fmt.Sprintf("%v", d.heapDev.toString()) + "\n"
	if d.heapDev != nil {
		result += "HeapDev: " + d.heapDev.toString() + "\n"
	}
	// heapMan.
	//result += "HeapMan: " + fmt.Sprintf("%v", d.heapMan.toString()) + "\n"
	if d.heapMan != nil {
		result += "HeapMan: " + d.heapMan.toString() + "\n"
	}
	// heapMix.
	//result += "HeapMix: " + fmt.Sprintf("%v", d.heapMix.toString()) + "\n"
	if d.heapMix != nil {
		result += "HeapMix: " + d.heapMix.toString() + "\n"
	}
	// scoreMap.
	// m := *(d.scoreMap)
	// for i := 0; i < len(m); i++ {
	// 	result += m[i] + "\n"
	// }
	return result
}

// Pair - structure used to store a pair of desks.
// The pairs can be either developer-developer, manager-manager or mixed.
// @node0, @node1: the nodes of this pair.
// @repl0, repl1: the replyers of this pair.
type Pair struct {
	node0, node1 *Node
}

// ConnectedComponent - structure used to store the information regarding a connected component.
// @id: the d of the component.
// @pairD: pairs of developer-developer desks.
// @pairM: pairs of manager-manager desks.
// @pairX: pairs of developer-manager desks.
// @single: nodes which have not been allocated to a pair.
type ConnectedComponent struct {
	ccid                int
	pairD, pairM, pairX []Pair
	single              []*Node
}

func (cc *ConnectedComponent) toString() string {
	result := strconv.Itoa(cc.ccid) + " pairD: " + strconv.Itoa(len(cc.pairD)) +
		" pairM: " + strconv.Itoa(len(cc.pairM)) + " pairX: " +
		strconv.Itoa(len(cc.pairX)) + " single: " + strconv.Itoa(len(cc.single))
	/*
		for _, d := range cc.pairD {
			result += strconv.Itoa(d.node0.position.x) + ":" + strconv.Itoa(d.node0.position.y) + "/" +
				strconv.Itoa(d.node1.position.x) + ":" + strconv.Itoa(d.node1.position.y) + " "
		}
		result += "\npairM: "
		for _, d := range cc.pairM {
			result += strconv.Itoa(d.node0.position.x) + ":" + strconv.Itoa(d.node0.position.y) + "/" +
				strconv.Itoa(d.node1.position.x) + ":" + strconv.Itoa(d.node1.position.y) + " "
		}
		result += "\npairX: "
		for _, d := range cc.pairX {
			result += strconv.Itoa(d.node0.position.x) + ":" + strconv.Itoa(d.node0.position.y) + "/" +
				strconv.Itoa(d.node1.position.x) + ":" + strconv.Itoa(d.node1.position.y) + " "
		}
	*/
	return result
}

// NewCC - function used to create a new ConnectedComponent variable.
func NewCC(id int) ConnectedComponent {
	return ConnectedComponent{id, []Pair{}, []Pair{}, []Pair{}, []*Node{}}
}
