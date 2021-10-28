package gset

type HashSet struct {
	set mapset.Set
}

func NewHashSet() *HashSet {
	var result HashSet
	result.set = mapset.NewSet()
	return &result
}

func (s *HashSet) Add(item interface{}) {
	s.set.Add(item)
}

func (s *HashSet) Contains(items ...interface{}) bool {
	return s.set.Contains(items...) //<-- Notice the 3 dot added after items
}

func (s *HashSet) Remove(item interface{}) {
	s.set.Remove(item)
}

// Clear clears all values in the set.
func (s *HashSet) Clear() {
	s.set.Clear()
}
