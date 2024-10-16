package miser

import (
	"testing"
)

func BenchmarkID(b *testing.B) {
	var ids map[ID]struct{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			newID := CreateID()
			id, ok := ids[newID]
			if ok {
				b.Fatal("duplicate id found:", id)
			}
		}
	})
}
