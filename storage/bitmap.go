package storage

import "sync"

/**
基于uint64数组实现的线程安全的位图
*/

type BitMap struct {
	rwLocks         []sync.RWMutex
	bitMapData      []uint64
	lockGranularity uint64
}

func NewBitMap(size uint64, lock_granularity uint64) *BitMap {
	bitMapDataLen := (size-1)/64 + 1
	rwLocksLen := (size-1)/lock_granularity + 1

	return &BitMap{
		rwLocks:         make([]sync.RWMutex, rwLocksLen),
		bitMapData:      make([]uint64, bitMapDataLen),
		lockGranularity: lock_granularity,
	}
}

func (b *BitMap) Get(pos uint64) bool {
	idx := pos / 64
	var mask uint64 = 1
	mask = mask << (pos % 64)

	lock_idx := pos / b.lockGranularity

	b.rwLocks[lock_idx].RLock()
	defer func() {
		b.rwLocks[lock_idx].RUnlock()
	}()

	return (b.bitMapData[idx] & mask) != 0
}

func (b *BitMap) Set(pos uint64, value bool) {
	idx := pos / 64
	var mask uint64 = 1
	mask = mask << (pos % 64)

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
