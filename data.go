package cirno

// data contains any data
// assigned to the shape.
type data struct {
	data interface{}
}

// Data returns whatever is saved in the
// data field of the shape.
func (data *data) Data() interface{} {
	return data.data
}

// SetData assigns new data to the shape.
func (data *data) SetData(newData interface{}) {
	data.data = newData
}
