package collection

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/google/btree"
	"github.com/petar/GoLLRB/llrb"
)

// Collection defines common interface for slice, llrb, and btree.
type Collection interface {
	Add([]byte)
	Get([]byte) []byte
}

type BytesSlice [][]byte

// Implement sort.Interface.
func (bs BytesSlice) Len() int {
	return len(bs)
}

func (bs BytesSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (bs BytesSlice) Less(i, j int) bool {
	return bytes.Compare(bs[i], bs[j]) < 0
}

// SortedSlice implements the Collection interface with sorted slice.
type SortedSlice struct {
	bs BytesSlice
}

func (s *SortedSlice) Add(v []byte) {
	s.bs = append(s.bs, v)
	sort.Sort(s.bs)
}

func (s *SortedSlice) Get(k []byte) []byte {
	n := sort.Search(len(s.bs), func(i int) bool {
		return bytes.Compare(s.bs[i], k) >= 0
	})
	if n >= len(s.bs) || !bytes.Equal(s.bs[n], k) {
		return nil
	}
	return s.bs[n]
}

// LLRB implements the Collection interface with llrb.LLRB.
type LLRB struct {
	lt *llrb.LLRB
}

func NewLLRB() *LLRB {
	return &LLRB{lt: llrb.New()}
}

type LBytesItem []byte

// Implement llrb.Item.
func (a LBytesItem) Less(b llrb.Item) bool {
	return bytes.Compare(a, b.(LBytesItem)) < 0
}

func (t *LLRB) Add(v []byte) {
	t.lt.ReplaceOrInsert(LBytesItem(v))
}

func (t *LLRB) Get(k []byte) []byte {
	v := t.lt.Get(LBytesItem(k))
	if v == nil {
		return nil
	}
	return v.(LBytesItem)
}

// BTree implements the Collection interface with btree.Btree.
type BTree struct {
	bt *btree.BTree
}

func NewBTree(degree int) *BTree {
	return &BTree{bt: btree.New(degree)}
}

type BBytesItem []byte

// Implement btree.Item.
func (a BBytesItem) Less(b btree.Item) bool {
	return bytes.Compare(a, b.(BBytesItem)) < 0
}

func (t *BTree) Add(v []byte) {
	t.bt.ReplaceOrInsert(BBytesItem(v))
}

func (t *BTree) Get(k []byte) []byte {
	v := t.bt.Get(BBytesItem(k))
	if v == nil {
		return nil
	}
	return v.(BBytesItem)
}

func main() {
	//	str := "abc"

	{
		var s SortedSlice
		key := []byte("abc")
		s.Add(key)
		fmt.Println(string(s.Get(key)))

		fmt.Println(s.bs)

		k := []byte("123")
		s.Add(k)
		fmt.Println(string(s.Get(k)))

		key = []byte("cdf")
		fmt.Println(string(s.Get(key)))

		s.Add(key)
		fmt.Println(string(s.Get(key)))

		fmt.Println("=====================")

		fmt.Println(string(s.bs[0]))
		fmt.Println(string(s.bs[1]))
		fmt.Println(string(s.bs[2]))
	}

	fmt.Println("=====================")

	{
		t := NewLLRB()
		key := []byte("abc")
		t.Add(key)
		fmt.Println(string(t.Get(key)))

		key = []byte("cdf")
		fmt.Println(string(t.Get(key)))

		t.Add(key)
		fmt.Println(string(t.Get(key)))
	}

	fmt.Println("=====================")

	{
		t := NewBTree(2)
		key := []byte("abc")
		t.Add(key)
		fmt.Println(string(t.Get(key)))

		key = []byte("cdf")
		fmt.Println(string(t.Get(key)))

		t.Add(key)
		//		fmt.Println(string(t.Get(key)))
	}
}
