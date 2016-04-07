package rbench

import (
	"fmt"
	"sort"
)

// Point represents a point in two-dimensional space.
type Point struct{ x, y float32 }

// DistanceToSqr returns the squared distance to the `other` point.
func (p *Point) DistanceToSqr(other *Point) float32 {
	dx, dy := (p.x - other.x), (p.y - other.y)
	return dx*dx + dy*dy
}

// String returns a textual representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("(%.4f, %.4f)", p.x, p.y)
}

// axisPoint is used for internal recordkeeping of points within an axis.
// It's a pair of the point and the value of that point's coordinate on
// the related axis.
type axisPoint struct {
	p     *Point
	value float32
}

// axis stores a sorted set of points along a one-dimensional line.
type axis struct {
	data  []axisPoint
	value func(*Point) float32

	generatedIndex bool
	indexed        map[*Point]int
}

// newAxis returns an axis created with the provided capacity. It is assumed
// that the axis will be filled with exactly `capacity` points before
// any other operations are done on it.
func newAxis(capacity uint, value func(*Point) float32) *axis {
	return &axis{
		data:  make([]axisPoint, capacity),
		value: func(p *Point) float32 { return p.y },
	}
}

// IndexFor returns the index of the point on the axis. It's assumed that the
// point will exist in the axis.
func (a *axis) IndexFor(p *Point) int {
	if !a.generatedIndex {
		a.indexed = map[*Point]int{}
		for i, pt := range a.data {
			a.indexed[pt.p] = i
		}
		a.generatedIndex = true
	}

	return a.indexed[p]
}

// ValueFor returns the point's coordinate on that axis.
func (a *axis) ValueFor(p *Point) float32 {
	return a.value(p)
}

// Insert adds a new point to the axis.
func (a *axis) Insert(p *Point) {
	val := a.value(p)
	i := sort.Search(len(a.data), func(i int) bool {
		if a.data[i].p == nil {
			return true
		}

		return a.data[i].value >= val
	})

	// We find the index the item is going to be inserted at, then we shift
	// everything over to make room for it. The list is filled from the left,
	// and we know that we'll never go over capacity.
	copy(a.data[i+1:], a.data[i:])
	a.data[i] = axisPoint{p: p, value: a.value(p)}
}

type Axdex struct {
	axes []*axis
}

// NewAxdex returns a new axis-based index with the provided capacity.
// It's assumed that you will insert exactly `capacity` points before
// running queries against the index.
func NewAxdex(capacity uint) *Axdex {
	a := &Axdex{
		axes: []*axis{
			newAxis(capacity, func(p *Point) float32 { return p.x }),
			newAxis(capacity, func(p *Point) float32 { return p.y }),
		},
	}

	return a
}

// Insert adds a new point into the Axdex.
func (a *Axdex) Insert(p *Point) {
	for _, axis := range a.axes {
		axis.Insert(p)
	}
}

type axResults struct {
	src   *Point
	data  []*Point
	worst float32
	count int
}

// Viable returns true if the provided value could possible be a coordinate
// of a nearest neighbor with coordinate src.
func (a *axResults) Viable(p *Point) (viable bool, distance float32) {
	d := p.DistanceToSqr(a.src)
	if a.data[a.count-1] == nil {
		return true, d
	}

	return d < a.worst, d
}

// HasPotential returns true if the difference between the center point and
// another point, given as delta, is less than the provided max and if it
// could possibly yield a viable point. Once this returns false for an axis
// points further out on that access will not have potential either.
func (a *axResults) HasPotential(delta, max float32) bool {
	if delta > max || -delta > max {
		return false
	}

	if a.data[a.count-1] == nil {
		return true
	}

	return delta*delta < a.worst
}

// GetResult returns a list of results from the list. It will returns as many
// non-nil results as it can, up to the provided count.
func (a *axResults) GetResult() []*Point {
	var i int
	for i < a.count && a.data[i] != nil {
		i++
	}

	return a.data[:i]
}

// Attempts to insert the point into the results.
func (a *axResults) Insert(p *Point) {
	for i := 0; i < a.count; i++ {
		if a.data[i] == p {
			return
		}

		if a.data[i] == nil {
			a.data[i] = p
			break
		}

		if a.src.DistanceToSqr(p) < a.src.DistanceToSqr(a.data[i]) {
			copy(a.data[i+1:], a.data[i:])
			a.data[i] = p
			break
		}
	}

	if a.data[a.count-1] != nil {
		a.worst = a.data[a.count-1].DistanceToSqr(a.src)
	}
}

// NearestN returns up the `n` nearest neighbors of the point, with a `max`
// search distance.
func (a *Axdex) NearestN(p *Point, n int, max float32) []*Point {
	results := &axResults{src: p, data: make([]*Point, n), count: n}
	results.Insert(p)

	for _, axis := range a.axes {
		idx := axis.IndexFor(p)
		var (
			left  = idx - 1
			right = idx + 1
			value = axis.ValueFor(p)
		)

		for {
			var (
				leftP  axisPoint
				rightP axisPoint

				leftViable  = false
				rightViable = false

				leftDistance  = float32(0)
				rightDistance = float32(0)
			)

			if left >= 0 {
				leftP = axis.data[left]
				leftViable, leftDistance = results.Viable(leftP.p)
				if !leftViable {
					left--
				}
			}
			if right < len(axis.data) {
				rightP = axis.data[right]
				rightViable, rightDistance = results.Viable(rightP.p)
				if !rightViable {
					right++
				}
			}

			if leftViable && (!rightViable || leftDistance < rightDistance) {
				results.Insert(leftP.p)
				left--
			} else if rightViable {
				results.Insert(rightP.p)
				right++
			}

			leftPotential := left >= 0 && results.HasPotential(value-leftP.value, max)
			rightPotential := right < len(axis.data) && results.HasPotential(value-rightP.value, max)
			if !(leftPotential || rightPotential) {
				break
			}

			if !leftPotential {
				left = -1
			}
			if !rightPotential {
				right = len(axis.data)
			}
		}
	}

	return results.GetResult()
}
