// author - 2019: https://github.com/howeih/Day-4-Counting-1-bits  https://github.com/barnybug/popcount

package gbit

func Count1BitsSlow32(value uint32) int {
	n := 0
	for value != 0 {
		value &= value - 1
		n++
	}
	return n
}

func Count1BitsSlow64(value uint64) int {
	n := 0
	for value != 0 {
		value &= value - 1
		n++
	}
	return n
}

func Count1BitsHamming32(i uint32) uint8 {
	i -= (i >> 1) & 0x55555555
	i = (i & 0x33333333) + ((i >> 2) & 0x33333333)
	i = (i + (i >> 4)) & 0x0F0F0F0F
	i = (i * 0x01010101) >> 24
	return uint8(i)
}

func Count1BitsHamming64(i uint64) uint8 {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	i = (i + (i >> 4)) & 0x0F0F0F0F0F0F0F0F
	i = (i * 0x0101010101010101) >> 56
	return uint8(i)
}
