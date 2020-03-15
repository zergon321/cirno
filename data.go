package cirno

type data struct {
	data interface{}
}

func (data *data) Data() interface{} {
	return data.data
}

func (data *data) SetData(newData interface{}) {
	data.data = newData
}