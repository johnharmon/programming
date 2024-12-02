package main

import (
	"fmt"
	"sort"
)

type ListNode struct {
	Next     *ListNode
	Previous *ListNode
	value    any
}

type ListHead struct {
	head   *ListNode
	tail   *ListNode
	length int
}

// Checks if a given positive or negative index is within the boundaries of the list
func (head *ListHead) IsIB(index int) bool {
	if index >= 0 {
		if index >= head.length {
			return false
		}
		return true
	}
	if -index > head.length {
		return false
	}
	return true
}

// Splices another linked list into the next of the given index of the calling linked list
func (head *ListHead) Splice(spliceHead *ListHead, index int) error {
	if head.IsIB(index) {
		return fmt.Errorf("Index out of bounds")
	}
	spliceStart, _ := head.Index(index)
	spliceEnd, endErr := head.Index(index + 1)
	if endErr != nil {
		return fmt.Errorf("Index out of bounds")
	}
	spliceStart.Next = spliceHead.head
	spliceEnd.Previous = spliceHead.tail
	head.length += spliceHead.length
	return nil
}

func (head *ListHead) Iter() func(func(*ListNode) bool) {
	return func(yield func(*ListNode) bool) {
		for current := head.head; current != nil; current = current.Next {
			if !yield(current) {
				return
			}
		}
	}
}

func (head *ListHead) RevIter() func(func(*ListNode) bool) {
	return func(yield func(*ListNode) bool) {
		for current := head.tail; current != nil; current = current.Previous {
			if !yield(current) {
				return
			}
		}
	}
}

func calcRealIndexWithBounds(length int, index int) (realIndex int, indexErr error) {
	if index > 0 && index >= length {
		indexErr = fmt.Errorf("Index out of bounds")
		return realIndex, indexErr
	} else if index < 0 && -index > length {
		indexErr = fmt.Errorf("Index out of bounds")
		return realIndex, indexErr
	} else {
		realIndex = calcRealIndex(length, index)
		return realIndex, nil
	}
}

func (head *ListHead) Update(newNode *ListNode, index int) (indexErr error) {
	realIndex, indexErr := calcRealIndexWithBounds(head.length, index)
	if indexErr != nil {
		return indexErr
	}
	if realIndex >= 0 {
		currentIndex := 0
		for current := head.head; currentIndex <= realIndex; current = current.Next {
			if currentIndex == realIndex {
				previousNode := current.Previous
				nextNode := current.Next
				previousNode.Next = nextNode
				nextNode.Previous = previousNode
			}
			currentIndex++
		}
	} else {
		currentIndex := -1
		for current := head.tail; currentIndex >= realIndex; current = current.Previous {
			if currentIndex == realIndex {
				previousNode := current.Previous
				nextNode := current.Next
				previousNode.Next = nextNode
				nextNode.Previous = previousNode
			}
			currentIndex--
		}
	}
	return indexErr
}

func (head *ListHead) UpdateVal(index int, value any) (updateErr error) {
	updateNode, updateErr := head.Index(index)
	if updateErr != nil {
		return updateErr
	}
	updateNode.value = value
	return nil
}

func (head *ListHead) Switch(index1 int, index2 int) (switchError error) {
	newNode1, indexErr1 := head.Index(index1)
	if indexErr1 != nil {
		switchError = indexErr1
		return switchError
	}
	newNode2, indexErr2 := head.Index(index2)
	if indexErr2 != nil {
		switchError = indexErr2
		return switchError
	}
	indexErr1 = head.UpdateVal(index1, newNode2.value)
	indexErr2 = head.UpdateVal(index2, newNode1.value)
	if indexErr1 != nil {
		switchError = indexErr1
		return switchError
	}
	if indexErr2 != nil {
		switchError = indexErr2
		return switchError
	}
	return switchError
}

func (head *ListHead) AppendVal(value any) {
	newNode := &ListNode{Next: nil, value: value}
	for current := head.head; current != nil; current = current.Next {
		if current.Next == nil {
			current.Next = newNode
		}
	}
}

func (head *ListHead) Append(newNode *ListNode) {
	for current := head.head; current != nil; current = current.Next {
		if current.Next == nil {
			current.Next = newNode
		}
	}
}

