//Append-only set of uint64 with tiny overhead.
//
//This is just append-only, fixed size set of uint64. It is faster than map and has lower memory footprint
//(essentially as slice of uint64 - 8 bytes per value, **but to gain comparable speed you should choose how many
//extra slots will be allocated**),
//It can be used only for uniformly distributed uint64 values. **It will panic if there are no slots left.**
//This is essentially hash table with open addressing, but without hash and with dead-simple addressing logic (lookup in nearest slots).
//Not concurrent-safe.
package smallset

const empty = uint64(0)

type Set struct {
	items        []uint64
	itemsPerCell uint64
	hasEmpty     bool
	itemsNum     int
}

func NewSet(size int) *Set {
	if size < 0 {
		panic("Wrong size for set (it should be greater or equals to zero)")
	}
	if size < 2 {
		size = 2
	}
	return &Set{
		items:        make([]uint64, size),
		itemsPerCell: (uint64(1) << 63) / uint64(size) * 2,
	}
}

// Add adds value to set and returns true if value was in Set before.
func (s *Set) Add(v uint64) bool {
	if v == empty {
		if s.hasEmpty {
			return true
		}
		s.itemsNum++
		s.hasEmpty = true
		return false
	}
	pos := s.getPosition(v)
	if s.items[pos] == v {
		return true
	}
	if s.items[pos] != empty {
		newpos, found := s.findRight(pos, v)
		if newpos < 0 {
			newpos, found := s.findLeft(pos, v)
			if newpos < 0 {
				panic("no space!")
			}
			if found {
				return true
			}
			pos = newpos
		} else {
			if found {
				return true
			}
			pos = newpos
		}
	}
	s.items[pos] = v
	s.itemsNum++
	return false
}

func (s *Set) Len() int {
	return s.itemsNum
}

func (s *Set) Cap() int {
	return len(s.items)
}

func (s *Set) findLeft(pos int, v uint64) (int, bool) {
	for n := pos - 1; n >= 0; n-- {
		if s.items[n] == v {
			return n, true
		} else {
			if s.items[n] == empty {
				return n, false
			}
		}
	}
	return -1, false
}

func (s *Set) findRight(pos int, v uint64) (int, bool) {
	for n, curv := range s.items[pos+1:] {
		if curv == v {
			return n + pos + 1, true
		} else {
			if curv == empty {
				return n + pos + 1, false
			}
		}
	}
	return -1, false
}

func (s *Set) getPosition(v uint64) int {
	return int(v / s.itemsPerCell)
}
