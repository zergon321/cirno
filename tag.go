package cirno

import "fmt"

// tag identifies the kind of the shape
// to collide other shapes.
type tag struct {
	// Identity determines which other
	// shapes can collide the present shape.
	identity int32
	// Mask determines which other shapes
	// the present shape can collide.
	mask int32
}

// ShouldCollide returns true if the shape should
// collide another one accodring to their tags.
func (t tag) ShouldCollide(other Shape) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the shape is nil")
	}

	return t.mask&other.GetIdentity() > 0, nil
}

// GetIdentity returns the valye of the shape identity.
func (t tag) GetIdentity() int32 {
	return t.identity
}

// SetIdentity assigns a new value to the tag identity.
func (t *tag) SetIdentity(newIdentity int32) {
	t.identity = newIdentity
}

// GetMask the value of the shape mask.
func (t tag) GetMask() int32 {
	return t.mask
}

// SetMask assigns a new value to the tag mask.
func (t *tag) SetMask(newMask int32) {
	t.mask = newMask
}
