// from: https://golangbyexample.com/maxheap-in-golang/

package main

import "fmt"
import "strconv"

////////////////////////////////////////////////////////////////////////////////

type arrType struct {
	value int      // the total potential of the replyers.
	r, s  *Replyer // the two replyers
}

type heapArrType struct {
	array []arrType
}

func (harr heapArrType) get(i int) arrType {
	return harr.array[i]
}

func (harr heapArrType) set(i int, v arrType) {
	harr.array[i] = v
}

////////////////////////////////////////////////////////////////////////////////

type maxheap struct {
	heapArray heapArrType
	size      int
	maxsize   int
}

func newMaxHeap(maxsize int) *maxheap {
	maxheap := &maxheap{
		heapArray: heapArrType{[]arrType{}},
		size:      0,
		maxsize:   maxsize,
	}
	return maxheap
}

func (m *maxheap) leaf(index int) bool {
	if index >= (m.size/2) && index <= m.size {
		return true
	}
	return false
}

func (m *maxheap) parent(index int) int {
	return (index - 1) / 2
}

func (m *maxheap) leftchild(index int) int {
	return 2*index + 1
}

func (m *maxheap) rightchild(index int) int {
	return 2*index + 2
}

func (m *maxheap) insert(item arrType) error {
	if m.size >= m.maxsize {
		return fmt.Errorf("Heap is full")
	}
	m.heapArray.array = append(m.heapArray.array, item)
	m.size++
	m.upHeapify(m.size - 1)
	return nil
}

func (m *maxheap) swap(first, second int) {
	temp := m.heapArray.get(first)
	m.heapArray.set(first, m.heapArray.get(second))
	m.heapArray.set(second, temp)
}

func (m *maxheap) upHeapify(index int) {
	for m.heapArray.get(index).value > m.heapArray.get(m.parent(index)).value {
		m.swap(index, m.parent(index))
	}
}

func (m *maxheap) downHeapify(current int) {
	if m.leaf(current) {
		return
	}
	largest := current
	leftChildIndex := m.leftchild(current)
	rightRightIndex := m.rightchild(current)
	//If current is smallest then return
	if leftChildIndex < m.size && m.heapArray.get(leftChildIndex).value > m.heapArray.get(largest).value {
		largest = leftChildIndex
	}
	if rightRightIndex < m.size && m.heapArray.get(rightRightIndex).value > m.heapArray.get(largest).value {
		largest = rightRightIndex
	}
	if largest != current {
		m.swap(current, largest)
		m.downHeapify(largest)
	}
	return
}

func (m *maxheap) buildMaxHeap() {
	for index := ((m.size / 2) - 1); index >= 0; index-- {
		m.downHeapify(index)
	}
}

func (m *maxheap) remove() arrType {
	top := m.heapArray.get(0)
	m.heapArray.set(0, m.heapArray.get(m.size-1))
	m.heapArray.array = m.heapArray.array[:(m.size)-1]
	m.size--
	m.downHeapify(0)
	return top
}

func (m maxheap) toString() string {
	result := ""
	for i := 0; i < m.size; i++ {
		result += strconv.Itoa(m.heapArray.get(i).value) + " "
	}
	return result
}
