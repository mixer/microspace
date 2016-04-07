package rbench

import (
	"sort"
)

type Axis struct {
	data    []*Point
	indexed map[*Point]int
	value   func(*Point) float32
}

func (a *Axis) IndexFor(p *Point) int {
	return a.indexed[p]
}

func (a *Axis) ValueFor(p *Point) float32 {
	return a.value(p)
}

func (a *Axis) Insert(p *Point) {
	val := a.value(p)
	i := sort.Search(len(a.data), func(i int) bool {
		if a.data[i] == nil {
			return true
		}

		return a.value(a.data[i]) >= val
	})

	// We find the index the item is going to be inserted at, then we shift
	// everything over to make room for it. The list is filled from the left,
	// and we know that we'll never go over capacity.
	copy(a.data[i+1:], a.data[i:])
	a.data[i] = p

	a.indexed[p] = i
}

type Axdex struct {
	axes []*Axis
}

// Creates a new axis-based index with a capacity.
func NewAxdex(capacity uint) *Axdex {
	a := &Axdex{
		axes: []*Axis{
			&Axis{data: make([]*Point, capacity), indexed: map[*Point]int{}, value: func(p *Point) float32 { return p.x }},
			&Axis{data: make([]*Point, capacity), indexed: map[*Point]int{}, value: func(p *Point) float32 { return p.y }},
		},
	}

	return a
}

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

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// Viable returns true if the provided value could possible be a coordinate
// of a nearest neighbor with coordinate src.
func (a *axResults) Viable(p *Point) (viable bool, distance float32) {
	d := p.DistanceTo(*a.src)
	if a.data[a.count-1] == nil {
		return true, d
	}

	return d < a.worst, d
}

func (a *axResults) HasPotential(src, pt float32) bool {
	if a.data[a.count-1] == nil {
		return true
	}

	return abs(src-pt) < a.worst
}

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

		if a.src.DistanceTo(*p) < a.src.DistanceTo(*a.data[i]) {
			copy(a.data[i+1:], a.data[i:])
			a.data[i] = p
			break
		}
	}

	if a.data[a.count-1] != nil {
		a.worst = a.data[a.count-1].DistanceTo(*a.src)
	}
}

type axisTracker struct {
	left, right int
	value       float32
}

func (a *Axdex) NearestN(p *Point, n int) []*Point {
	results := &axResults{src: p, data: make([]*Point, n), count: n}
	results.Insert(p)

	ats := make([]axisTracker, len(a.axes))
	for i, axis := range a.axes {
		idx := axis.IndexFor(p)
		ats[i] = axisTracker{
			left:  idx - 1,
			right: idx + 1,
			value: axis.ValueFor(p),
		}
	}

	for i, axis := range a.axes {
		idx := axis.IndexFor(p)
		var (
			left  = idx - 1
			right = idx + 1
			value = axis.ValueFor(p)
		)

		for {
			var (
				leftViable  = false
				rightViable = false

				leftDistance  = float32(0)
				rightDistance = float32(0)
			)

			if left >= 0 {
				leftViable, leftDistance = results.Viable(axis.data[left])
				if !leftViable {
					left--
				}
			}
			if right < len(axis.data) {
				rightViable, rightDistance = results.Viable(axis.data[right])
				if !rightViable {
					right++
				}
			}

			if leftViable && (!rightViable || leftDistance < rightDistance) {
				results.Insert(axis.data[left])
				left--
			} else if rightViable {
				results.Insert(axis.data[right])
				right++
			}

			leftPotential := left >= 0 && results.HasPotential(value, axis.ValueFor(axis.data[left]))
			rightPotential := right < len(axis.data) && results.HasPotential(value, axis.ValueFor(axis.data[right]))
			if !(leftPotential || rightPotential) {
				break
			}

			if !leftPotential {
				ats[i].left = -1
			}
			if !rightPotential {
				ats[i].right = len(axis.data)
			}
		}
	}

	return results.GetResult()
}
