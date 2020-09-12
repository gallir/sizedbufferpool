package sizedbufferpool

import (
	"sync"
)

// SizedBufferPool Is the strct to kep the pools for []byte sync.Pool
type SizedBufferPool struct {
	pools     []sync.Pool
	base      uint // Smallest bucket size
	powerBase uint // In base 2
	n         uint // NUmber of pools
}

type SizedBuffer struct {
	B []byte
}

// New returns a SizedBufferPool. It allows to split []byte in
// buckets according to its size
// minSize: the smallest buffer size, for example 4096
// buckets: number of pools for diffferen sizes, each os is twice the size of the previous one
// Actual example: sizedbufferpool.New(4096, 8)
func New(minSize uint, buckets uint) (pool *SizedBufferPool) {
	pool = &SizedBufferPool{}
	if minSize < 2 {
		minSize = 2
	}
	if buckets < 1 {
		buckets = 1
	}
	pool.base = minSize
	for minSize > 1 {
		minSize = minSize >> 1
		pool.powerBase++
	}
	if 1<<pool.powerBase < uint(minSize) {
		pool.powerBase++
	}

	pool.n = buckets
	pool.pools = make([]sync.Pool, pool.n)

	return

}

// Get return a []byte of the specified size
func (p *SizedBufferPool) Get(s int) *SizedBuffer {
	i := p.index(uint(s))
	v := p.pools[i].Get()
	if v == nil {
		newCap := p.cap(i)
		if s > newCap {
			newCap = s
		}
		return &SizedBuffer{
			B: make([]byte, s, newCap),
		}
	}

	b := v.(*SizedBuffer)
	if cap(b.B) >= s {
		b.B = b.B[:s]
		return b
	}

	// The size is smaller, return it to the pool and create another one
	p.Put(b) // Put it back into the right pool
	newCap := p.cap(i)
	if s > newCap {
		newCap = s
	}
	return &SizedBuffer{
		B: make([]byte, s, newCap),
	}
}

// Put stores []bytes in its corresponding bucket
func (p *SizedBufferPool) Put(b *SizedBuffer) {
	if cap(b.B) == 0 {
		return
	}
	p.pools[p.index(uint(cap(b.B)))].Put(b)
}

func (p *SizedBufferPool) index(n uint) uint {
	n--
	n >>= p.powerBase
	idx := uint(0)
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= p.n {
		return p.n - 1
	}
	return idx
}

func (p *SizedBufferPool) cap(i uint) int {
	return 1 << (p.powerBase + i)
}
