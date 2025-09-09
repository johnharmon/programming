package main

import (
	"iter"
)

type OrderedSet[T comparable] struct {
	statusBitMap []uint64
	items        []T
	indexMap     map[T]int
	liveElements int
	tombstones   int
}

func (oSet *OrderedSet[T]) Append(elem T) bool {
	if _, ok := oSet.indexMap[elem]; !ok {
		oSet.items = append(oSet.items, elem)
		oSet.indexMap[elem] = len(oSet.items) - 1
		oSet.liveElements += 1
		oSet.appendBitMap(oSet.indexMap[elem])
		return true
	}
	return false
}

func (oSet *OrderedSet[T]) appendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(oSet.statusBitMap) {
		oSet.statusBitMap = append(oSet.statusBitMap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	oSet.statusBitMap[uintIndex] |= bitMask
}

func (oSet *OrderedSet[T]) deleteBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex < len(oSet.statusBitMap) {
		bitMask := uint64(1) << bitOffset
		oSet.statusBitMap[uintIndex] &^= bitMask
	}
}

func (oSet *OrderedSet[T]) Delete(elem T) bool {
	if idx, ok := oSet.indexMap[elem]; ok {
		oSet.deleteBitMap(idx)
		delete(oSet.indexMap, elem)
		return true
	}
	return false
}

func (oSet *OrderedSet[T]) isAlive(elem T) bool {
	if _, ok := oSet.indexMap[elem]; ok {
		return true
	}
	return false
}

func (oSet *OrderedSet[T]) In(elem T) bool {
	if _, ok := oSet.indexMap[elem]; ok {
		return true
	} else {
		return false
	}
}

func (oSet *OrderedSet[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range oSet.items {
			if oSet.isAlive(item) {
				if !(yield(item)) {
					return
				}
			}
		}
	}
}

func (oSet *OrderedSet[T]) IterIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		liveIdx := 0
		for _, item := range oSet.items {
			if oSet.isAlive(item) {
				if !(yield(liveIdx, item)) {
					return
				}
				liveIdx++
			}
		}
	}
}
