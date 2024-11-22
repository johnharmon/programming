package main

import (
	"fmt"
	"math"
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

func calcRealIndex(length int, index int) (realIndex int) {
	if index >= 0 {
		if index >= length/2 {
			realIndex = (length - index) * -1
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
	if index >= head.length || -index > head.length {
		indexErr = fmt.Errorf("Index Error out of bounds")
	} else {
		realIndex := calcRealIndex(head.length, index)
		if realIndex >= 0 {
			node, indexErr = head.fwdIndex(realIndex)
		} else {
			node, indexErr = head.revIndex(realIndex)
		}
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
