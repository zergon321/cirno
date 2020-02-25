package cirno

// Tag identifies the kind of the shape
// to collide other shapes.
type Tag int32

// ShouldCollide returns true if the shape should
// collide another one accodring to their tags.
func (tag Tag) ShouldCollide(other Shape) bool {
	return tag&other.GetTag() > 0
}

// SetTag assign a new value to the tag.
func (tag *Tag) SetTag(newTag Tag) {
	*tag = newTag
}

// GetTag returns the value if the tag.
func (tag Tag) GetTag() Tag {
	return tag
}
