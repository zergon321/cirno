package cirno

// Tag identifies the kind of the shape
// to collide other shapes.
type Tag struct {
	identity int32
	mask     int32
}

// ShouldCollide returns true if the shape should
// collide another one accodring to their tags.
func (tag Tag) ShouldCollide(other Shape) bool {
	return tag.mask&other.GetIdentity() > 0
}

// GetIdentity returns the valye of the shape identity.
func (tag Tag) GetIdentity() int32 {
	return tag.identity
}

// SetIdentity assigns a new value to the tag identity.
func (tag *Tag) SetIdentity(newIdentity int32) {
	tag.identity = newIdentity
}

// GetMask the value of the shape mask.
func (tag Tag) GetMask() int32 {
	return tag.mask
}

// SetMask assigns a new value to the tag mask.
func (tag *Tag) SetMask(newMask int32) {
	tag.mask = newMask
}
