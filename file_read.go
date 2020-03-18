package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

// readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func readHW(r *bufio.Reader) (int, int) {
	line, err := readln(r)
	check(err)
	i := strings.Index(line, " ")
	W, err := strconv.ParseInt(line[:i], intConvBase10, intBitSize64)
	check(err)
	H, err := strconv.ParseInt(line[i+1:], intConvBase10, intBitSize64)
	check(err)
	return int(W), int(H)
}

func readMap(r *bufio.Reader, W, H int) *[][]Node {
	layout := make([][]Node, H)
	for i := 0; i < H; i++ {
		layout[i] = make([]Node, W)
		line, err := readln(r)
		check(err)
		for j := 0; j < W; j++ {
			layout[i][j] = Node{0, line[j], nil}
		}
	}
	return &layout
}

func readDevs(r *bufio.Reader, data *Data, cid, sid *int) *[]Replyer {
	line, err := readln(r)
	check(err)
	D, err := strconv.ParseInt(line, intConvBase10, intBitSize64) // numbers of developers.
	check(err)
	// read D developers.
	devs := make([]Replyer, D)
	for i := 0; i < int(D); i++ {
		line, err = readln(r)
		check(err)
		// company.
		c := strings.Index(line, " ")
		C := line[:c] // company name.
		ci := data.companies[C]
		if ci == 0 {
			data.companies[C] = *cid
			*cid = *cid + 1
		}
		// bonus.
		b := strings.Index(line[c+1:], " ")
		b = b + c + 1
		B, err := strconv.ParseInt(line[c+1:b], intConvBase10, intBitSize64) // bonus.
		check(err)
		// number of skills.
		s := strings.Index(line[b+1:], " ")
		s = s + b + 1
		S, err := strconv.ParseInt(line[b+1:s], intConvBase10, intBitSize64) // number of skills.
		check(err)
		// skills.
		setOfSkills := make([]int, S)
		for j := 0; j < int(S); j++ {
			t := strings.IndexAny(line[s+1:], " \n")
			if t == -1 {
				t = len(line)
			} else {
				t = t + s + 1
			}
			T := line[s+1 : t]
			si := data.skills[T]
			if si == 0 {
				// new skill.
				data.skills[T] = *sid
				*sid = *sid + 1
			}
			setOfSkills[j] = data.skills[T]
			s = t
		}
		// create new developer.
		devs[i] = Replyer{'d', data.companies[C], int(B), setOfSkills}
	}
	return &devs
}

func readMans(r *bufio.Reader, data *Data, cid, sid *int) *[]Replyer {
	line, err := readln(r)
	check(err)
	M, err := strconv.ParseInt(line, intConvBase10, intBitSize64) // numbers of managers.
	check(err)
	// read managers.
	mans := make([]Replyer, M)
	for i := 0; i < int(M); i++ {
		line, err = readln(r)
		check(err)
		// company.
		c := strings.Index(line, " ")
		C := line[:c] // company name.
		ci := data.companies[C]
		if ci == 0 {
			data.companies[C] = *cid
			*cid = *cid + 1
		}
		// bonus.
		b := strings.Index(line[c+1:], " ")
		if b == -1 {
			b = len(line)
		} else {
			b = b + c + 1
		}
		B, err := strconv.ParseInt(line[c+1:b], intConvBase10, intBitSize64) // bonus.
		check(err)
		// create new manager.
		mans[i] = Replyer{'m', data.companies[C], int(B), nil}
	}
	return &mans
}

func readFile(path string, fname string) *Data {
	// open file.
	file, err := os.Open(path + fname)
	check(err)
	defer file.Close()
	r := bufio.NewReader(file)
	// read M and N.
	W, H := readHW(r)
	// read the map.
	layout := readMap(r, W, H)
	office := Office{int(W), int(H), *layout}
	// read the developers.
	companies := make(map[string]int)
	cid := 1 // company id.
	skills := make(map[string]int)
	sid := 1 // skill id.
	data := Data{office, nil, nil, companies, skills, nil, nil, nil, nil}
	devs := readDevs(r, &data, &cid, &sid)
	// read the managers.
	mans := readMans(r, &data, &cid, &sid)
	data.devs = *devs
	data.mans = *mans
	return &data
}
