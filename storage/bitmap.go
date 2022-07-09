package storage

import "sync"

type BitMap interface {
	Get(uint64) bool
	Set(uint64, bool)
}

func NewBitMap(size uint64, lock_granularity uint64) BitMap {
	bitMapDataLen := (size-1)>>6 + 1
	if lock_granularity == 0 {
		return &bitMap{
			bitMapData: make([]uint64, bitMapDataLen),
		}
	}

	rwLocksLen := (size-1)/lock_granularity + 1
	return &threadSafeBitMap{
		rwLocks:         make([]sync.RWMutex, rwLocksLen),
		bitMapData:      make([]uint64, bitMapDataLen),
		lockGranularity: lock_granularity,
	}
}

/**
基于uint64数组实现的线程安全的位图
*/

type threadSafeBitMap struct {
	rwLocks         []sync.RWMutex
	bitMapData      []uint64
	lockGranularity uint64
}

func (b *threadSafeBitMap) Get(pos uint64) bool {
	idx := pos >> 6
	var mask uint64 = 1
	mask = mask << (pos & 63)

	lock_idx := pos / b.lockGranularity

	b.rwLocks[lock_idx].RLock()
	defer func() {
		b.rwLocks[lock_idx].RUnlock()
	}()

	return (b.bitMapData[idx] & mask) != 0
}

func (b *threadSafeBitMap) Set(pos uint64, value bool) {
	idx := pos >> 6
	var mask uint64 = 1
	mask = mask << (pos & 63)

	lock_idx := pos / b.lockGranularity
	b.rwLocks[lock_idx].Lock()
	defer func() {
		b.rwLocks[lock_idx].Unlock()
	}()

	if value == true {
		b.bitMapData[idx] = b.bitMapData[idx] | mask
	} else {
		b.bitMapData[idx] = b.bitMapData[idx] & ^mask
	}
}

/**
基于uint64的位图
*/

type bitMap struct {
	bitMapData []uint64
}

func (b *bitMap) Get(pos uint64) bool {
	idx := pos >> 6
	var mask uint64 = 1
	mask = mask << (pos & 63)

	return (b.bitMapData[idx] & mask) != 0
}

func (b *bitMap) Set(pos uint64, value bool) {
	idx := pos >> 6
	var mask uint64 = 1
	mask = mask << (pos & 63)

	if value == true {
		b.bitMapData[idx] = b.bitMapData[idx] | mask
	} else {
		b.bitMapData[idx] = b.bitMapData[idx] & ^mask
	}
}
