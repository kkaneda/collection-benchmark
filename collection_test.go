package collection

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func newKey(s string) []byte {
	return []byte(s)
}

func testAddAndGet(t *testing.T, c Collection) {
	k1 := newKey("aaa")
	if c.Get(k1) != nil {
		t.Error("Unexpected")
	}
	c.Add(k1)
	if !bytes.Equal(c.Get(k1), k1) {
		t.Error("Unexpected")
	}

	k2 := newKey("ccc")
	if c.Get(k2) != nil {
		t.Error("Unexpected")
	}
	c.Add(k2)
	if !bytes.Equal(c.Get(k2), k2) {
		t.Error("Unexpected")
	}

	k3 := newKey("bbb")
	if c.Get(k3) != nil {
		t.Error("Unexpected")
	}
	c.Add(k3)
	if !bytes.Equal(c.Get(k1), k1) {
		t.Error("Unexpected")
	}
	if !bytes.Equal(c.Get(k2), k2) {
		t.Error("Unexpected")
	}
}

func TestSortedSliceAddAndGet(t *testing.T) {
	testAddAndGet(t, &SortedSlice{})
}

func TestLLRBAddAndGet(t *testing.T) {
	testAddAndGet(t, NewLLRB())
}

func TestBtreeAddAndGet(t *testing.T) {
	testAddAndGet(t, NewBTree(2))
}

func TestSortedSliceRandomGet(t *testing.T) {
	s := &SortedSlice{}
	numElems := 1024

	keys := make([][]byte, numElems)
	for i := 0; i < numElems; i++ {
		key := newKey(fmt.Sprintf("a%d", rand.Int63()))
		s.Add(key)
		keys[i] = key
	}
	p := rand.Perm(len(keys))
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(s.Get(key), key) {
			t.Fatal("Not found: %v", key)
		}
	}
}

// benchmarkGet creates b.N elements and randomly look up the elements.
func benchmarkGet(b *testing.B, c Collection) {
	b.StopTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		key := newKey(fmt.Sprintf("a%d", rand.Int63()))
		c.Add(key)
		keys[i] = key
	}
	p := rand.Perm(len(keys))
	b.StartTimer()
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(c.Get(key), key) {
			b.Fatal("Not found: %v", key)
		}
	}
}

//func BenchmarkSortedSlice(b *testing.B) {
//	benchmarkGet(b, &SortedSlice{})
//}

func BenchmarkLLRBGet(b *testing.B) {
	benchmarkGet(b, NewLLRB())
}

func BenchmarkBtree2Get(b *testing.B) {
	benchmarkGet(b, NewBTree(2))
}

func BenchmarkBtree3Get(b *testing.B) {
	benchmarkGet(b, NewBTree(3))
}

func BenchmarkBtree4Get(b *testing.B) {
	benchmarkGet(b, NewBTree(4))
}
