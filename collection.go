package collection

import (
	"bytes"
	"sort"

	"github.com/google/btree"
	"github.com/petar/GoLLRB/llrb"
)

// Collection defines common interface for slice, llrb, and btree.
type Collection interface {
	Add([]byte)
	Get([]byte) []byte
	Delete([]byte) []byte
	// Freeze is called when no further mutation will not
	// happen. No-op except LazySortedSlice.
	Freeze()
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

func (s *SortedSlice) Delete(k []byte) []byte {
	n := sort.Search(len(s.bs), func(i int) bool {
		return bytes.Compare(s.bs[i], k) >= 0
	})
	if n >= len(s.bs) || !bytes.Equal(s.bs[n], k) {
		return nil
	}
	v := s.bs[n]
	s.bs = append(s.bs[:n], s.bs[n+1:]...)
	return v
}

func (s *SortedSlice) Freeze() {
}

// LazySortedSlice implements the Collection interface and
// do sorting when requested.
type LazySortedSlice struct {
	bs     BytesSlice
	frozen bool
}

func (s *LazySortedSlice) Add(v []byte) {
	if s.frozen {
		panic("No mutation is allowed")
	}
	s.bs = append(s.bs, v)
}

func (s *LazySortedSlice) Get(k []byte) []byte {
	if !s.frozen {
		panic("Not yet frozen")
	}
	n := sort.Search(len(s.bs), func(i int) bool {
		return bytes.Compare(s.bs[i], k) >= 0
	})
	if n >= len(s.bs) || !bytes.Equal(s.bs[n], k) {
		return nil
	}
	return s.bs[n]
}

func (s *LazySortedSlice) Delete(k []byte) []byte {
	n := sort.Search(len(s.bs), func(i int) bool {
		return bytes.Compare(s.bs[i], k) >= 0
	})
	if n >= len(s.bs) || !bytes.Equal(s.bs[n], k) {
		return nil
	}
	v := s.bs[n]
	s.bs = append(s.bs[:n], s.bs[n+1:]...)
	return v
}

func (s *LazySortedSlice) Freeze() {
	sort.Sort(s.bs)
	s.frozen = true
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

func (t *LLRB) Delete(k []byte) []byte {
	v := t.lt.Delete(LBytesItem(k))
	if v == nil {
		return nil
	}
	return v.(LBytesItem)
}

func (t *LLRB) Freeze() {
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

func (t *BTree) Delete(k []byte) []byte {
	v := t.bt.Delete(BBytesItem(k))
	if v == nil {
		return nil
	}
	return v.(BBytesItem)
}

func (t *BTree) Freeze() {
}
