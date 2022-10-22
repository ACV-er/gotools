package units

func BitCount32(n uint32) uint {
	n = n - ((n >> 1) & 0x55555555)
	n = (n & 0x33333333) + ((n >> 2) & 0x33333333)
	n = (n + (n >> 4)) & 0x0f0f0f0f
	n = n + (n >> 8)
	n = n + (n >> 16)

	return uint(n & 0xff)
}

func BitCount64(n uint64) uint {
	n = n - ((n >> 1) & 0x5555555555555555)
	n = (n & 0x3333333333333333) + ((n >> 2) & 0x3333333333333333)
	n = (n + (n >> 4)) & 0x0f0f0f0f0f0f0f0f
	n = n + (n >> 8)
	n = n + (n >> 16)
	n = n + (n >> 32)

	return uint(n & 0xff)
}
