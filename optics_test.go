package microspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOPTICS(t *testing.T) {
	tt := []struct {
		data     [][]float32
		clusters [][]int
	}{
		{
			data: [][]float32{
				{1, 1}, {0, 1}, {1, 0},
				{10, 10}, {10, 11}, {11, 10},
				{50, 50}, {51, 50}, {50, 51}, {50, 49},
				{100, 100},
			},
			clusters: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8, 9},
			},
		},
	}

	for _, tc := range tt {
		points := []*Point{}
		index := NewAxdex(uint(len(tc.data)))
		for _, coord := range tc.data {
			point := &Point{X: coord[0], Y: coord[1]}
			index.Insert(point)
			points = append(points, point)
		}

		clusters := OPTICS(index, 4, 3)
		assert.Equal(t, len(clusters), len(tc.clusters))
		for i, expectation := range tc.clusters {
			expected := []*Point{}
			for _, idx := range expectation {
				expected = append(expected, points[idx])
			}

			assert.Equal(t, expected, clusters[i].Points)
		}
	}
}
