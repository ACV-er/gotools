package units

import (
	"math/rand"
	"testing"
)

func bitCount(n uint64) uint {
	var rel uint = 0
	for n != 0 {
		rel += uint(n & 1)
		n >>= 1
	}

	return rel
}

func getData() []uint64 {
	data := []uint64{}
	for i := 0; i < 100000; i++ {
		data = append(data, rand.Uint64())
	}
	return data
}

func get32Data() []uint64 {
	data := []uint64{}
	for i := 0; i < 100000; i++ {
		data = append(data, uint64(rand.Uint32()))
	}
	return data
}

func TestBitCount32(t *testing.T) {
	data := get32Data()

	for _, n := range data {
		if got := BitCount32(uint32(n)); got != bitCount(n) {
			t.Errorf("BitCount32() = %v, want %v", got, bitCount(n))
		}
	}
}

func TestBitCount64(t *testing.T) {
	data := getData()

	for _, n := range data {
		if got := BitCount64(n); got != bitCount(n) {
			t.Errorf("BitCount32() = %v, want %v", got, bitCount(n))
		}
	}
}
