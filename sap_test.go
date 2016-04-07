package rbench

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pointDistanceList struct {
	center *Point
	list   []*Point
}

func (p pointDistanceList) Len() int {
	return len(p.list)
}

func (p pointDistanceList) Less(i, j int) bool {
	return p.list[i].DistanceToSqr(p.center) < p.list[j].DistanceToSqr(p.center)
}

func (p pointDistanceList) Swap(i, j int) {
	p.list[i], p.list[j] = p.list[j], p.list[i]
}

func TestIndexNearest(t *testing.T) {
	count := 100
	delta := 0.000001
	tr := NewAxdex(uint(count))

	points := []*Point{}
	for i := 0; i < count; i++ {
		p := &Point{rand.Float32(), rand.Float32()}
		points = append(points, p)
		tr.Insert(p)
	}

	finalizeIndex(tr)

	for _, axis := range tr.axes {
		for i := 1; i < len(axis.data); i++ {
			assert.True(t, axis.ValueFor(axis.data[i].p) >= axis.ValueFor(axis.data[i-1].p))
			assert.Equal(t, i, axis.IndexFor(axis.data[i].p))
		}
	}

	testLast := 5
	for _, p := range points {
		n := tr.NearestN(p, testLast, 0.25)
		pdl := pointDistanceList{center: p, list: points}
		sort.Sort(pdl)

		list := pdl.list[:testLast]
		if len(n) < testLast || len(list) < testLast {
			t.Fatalf("Invalid nearest for point %s:\n\tResults:   %s\n\tExpecting: %s\n", p, n, list)
		}

		for k := 0; k < testLast; k++ {
			if math.Abs(float64(list[k].x-n[k].x)) > delta || math.Abs(float64(list[k].y-n[k].y)) > delta {
				t.Fatalf("Invalid nearest for point %s:\n\tResults:   %s\n\tExpecting: %s\n\tGot: %s, expected %s\n", p, n, list, n[k], list[k])
			}
		}
	}
}

func finalizeIndex(t *Axdex) {
	for _, axis := range t.axes {
		axis.runSort()
	}
}

func generateIndex(n int) *Axdex {
	t := NewAxdex(uint(n))
	for k := 0; k < n; k++ {
		t.Insert(&Point{rand.Float32(), rand.Float32()})
	}

	return t
}

func benchIndexCreate(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		finalizeIndex(generateIndex(n))
	}
}

func benchIndexNearest(b *testing.B, n int) {
	t := generateIndex(n)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.NearestN(&Point{0.5, 0.5}, 3, 0.25)
	}
}

func benchIndexNearestWorstCase(b *testing.B, n int) {
	t := NewAxdex(uint(n))
	for k := 0; k < n; k++ {
		t.Insert(&Point{0.6, 0.6})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.NearestN(&Point{0.5, 0.5}, 3, 0.25)
	}
}

func BenchmarkIndexCreate10(b *testing.B)    { benchIndexCreate(b, 10) }
func BenchmarkIndexCreate100(b *testing.B)   { benchIndexCreate(b, 100) }
func BenchmarkIndexCreate1000(b *testing.B)  { benchIndexCreate(b, 1000) }
func BenchmarkIndexCreate10000(b *testing.B) { benchIndexCreate(b, 10000) }

func BenchmarkIndexNearest10(b *testing.B)    { benchIndexNearest(b, 10) }
func BenchmarkIndexNearest100(b *testing.B)   { benchIndexNearest(b, 100) }
func BenchmarkIndexNearest1000(b *testing.B)  { benchIndexNearest(b, 1000) }
func BenchmarkIndexNearest10000(b *testing.B) { benchIndexNearest(b, 10000) }

func BenchmarkIndexNearestWorstCase10(b *testing.B)    { benchIndexNearestWorstCase(b, 10) }
func BenchmarkIndexNearestWorstCase100(b *testing.B)   { benchIndexNearestWorstCase(b, 100) }
func BenchmarkIndexNearestWorstCase1000(b *testing.B)  { benchIndexNearestWorstCase(b, 1000) }
func BenchmarkIndexNearestWorstCase10000(b *testing.B) { benchIndexNearestWorstCase(b, 10000) }
