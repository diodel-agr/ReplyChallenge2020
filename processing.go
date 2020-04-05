package main

// workingPotential - function used to compute the working potential of @r and @s replyers.
// @r and @s: thw two replyers.
// @skillCount: the number of skills.
// @return: the working potentian of the two replyers.
func (r Replyer) workingPotential(s Replyer, skillCount int) int {
	// skills := make(map[int]int)
	// common := 0
	// distinct := 0
	// for i := 0; i < len(r.skills); i++ {
	// 	skills[r.skills[i]] = 1
	// 	distinct++
	// }
	// for i := 0; i < len(s.skills); i++ {
	// 	if skills[s.skills[i]] == 1 {
	// 		common++
	// 		distinct--
	// 	} else {
	// 		distinct++
	// 	}
	// }
	//
	skills := make([]uint8, 1+skillCount)
	common := 0
	distinct := 0
	for i := 0; i < len(r.skills); i++ {
		skills[r.skills[i]]++
		distinct++
	}
	for i := 0; i < len(s.skills); i++ {
		if skills[s.skills[i]] == 0 {
			distinct++
		} else if skills[s.skills[i]] == 1 {
			common++
			distinct--
		}
	}
	//
	// determine the common and different skills.
	// skills := make([]uint8, 1+skillCount)
	// size := len(r.skills)
	// if len(s.skills) < len(r.skills) {
	// 	size = len(s.skills)
	// }
	// i := 0
	// for ; i < size; i++ {
	// 	skills[r.skills[i]]++
	// 	skills[s.skills[i]]++
	// }
	// src := r.skills[:]
	// if len(s.skills) > len(r.skills) {
	// 	src = s.skills[:]
	// 	size = len(s.skills)
	// } else {
	// 	size = len(r.skills)
	// }
	// for ; i < size; i++ {
	// 	skills[src[i]]++
	// }
	// // count the common and different skills.
	// common := 0
	// distinct := 0
	// for i = 0; i < len(skills); i++ {
	// 	if skills[i] == 1 {
	// 		distinct++
	// 	} else if skills[i] == 2 {
	// 		common++
	// 	}
	// }
	// count working potential.
	return common * (distinct - common)
}

// bonusPotential - function used to compute the bonus potential of 2 replyers.
// @r and @s: the two replyers.
// @return: the bonus potential.
func (r Replyer) bonusPotential(s Replyer) int {
	bonus := 0
	if r.company == s.company {
		bonus = r.bonus * s.bonus
	}
	return bonus
}

// totalPotential - function used to compute the total potential of two workers.
// @r and @s the two replyers.
// @return: the total potential.
func (r Replyer) totalPotential(s Replyer, skillCount int) int {
	tp := r.workingPotential(s, skillCount) + r.bonusPotential(s)
	return tp
}

func (d *Data) updateScoreMap(r, s *Replyer, score int) {
	// update r -> s.
	//obtain map.
	sm := (*d.scoreMap)[r]
	// initialise map.
	if sm == nil {
		sm = make(map[*Replyer]int)
		(*d.scoreMap)[r] = sm
	}
	// update map.
	sm[s] = score
	// update s -> r.
	//obtain map.
	sm = (*d.scoreMap)[s]
	// initialise map.
	if sm == nil {
		sm = make(map[*Replyer]int)
		(*d.scoreMap)[s] = sm
	}
	// update map.
	sm[r] = score
}

// computeTotalPotential - function used to compute the total potential of all pairs of replyers (developers and managers)
// and the max-heaps.
func (d *Data) computeTotalPotential() {
	size := len(d.devs) + len(d.mans)
	maxsize := int((size * size) / 2)
	d.heapDev = newMaxHeap(maxsize)
	d.heapMan = newMaxHeap(maxsize)
	d.heapMix = newMaxHeap(maxsize)
	sm := make(map[*Replyer]map[*Replyer]int)
	d.scoreMap = &sm
	skillCount := len(d.skills)
	// devs.
	for i := 0; i < len(d.devs); i++ {
		// developers.
		for j := i + 1; j < len(d.devs); j++ {
			value := d.devs[i].totalPotential(d.devs[j], skillCount)
			heapEl := arrType{value, &d.devs[i], &d.devs[j]}
			d.heapDev.insert(heapEl)
			// update scoreMap.
			d.updateScoreMap(&d.devs[i], &d.devs[j], value)
		}
		// managers.
		for j := 0; j < len(d.mans); j++ {
			value := d.devs[i].totalPotential(d.mans[j], skillCount)
			heapEl := arrType{value, &d.devs[i], &d.mans[j]}
			d.heapMix.insert(heapEl)
			// update scoreMap.
			d.updateScoreMap(&d.devs[i], &d.devs[j], value)
		}
	}
	// mans.
	n := len(d.mans) - 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < len(d.mans); j++ {
			value := d.mans[i].totalPotential(d.mans[j], skillCount)
			heapEl := arrType{value, &d.mans[i], &d.mans[j]}
			d.heapMan.insert(heapEl)
			// update scoreMap.
			d.updateScoreMap(&d.devs[i], &d.devs[j], value)
		}
	}
	d.heapDev.buildMaxHeap()
	d.heapMan.buildMaxHeap()
	d.heapMix.buildMaxHeap()
}
