package save

import (
	"fmt"
	"math"
)

// This implement the format since Minecraft 1.16
type BitStorage struct {
	data []uint64
	mask uint64

	bits, size    int
	valuesPerLong int
}

// NewBitStorage create a new BitStorage, // TODO: document
func NewBitStorage(bits, size int, arrl []uint64) (b *BitStorage) {
	b = &BitStorage{
		mask:          1<<bits - 1,
		bits:          bits,
		size:          size,
		valuesPerLong: 64 / bits,
	}
	dataLen := (size + b.valuesPerLong - 1) / b.valuesPerLong
	if arrl != nil {
		if len(arrl) != dataLen {
			panic(fmt.Errorf("invalid length given for storage, got: %d but expected: %d", len(arrl), dataLen))
		}
		b.data = arrl
	} else {
		b.data = make([]uint64, dataLen)
	}
	return
}

func (b *BitStorage) cellIndex(n int) int {
	elemPerLong := 64 / b.bits
	return n / elemPerLong
}

func (b *BitStorage) Swap(i, v int) (old int) {
	if i < 0 || i > b.size-1 ||
		v < 0 || uint64(v) > b.mask {
		panic("out of bounds")
	}
	c := b.cellIndex(i)
	l := b.data[c]
	offset := uint64((i - c*b.valuesPerLong) * b.bits)
	old = int(l >> offset & b.mask)
	b.data[c] = l&(b.mask<<offset^math.MaxUint64) | (uint64(v)&b.mask)<<offset
	return
}

func (b *BitStorage) Set(i, v int) {
	if i < 0 || i > b.size-1 ||
		v < 0 || uint64(v) > b.mask {
		panic("out of bounds")
	}
	c := b.cellIndex(i)
	l := b.data[c]
	offset := (i - c*b.valuesPerLong) * b.bits
	b.data[c] = l&(b.mask<<offset^math.MaxUint64) | (uint64(v)&b.mask)<<offset
}

func (b *BitStorage) Get(i int) int {
	if i < 0 || i > b.size-1 {
		panic("out of bounds")
	}
	c := b.cellIndex(i)
	l := b.data[c]
	offset := (i - c*b.valuesPerLong) * b.bits
	return int(l >> offset & b.mask)
}
