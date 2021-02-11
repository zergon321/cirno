package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestCreateRectangle(t *testing.T) {
	_, err := cirno.NewRectangle(cirno.NewVector(0, 0), 0, 13, 0)
	assert.NotNil(t, err)
	_, err = cirno.NewRectangle(cirno.NewVector(0, 0), 13, 0, 0)
	assert.NotNil(t, err)
}
