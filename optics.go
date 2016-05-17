package microspace

import (
	"container/heap"
	"fmt"
	"sort"
)

// opticPQ is a priority queue where points with a lower reach distance
// are ordered first.
type opticPQ struct {
	points []*opticPoint
}

// Len implements sort.Interface.Len
func (o *opticPQ) Len() int {
	return len(o.points)
}

// Less implements sort.Interface.Less
func (o *opticPQ) Less(a, b int) bool {
	return o.points[a].reachDist < o.points[b].reachDist
}

// Swap implements sort.Interface.Swap
func (o *opticPQ) Swap(a, b int) {
	o.points[a], o.points[b] = o.points[b], o.points[a]
}

// Push implements sort.Interface.Push
func (o *opticPQ) Push(x interface{}) {
	o.points = append(o.points, x.(*opticPoint))
}

// Pop implements sort.Interface.Pop
func (o *opticPQ) Pop() interface{} {
	var first *opticPoint
	first, o.points = o.points[0], o.points[1:]
	return first
}

// Returns the index of the point in the optic queue.
func (o *opticPQ) IndexOf(p *opticPoint) int {
	idx := sort.Search(len(o.points), func(i int) bool {
		return o.points[i].reachDist < p.reachDist
	})

	for idx++; o.points[idx] != p; idx++ {
	}

	return idx
}

var _ heap.Interface = new(opticPQ)

type opticPoint struct {
	point     *Point
	processed bool
	reachDist float32
}

type optics struct {
	index      Index
	points     []*opticPoint
	clusters   []*Cluster
	pointIndex map[*Point]int

	epsilon   float32
	minPoints int
}

type Cluster struct {
	Points []*Point
}

func (c *Cluster) add(point *Point) {
	c.Points = append(c.Points, point)
}

// OPTICS executes a cluster analysis on the spatial index, where
// `epsilon` is the max search distance and `minPoints` is the smallest
// number of individual points needed to qualify as a cluster.
func OPTICS(idx Index, epsilon float32, minPoints int) []*Cluster {
	opts := &optics{
		index:      idx,
		epsilon:    epsilon,
		minPoints:  minPoints,
		pointIndex: make(map[*Point]int),
	}

	points := idx.Points()
	opts.points = make([]*opticPoint, len(points))
	for i, point := range points {
		opts.points[i] = &opticPoint{point: point, reachDist: -1}
		opts.pointIndex[point] = i
	}

	opts.Run(epsilon, minPoints)
	return opts.clusters
}

func (o *optics) Run(epsilon float32, minPoints int) {
	for _, op := range o.points {
		if op.processed {
			continue
		}

		cluster := &Cluster{}
		cluster.add(op.point)
		o.clusters = append(o.clusters, cluster)
		op.processed = true

		cdsq := o.squaredDistanceToCore(op.point)
		if cdsq < 0 {
			continue
		}

		neighbors := o.index.NearestN(op.point, -1, o.epsilon)
		fmt.Printf("N(%s) => %s\n", op.point, neighbors)
		queue := &opticPQ{}
		o.updateQueue(op.point, cdsq, neighbors, queue)
		o.expandCluster(cluster, queue)
	}
}

func (o *optics) getRecordForPoint(p *Point) *opticPoint {
	return o.points[o.pointIndex[p]]
}

func (o *optics) updateQueue(p *Point, cdsq float32, neighbors []*Point, queue *opticPQ) {
	for _, neighbor := range neighbors {
		op := o.getRecordForPoint(neighbor)
		if op.processed {
			continue
		}

		distance := p.DistanceToSqr(op.point)
		if cdsq > distance {
			distance = cdsq
		}

		if op.reachDist < 0 {
			op.reachDist = distance
			heap.Push(queue, op)
		} else if distance < op.reachDist {
			idx := queue.IndexOf(op)
			op.reachDist = distance
			heap.Fix(queue, idx)
		}
	}
}

func (o *optics) expandCluster(cluster *Cluster, queue *opticPQ) {
	for _, op := range queue.points {
		if op.processed {
			continue
		}

		op.processed = true
		cluster.add(op.point)

		cdsq := o.squaredDistanceToCore(op.point)
		if cdsq < 0 {
			continue
		}

		neighbors := o.index.NearestN(op.point, -1, o.epsilon)
		o.updateQueue(op.point, cdsq, neighbors, queue)
		o.expandCluster(cluster, queue)
		return
	}
}

// squaredDistanceToCore returns the square of a point's distance to the
// core of the nearest cluster. It returns -1 if a cluster can't be found
// within the epsilon radius.
func (o *optics) squaredDistanceToCore(p *Point) float32 {
	n := o.index.NearestN(p, o.minPoints, o.epsilon)
	if len(n) == o.minPoints {
		return p.DistanceToSqr(n[len(n)-1])
	}

	return -1
}
