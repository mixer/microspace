package rbench

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestTreeNearestEmpty(t *testing.T) {
	// assert.Equal(t, []Point{}, generateTree(0).NearestN(&Point{0, 0}, 5))
}

func TestTreeNearest(t *testing.T) {
	inc := 0.01
	delta := (0.000001)
	tr := NewAxdex(uint(1 / inc))

	previous := []*Point{}
	for i := float64(0); i < 1; i += inc {
		p := &Point{float32(i), float32(math.Sin(i * math.Pi))}
		previous = append(previous, p)
		tr.Insert(p)
	}

	testLast := 5
	for _, p := range previous {
		n := tr.NearestN(p, testLast)
		pdl := pointDistanceList{center: p, list: previous}
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

func generateTree(n int) *Axdex {
	t := NewAxdex(uint(n))
	for k := 0; k < n; k++ {
		t.Insert(&Point{rand.Float32(), rand.Float32()})
	}

	return t
}

func benchTreeCreate(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		generateTree(n)
	}
}

func benchTreeNearest(b *testing.B, n int) {
	t := generateTree(n)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.NearestN(&Point{0.5, 0.5}, 3)
	}
}

func benchTreeNearestWorstCase(b *testing.B, n int) {
	t := NewAxdex(uint(n))
	for k := 0; k < n; k++ {
		t.Insert(&Point{0.5, 0.5})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.NearestN(&Point{0.5, 0.5}, 3)
	}
}

func BenchmarkTreeCreate10(b *testing.B)    { benchTreeCreate(b, 10) }
func BenchmarkTreeCreate100(b *testing.B)   { benchTreeCreate(b, 100) }
func BenchmarkTreeCreate1000(b *testing.B)  { benchTreeCreate(b, 1000) }
func BenchmarkTreeCreate10000(b *testing.B) { benchTreeCreate(b, 10000) }

func BenchmarkTreeNearest10(b *testing.B)    { benchTreeNearest(b, 10) }
func BenchmarkTreeNearest100(b *testing.B)   { benchTreeNearest(b, 100) }
func BenchmarkTreeNearest1000(b *testing.B)  { benchTreeNearest(b, 1000) }
func BenchmarkTreeNearest10000(b *testing.B) { benchTreeNearest(b, 10000) }

func BenchmarkTreeNearestWorstCase10(b *testing.B)    { benchTreeNearestWorstCase(b, 10) }
func BenchmarkTreeNearestWorstCase100(b *testing.B)   { benchTreeNearestWorstCase(b, 100) }
func BenchmarkTreeNearestWorstCase1000(b *testing.B)  { benchTreeNearestWorstCase(b, 1000) }
func BenchmarkTreeNearestWorstCase10000(b *testing.B) { benchTreeNearestWorstCase(b, 10000) }
