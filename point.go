package microspace

import "fmt"

// Point represents a point in two-dimensional space.
type Point struct{ X, Y float32 }

// DistanceToSqr returns the squared distance to the `other` point.
func (p *Point) DistanceToSqr(other *Point) float32 {
	dx, dy := (p.X - other.X), (p.Y - other.Y)
	return dx*dx + dy*dy
}

// String returns a textual representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("(%.4f, %.4f)", p.X, p.Y)
}
