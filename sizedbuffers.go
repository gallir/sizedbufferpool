package sizedbufferpool

import (
	"sync"
)

type SizedBufferPool struct {
	pools     []sync.Pool
	base      uint
	powerBase uint
	chunks    uint
}

// New returns a SizedBufferPool. It allows to split []byte in
// buckets according to its size
func New(minSize uint, chunks uint) (pool *SizedBufferPool) {
	pool = &SizedBufferPool{}
	if minSize < 2 {
		minSize = 2
	}
	if chunks < 1 {
		chunks = 1
	}
	pool.base = minSize
	for minSize > 1 {
		minSize = minSize >> 1
		pool.powerBase++
	}
	if 1<<pool.powerBase < uint(minSize) {
		pool.powerBase++
	}

	pool.chunks = chunks
	pool.pools = make([]sync.Pool, pool.chunks)

	return

}

// Get return a []byte of the specified size
func (p *SizedBufferPool) Get(s int) []byte {
	i := p.index(uint(s))
	v := p.pools[i].Get()
	if v == nil {
		return make([]byte, s, p.cap(i))
	}

	b := v.([]byte)
	if cap(b) >= s {
		return b[:s]
	}

	// The size is smaller, return it to the pool and create another one
	p.Put(b) // Put it back into the right pool
	newCap := p.cap(i)
	if s > newCap {
		newCap = s
	}
	return make([]byte, s, newCap)
}

// Put stores []bytes in its corresponding bucket
func (p *SizedBufferPool) Put(b []byte) {
	p.pools[p.index(uint(cap(b)))].Put(b)
}

func (p *SizedBufferPool) index(n uint) uint {
	n--
	n >>= p.powerBase
	idx := uint(0)
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= p.chunks {
		return p.chunks - 1
	}
	return idx
}

func (p *SizedBufferPool) cap(i uint) int {
	return 1 << (p.powerBase + i)
}
