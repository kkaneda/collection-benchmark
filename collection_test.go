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

func testBasic(t *testing.T, c Collection) {
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

	if !bytes.Equal(c.Delete(k2), k2) {
		t.Error("Unexpected")
	}
	if c.Get(k2) != nil {
		t.Error("Unexpected")
	}

	k4 := newKey("000")
	if c.Get(k4) != nil {
		t.Error("Unexpected")
	}
	c.Add(k4)
	if !bytes.Equal(c.Get(k4), k4) {
		t.Error("Unexpected")
	}
}

func TestBasic(t *testing.T) {
	testBasic(t, &SortedSlice{})
}

func TestLLRBAddAndGet(t *testing.T) {
	testBasic(t, NewLLRB())
}

func TestBtreeAddAndGet(t *testing.T) {
	testBasic(t, NewBTree(2))
}

func TestSortedSliceRandomGet(t *testing.T) {
	s := &SortedSlice{}
	numElems := 1024

	keys := make([][]byte, numElems)
	for i := 0; i < numElems; i++ {
		key := newKey(fmt.Sprintf("a%020d", rand.Int63()))
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

	p = rand.Perm(len(keys))
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(s.Delete(key), key) {
			t.Fatal("Not found: %v", key)
		}
	}
}

func TestLazySortedSliceRandomGet(t *testing.T) {
	s := &LazySortedSlice{}
	numElems := 1024

	keys := make([][]byte, numElems)
	for i := 0; i < numElems; i++ {
		key := newKey(fmt.Sprintf("a%020d", rand.Int63()))
		s.Add(key)
		keys[i] = key
	}
	s.Freeze()

	p := rand.Perm(len(keys))
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(s.Get(key), key) {
			t.Fatal("Not found: %v", key)
		}
	}

	p = rand.Perm(len(keys))
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(s.Delete(key), key) {
			t.Fatal("Not found: %v", key)
		}
	}
}

// benchmarkGet creates b.N elements and randomly looks up the elements.
func benchmarkGet(b *testing.B, c Collection) {
	b.StopTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		key := newKey(fmt.Sprintf("a%020d", rand.Int63()))
		c.Add(key)
		keys[i] = key
	}
	p := rand.Perm(len(keys))
	c.Freeze()
	b.StartTimer()
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(c.Get(key), key) {
			b.Fatal("Not found: %v", key)
		}
	}
}

func BenchmarkGetSortedSlice(b *testing.B) {
	benchmarkGet(b, &SortedSlice{})
}

//func BenchmarkGetLazySortedSlice(b *testing.B) {
//	benchmarkGet(b, &LazySortedSlice{})
//}
func BenchmarkGetLLRB(b *testing.B) {
	benchmarkGet(b, NewLLRB())
}
func BenchmarkGetBtree2(b *testing.B) {
	benchmarkGet(b, NewBTree(2))
}
func BenchmarkGetBtree4(b *testing.B) {
	benchmarkGet(b, NewBTree(4))
}
func BenchmarkGetBtree8(b *testing.B) {
	benchmarkGet(b, NewBTree(8))
}
func BenchmarkGetBtree16(b *testing.B) {
	benchmarkGet(b, NewBTree(16))
}
func BenchmarkGetBtree32(b *testing.B) {
	benchmarkGet(b, NewBTree(32))
}
func BenchmarkGetBtree64(b *testing.B) {
	benchmarkGet(b, NewBTree(64))
}

// benchmarkAdd creates b.N elements.
func benchmarkAdd(b *testing.B, c Collection) {
	for i := 0; i < b.N; i++ {
		key := newKey(fmt.Sprintf("a%020d", rand.Int63()))
		c.Add(key)
	}
}

func BenchmarkAddSortedSlice(b *testing.B) {
	benchmarkAdd(b, &SortedSlice{})
}
func BenchmarkAddLLRB(b *testing.B) {
	benchmarkAdd(b, NewLLRB())
}
func BenchmarkAddBtree2(b *testing.B) {
	benchmarkAdd(b, NewBTree(2))
}
func BenchmarkAddBtree4(b *testing.B) {
	benchmarkAdd(b, NewBTree(4))
}
func BenchmarkAddBtree8(b *testing.B) {
	benchmarkAdd(b, NewBTree(8))
}
func BenchmarkAddBtree16(b *testing.B) {
	benchmarkAdd(b, NewBTree(16))
}
func BenchmarkAddBtree32(b *testing.B) {
	benchmarkAdd(b, NewBTree(32))
}
func BenchmarkAddBtree64(b *testing.B) {
	benchmarkAdd(b, NewBTree(64))
}

// benchmarkDelete deletes b.N elements.
func benchmarkDelete(b *testing.B, c Collection) {
	b.StopTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		key := newKey(fmt.Sprintf("a%020d", rand.Int63()))
		c.Add(key)
		keys[i] = key
	}
	c.Freeze()
	p := rand.Perm(len(keys))
	b.StartTimer()
	for i, _ := range keys {
		key := keys[p[i]]
		if !bytes.Equal(c.Delete(key), key) {
			b.Fatal("Not found: %v", key)
		}
	}
}

// Too slow
//func BenchmarkDeleteSortedSlice(b *testing.B) {
//	benchmarkDelete(b, &SortedSlice{})
//}
func BenchmarkDeleteLazySortedSlice(b *testing.B) {
	benchmarkDelete(b, &LazySortedSlice{})
}
func BenchmarkDeleteLLRB(b *testing.B) {
	benchmarkDelete(b, NewLLRB())
}
func BenchmarkDeleteBtree2(b *testing.B) {
	benchmarkDelete(b, NewBTree(2))
}
func BenchmarkDeleteBtree4(b *testing.B) {
	benchmarkDelete(b, NewBTree(4))
}
func BenchmarkDeleteBtree8(b *testing.B) {
	benchmarkDelete(b, NewBTree(8))
}
func BenchmarkDeleteBtree16(b *testing.B) {
	benchmarkDelete(b, NewBTree(16))
}
func BenchmarkDeleteBtree32(b *testing.B) {
	benchmarkDelete(b, NewBTree(32))
}
func BenchmarkDeleteBtree64(b *testing.B) {
	benchmarkDelete(b, NewBTree(64))
}
