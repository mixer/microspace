package rbench

import (
	// "math"
	// "sort"
	"fmt"
)

type Point struct{ x, y float32 }

func (p Point) DistanceTo(other Point) float32 {
	dx, dy := (p.x - other.x), (p.y - other.y)
	return dx*dx + dy*dy
}

func (p Point) String() string {
	return fmt.Sprintf("(%.4f, %.4f)", p.x, p.y)
}

// type Tree struct {
// 	// The data is a list of buckets of points. The top left corner of the
// 	// bucket can be found at (index % resolution * xbucket,
// 	// index / resolution * ybucket).
// 	data [][]Point

// 	resolution           int
// 	topLeft, bottomRight Point
// 	xbucket, ybucket     float64
// }

// // NewTree creates a new grid-like structure in the provided bounds
// // range, with the give resolution. Greater resolutions will lend better
// // performance at the cost of more memory usage.
// func NewTree(topLeft, bottomRight Point, resolution uint) *Tree {
// 	return &Tree{
// 		data:        make([][]Point, (resolution+1)*(resolution+1)),
// 		topLeft:     topLeft,
// 		bottomRight: bottomRight,
// 		resolution:  int(resolution),
// 		xbucket:     (bottomRight.x - topLeft.x) / float64(resolution),
// 		ybucket:     (bottomRight.y - topLeft.y) / float64(resolution),
// 	}
// }

// // inBounds returns whether the x value is within the left/right bounds.
// func inBounds(left, x, right float64) bool {
// 	if x < left {
// 		return false
// 	}
// 	if x >= right {
// 		return false
// 	}

// 	return true
// }

// // absMax returns the absolute maximum of the provided values,
// func absMax(values ...float64) float64 {
// 	max := values[0]
// 	for i := 0; i < len(values); i++ {
// 		a := math.Abs(values[i])
// 		if a > max {
// 			max = a
// 		}
// 	}

// 	return max
// }

// // coordsToIndex returns the index of the bucket that the provided coordinates
// // fall into.
// func (q Tree) coordsToIndex(x, y float64) int {
// 	return int((x-q.topLeft.x)/q.xbucket) + int((y-q.topLeft.y)/q.ybucket)*q.resolution
// }

// // Insert adds a new point to the Tree.
// func (q Tree) Insert(p Point) {
// 	idx := q.coordsToIndex(p.x, p.y)
// 	q.data[idx] = append(q.data[idx], p)
// }

type pointDistanceList struct {
	center *Point
	list   []*Point
}

func (p pointDistanceList) Len() int {
	return len(p.list)
}

func (p pointDistanceList) Less(i, j int) bool {
	return p.PointsLess(p.list[i], p.list[j])
}

func (p pointDistanceList) PointsLess(i, j *Point) bool {
	return i.DistanceTo(*p.center) <= j.DistanceTo(*p.center)
}

func (p pointDistanceList) Swap(i, j int) {
	p.list[i], p.list[j] = p.list[j], p.list[i]
}

func (p *pointDistanceList) Merge(r []*Point) {
	next := make([]*Point, 0, len(p.list)+len(r))
	for len(p.list) > 0 || len(r) > 0 {
		if len(p.list) == 0 {
			next = append(next, r...)
			break
		}
		if len(r) == 0 {
			next = append(next, p.list...)
			break
		}

		if p.PointsLess(p.list[0], r[0]) {
			next = append(next, p.list[0])
			p.list = p.list[1:]
		} else {
			next = append(next, r[0])
			p.list = p.list[1:]
		}
	}

	p.list = next
}

// // NearestN the approximate nearest n int integers to the provided point.
// func (q Tree) NearestN(p Point, n int) []Point {
// 	results := &pointDistanceList{center: p, list: make([]Point, 0, n)}

// 	center := q.coordsToIndex(p.x, p.y)
// 	left := (center / q.resolution) * q.resolution
// 	right := left + q.resolution
// 	candidates := pointDistanceList{
// 		center: p,
// 		list:   q.data[center][:],
// 	}

// 	for epsilon := 1; epsilon < q.resolution; epsilon++ {
// 		for xi := -epsilon; xi <= epsilon; xi++ {
// 			x := center + xi
// 			if x < left || x >= right {
// 				continue
// 			}

// 			var y int
// 			if xi < 0 {
// 				y = (epsilon + xi) * q.resolution
// 			} else {
// 				y = (epsilon - xi) * q.resolution
// 			}

// 			if x+y < len(q.data) {
// 				candidates.list = append(candidates.list, q.data[x+y]...)
// 			}
// 			if y > 0 && x-y >= 0 {
// 				candidates.list = append(candidates.list, q.data[x-y]...)
// 			}
// 		}

// 		sort.Sort(candidates)

// 		// If we have at least how many candidates the called asked for,
// 		// add them to the results and stop looping. We're done.
// 		if candidates.Len() > n-results.Len() {
// 			results.Merge(candidates.list[:n-results.Len()])
// 			break
// 		}

// 		// Otherwise merge what we have, clear the list, and look further.
// 		results.Merge(candidates.list)
// 		candidates.list = []Point(nil)
// 	}

// 	return results.list
// }
