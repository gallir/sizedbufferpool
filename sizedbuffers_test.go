package sizedbufferpool

import (
	"testing"
)

func TestSizedBufferPool(t *testing.T) {
	base := 12
	chunks := uint(8)
	baseBytes := uint(1 << base)
	max := ((1 << base) << (chunks - 1)) + (1 << base)
	p := New(baseBytes, chunks)
	t.Run("Get_Put_1024", func(t *testing.T) {
		for s := 1; s < max; s++ {
			test(t, p, baseBytes, s)
		}
		for s := max; s > 0; s-- {
			test(t, p, baseBytes, s)
		}
	})
}

func TestSizedBufferPoolReverse(t *testing.T) {
	base := 12
	chunks := uint(8)
	baseBytes := uint(1 << base)
	max := ((1 << base) << (chunks - 1)) + (1 << base)
	p := New(baseBytes, chunks)
	t.Run("Get_Put_1024", func(t *testing.T) {
		for s := max; s > 0; s-- {
			test(t, p, baseBytes, s)
		}
		for s := 1; s < max; s++ {
			test(t, p, baseBytes, s)
		}
	})
}

func test(t *testing.T, p *SizedBufferPool, base uint, s int) {
	b := p.Get(s)
	if len(b.B) != s {
		t.Errorf("Get() = %d, want %d", len(b.B), s)
	}

	if cap(b.B) < s {
		t.Errorf("Get() = %d < %d", cap(b.B), s)
	}

	if cap(b.B) > 2*s+int(base) {
		t.Errorf("Get() = %d > 2 * %d", cap(b.B), s)
	}
	p.Put(b)

}
