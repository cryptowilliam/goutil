package gset

import "github.com/emirpasic/gods/sets/hashset"

type HashSetBK struct {
	set *hashset.Set
}

func NewHashSetBK() *HashSetBK {
	var result HashSetBK
	result.set = hashset.New()
	return &result
}

func (s *HashSetBK) Add(items ...interface{}) {
	s.set.Add(items...) //<-- Notice the 3 dot added after items
}

func (s *HashSetBK) Contains(items ...interface{}) bool {
	return s.set.Contains(items...) //<-- Notice the 3 dot added after items
}

func (s *HashSetBK) Remove(items ...interface{}) {
	s.set.Remove(items...) //<-- Notice the 3 dot added after items
}

func (s *HashSetBK) Size() int {
	return s.set.Size()
}

// Clear clears all values in the set.
func (s *HashSetBK) Clear() {
	s.set.Clear()
}
