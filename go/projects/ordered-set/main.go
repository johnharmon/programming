package main

import (
	"iter"
	"sync"
)

type OrderedSet[T comparable] struct {
	bitmap     []uint64
	items      []T
	indexMap   map[T]int
	liveCount  int
	tombstones int
	compacting bool
	rwLock     *sync.RWMutex
}

func (oSet *OrderedSet[T]) Append(elem T) bool {
	if _, ok := oSet.indexMap[elem]; !ok {
		oSet.items = append(oSet.items, elem)
		oSet.indexMap[elem] = len(oSet.items) - 1
		oSet.liveCount += 1
		oSet.appendBitMap(oSet.indexMap[elem])
		return true
	}
	return false
}

func (oSet *OrderedSet[T]) appendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(oSet.bitmap) {
		oSet.bitmap = append(oSet.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	oSet.bitmap[uintIndex] |= bitMask
}

func (oSet *OrderedSet[T]) compact() {
	oSet.rwLock.RLock()
	oSet.compacting = true
	bitmapSnapshot := make([]uint64, len(oSet.bitmap))
	copy(bitmapSnapshot, oSet.bitmap)
	newItems := make([]T, len(oSet.items))
	liveCount := 0
	for idx := range oSet.items {
		if isAlive(bitmapSnapshot, idx) {
			newItems[liveCount] = oSet.items[idx]
			liveCount++
		}
	}
	oSet.rwLock.RUnlock()
	newBitmap := make([]uint64, (liveCount+63)/64)
	bitmapRemainder := liveCount % 64
	for idx := range newBitmap {
		newBitmap[idx] = ^uint64(0)
	}
	if bitmapRemainder != 0 {
		newBitmap[len(newBitmap)-1] = (uint64(1) << bitmapRemainder) - 1
	}
}

func (oSet *OrderedSet[T]) deleteBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex < len(oSet.bitmap) {
		bitMask := uint64(1) << bitOffset
		oSet.bitmap[uintIndex] &^= bitMask
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

func (oSet *OrderedSet[T]) isValAlive(elem T) bool {
	if _, ok := oSet.indexMap[elem]; ok {
		return true
	}
	return false
}

// func getBitMapIndex(idx int) uint64
func (oSet *OrderedSet[T]) isIdxAlive(idx int) bool {
	if idx < 0 || idx > len(oSet.bitmap) {
		return false
	}
	uintIndex := idx / 64
	bitOffset := idx % 64
	bitMask := uint64(1) << bitOffset
	return (oSet.bitmap[uintIndex] & bitMask) != 0
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
			if oSet.isValAlive(item) {
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
			if oSet.isValAlive(item) {
				if !(yield(liveIdx, item)) {
					return
				}
				liveIdx++
			}
		}
	}
}

func isAlive(bitmap []uint64, idx int) bool {
	if idx < 0 || idx > len(bitmap) {
		return false
	}
	uintIndex := idx / 64
	bitOffset := idx % 64
	bitMask := uint64(1) << bitOffset
	return (bitmap[uintIndex] & bitMask) != 0
}