func (head *ListHead) revIndex(index int) (node *ListNode, indexErr error) {
	if index >= 0 {
		return head.fwdIndex(index)
	} else {
		current := head.tail
		for i := -1; i >= index; i-- {
			if current.Previous != nil {
				current = current.Previous
			} else {
				return nil, fmt.Errorf("Invalid index given")
			}
		}
		return current, nil
	}
}

func (head *ListHead) fwdIndex(index int) (node *ListNode, indexErr error) {
	if index < 0 {
		return head.revIndex(index)
	} else {
		current := head.tail
		for i := 0; i >= index; i++ {
			if current.Previous != nil {
				current = current.Next
			} else {
				return nil, fmt.Errorf("Invalid index given")
			}
		}
		return current, nil
	}
}

// Calculates the closest real pos/negative index for a node based on the index passed
// and the length of the list, this may cause it to switch the pos/neg value of the index if it would go past the
// middle node in the basic traversal called by the passed index
func calcRealIndex(length int, index int) (realIndex int) {
	if index >= 0 {
		// This check will mean that it will reverse index into the middle nodes on odd sized linked lists
		if index >= length/2 {
			realIndex = -(length - index)
		} else {
			realIndex = index
		}
	} else {
		absIndex := -index
		if absIndex > length/2 {
			realIndex = length + index
		} else {
			realIndex = index
		}
	}
	return realIndex
}

func (head *ListHead) Index(index int) (node *ListNode, indexErr error) {
	realIndex, indexErr := calcRealIndexWithBounds(head.length, index)
	if indexErr != nil {
		return node, indexErr
	} else if realIndex >= 0 {
		node, indexErr = head.fwdIndex(realIndex)
	} else {
		node, indexErr = head.revIndex(realIndex)
	}
	return node, indexErr
}

//	if int(absIndex) > head.length {
//		return node, fmt.Errorf("Index Error out of bounds")
//	}
//if index >= 0 {
//	if index >= head.length {
//		return node, fmt.Errorf("Index Error out of bounds")
//	} else if index >= head.length/2 {
//		negIndex := -1 * (head.length - index)
//		return head.revIndex(negIndex)
//	}
//	return head.fwdIndex(index)
//} else {
//	absIndex := index * -1
//	if absIndex >= head.length {
//		return node, fmt.Errorf("Index Error out of bounds")
//	} else if absIndex < head.length/2 {
//		posIndex := head.length + index
//		return head.fwdIndex(posIndex)
//	}
//	return head.revIndex(index)
//}
//}

func (head *ListHead) Delete(index int) bool {
	if head == nil {
		return false
	}
	currentIndex := 0
	if index >= head.length {
		return false
	}
	var previousNode *ListNode
	var nextNode *ListNode
	for current := head.head; current != nil; current = current.Next {
		nextNode = current.Next
		if currentIndex == index && currentIndex == 0 {
			if nextNode != nil {
				nextNode.Previous = nil
			}
			head.head = current.Next
			head.length -= 1
			break
		} else if currentIndex == index {
			previousNode.Next = current.Next
			if nextNode != nil {
				nextNode.Previous = previousNode
			}
			head.length -= 1
			break
		} else {
			previousNode = current
			currentIndex++
		}
	}
	return true
}

func fmap[T any, R any](slice []T, f func(T) R) (r []R) {
	for _, item := range slice {
		result := f(item)
		r = append(r, result)
	}
	return r
}

func filter[T any](filterFunc func(T) bool, slice []T) (filtered []T) {
	for _, item := range slice {
		result := filterFunc(item)
		if result {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func BasicSorted[T any](slice []T, valueFunc func(T) int) []T {
	type itemAndValue struct {
		originalIndex int
		item          []T
	}
	itemsAndValues := make([]itemAndValue, len(slice))

	valueAndIdx := map[int]itemAndValue{}
	sortedSlice := make([]T, len(slice))
	for idx, item := range slice {
		value := valueFunc(item)
		valueAndIdx[value] = idx
		values = append(values, value)
	}
	sort.Ints(values)
	for idx, value := range values {
		sortedSlice[idx] = slice[valueAndIdx[value]]
	}
	return sortedSlice
}

func main() {
}
