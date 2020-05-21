package cirno

// NormalToCircle returns the normal from the given circle
// to the other circle.
func (circle *Circle) NormalToCircle(other *Circle) Vector {
	return other.center.Subtract(circle.center).Normalize()
}
