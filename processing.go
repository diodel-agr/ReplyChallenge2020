package main

// workingPotential - function used to compute the working potential of @r and @s replyers.
// @r and @s: thw two replyers.
// @skillCount: the number of skills.
// @return: the working potentian of the two replyers.
func (r Replyer) workingPotential(s Replyer, skillCount int) int {
	// determine the common and different skills.
	skills := make([]int8, 1+skillCount)
	size := len(r.skills)
	if len(s.skills) < len(r.skills) {
		size = len(s.skills)
	}
	i := 0
	for ; i < size; i++ {
		skills[r.skills[i]]++
		skills[s.skills[i]]++
	}
	src := r.skills[:]
	if len(s.skills) > len(r.skills) {
		src = s.skills[:]
		size = len(s.skills)
	} else {
		size = len(r.skills)
	}
	for ; i < size; i++ {
		skills[src[i]]++
	}
	// count the common and different skills.
	common := 0
	distinct := 0
	for i = 0; i < len(skills); i++ {
		if skills[i] == 1 {
			distinct++
		} else if skills[i] == 2 {
			common++
		}
	}
	// count working potential.
	wp := common * (distinct - common)
	return wp
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

// computeTotalPotential - function used to compute the total potential of all pairs of replyers (developers and managers)
func (d Data) computeTotalPotential() *maxheap {
	size := len(d.devs) + len(d.mans)
	maxsize := int((size * size) / 2)
	maxHeap := newMaxHeap(maxsize)
	skillCount := len(d.skills)
	// devs.
	for i := 0; i < len(d.devs); i++ {
		// developers.
		for j := i + 1; j < len(d.devs); j++ {
			value := d.devs[i].totalPotential(d.devs[j], skillCount)
			heapEl := arrType{value, &d.devs[i], &d.devs[j]}
			maxHeap.insert(heapEl)
		}
		// managers.
		for j := 0; j < len(d.mans); j++ {
			value := d.devs[i].totalPotential(d.mans[j], skillCount)
			heapEl := arrType{value, &d.devs[i], &d.mans[j]}
			maxHeap.insert(heapEl)
		}
	}
	// mans.
	n := len(d.mans) - 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < len(d.mans); j++ {
			value := d.mans[i].totalPotential(d.mans[j], skillCount)
			heapEl := arrType{value, &d.mans[i], &d.mans[j]}
			maxHeap.insert(heapEl)
		}
	}
	maxHeap.buildMaxHeap()
	return maxHeap
}
