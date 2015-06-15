package smallset

import (
	"math/rand"
	"testing"
	"testing/quick"
)

func TestSetQuick(t *testing.T) {
	if err := quick.Check(testSet, &quick.Config{MaxCount: 1000000}); err != nil {
		t.Error(err)
	}
}

func testSet(values []uint64) bool {
	mapSet := make(map[uint64]struct{}, len(values))
	for _, v := range values {
		mapSet[v] = struct{}{}
	}

	s := NewSet(len(values))
	for _, v := range values {
		s.Add(v)
	}

	if len(mapSet) != s.Len() {
		return false
	}
	for v := range mapSet {
		if !s.Add(v) {
			return false
		}
	}
	return true
}

func BenchmarkSet(b *testing.B) {
	data := generateData(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	s := NewSet(len(data) + len(data)/8)
	var aggregatedWas bool
	for _, v := range data {
		was := s.Add(v)
		aggregatedWas = aggregatedWas || was
	}
	b.StopTimer()
	b.Log("aggregatedWas is", aggregatedWas)
}

func BenchmarkMap(b *testing.B) {
	data := generateData(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	s := make(map[uint64]struct{}, len(data))
	var aggregatedWas bool
	for _, v := range data {
		_, was := s[v]
		s[v] = struct{}{}
		aggregatedWas = aggregatedWas || was
	}
	b.StopTimer()
	b.Log("aggregatedWas is", aggregatedWas) // show value so it will not be optimized away
}

func generateData(size int) []uint64 {
	data := make([]uint64, size)
	for n := range data {
		data[n] = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	}
	return data
}
