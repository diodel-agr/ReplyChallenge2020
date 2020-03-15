package main

import (
	"fmt"
	"strconv"
)

const (
	IntBitSize64  int = 64
	IntConvBase10 int = 10
)

// Replyer - structure used to store the information of a replyer (developer or manager)
// @replType: m (manager) or d (developer)
// @company: the id of the company.
// @bonus: the replyer bonus.
// @skills: slice storing the skill ids.
type Replyer struct {
	replType byte
	company  int
	bonus    int
	skills   []int
}

func (r Replyer) toString() string {
	result := string(r.replType) + " " +
		strconv.FormatInt(int64(r.company), IntConvBase10) + " " +
		strconv.FormatInt(int64(r.bonus), IntConvBase10)
	if r.skills != nil {
		result += " ["
		for i := 0; i < len(r.skills); i++ {
			result += " " + strconv.FormatInt(int64(r.skills[i]), IntConvBase10)
		}
		result += " ]"
	}
	return result

}

// Node - structure used to store the type of the floor element: wall, developer desk or manager desk.
// @nodeType: character used to code the type of the element.
// @occupant: pointer to a Replyer object.
type Node struct {
	nodeType byte // '#' or '_' or 'M'.
	occupant *Replyer
}

func (n Node) toString() string {
	return string(n.nodeType)
}

// Office - structure used to store the information regarding the floor.
// @W: floor width.
// @H: floor heigth.
// @vacant: the number of free desks.
// layout: the map of the floor.
type Office struct {
	W, H   int
	vacant int // number of free places.
	layout [][]Node
}

func (o Office) toString() string {
	result := "Office size: " + strconv.FormatInt(int64(o.H), IntConvBase10) +
		"x" + strconv.FormatInt(int64(o.W), IntConvBase10) + "\n"
	result = result + "Vacant places: " + strconv.FormatInt(int64(o.vacant), IntConvBase10) + "\n"
	for i := 0; i < o.H; i++ {
		for j := 0; j < o.W; j++ {
			result = result + o.layout[i][j].toString()
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
}

func (d Data) toString() string {
	// map.
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
	result += "\nSkills: " + fmt.Sprintf("%v", d.skills)
	return result
}
