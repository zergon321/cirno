package cirno

import (
	"fmt"
)

// FilterByIdentity returns all the shapes
// matching the specified identity template.
func (shapes Shapes) FilterByIdentity(identity int32) Shapes {
	filteredShapes := make(Shapes, 0)

	for shape := range shapes {
		if shape.GetIdentity()&identity == identity {
			filteredShapes.Insert(shape)
		}
	}

	return filteredShapes
}

// FilterByMask returns all the shapes
// matching the specified mask template.
func (shapes Shapes) FilterByMask(mask int32) Shapes {
	filteredShapes := make(Shapes, 0)

	for shape := range shapes {
		if shape.GetMask()&mask == mask {
			filteredShapes.Insert(shape)
		}
	}

	return filteredShapes
}

// FilterByCollisionRight returns all the shapes the given shape
// should collide.
func (shapes Shapes) FilterByCollisionRight(shape Shape) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		shouldCollide, err := shape.ShouldCollide(item)

		if err != nil {
			return nil, err
		}

		if shouldCollide {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes, nil
}

// FilterByCollisionLeft returns all the shapes that should
// collide the given shape.
func (shapes Shapes) FilterByCollisionLeft(shape Shape) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		shouldCollide, err := item.ShouldCollide(shape)

		if err != nil {
			return nil, err
		}

		if shouldCollide {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes, nil
}

// FilterByCollision returns all the shapes that should collide
// or get collided by the given shape.
func (shapes Shapes) FilterByCollision(shape Shape) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		shouldCollideLeft, err := item.ShouldCollide(shape)

		if err != nil {
			return nil, err
		}

		shouldCollideRight, err := shape.ShouldCollide(item)

		if err != nil {
			return nil, err
		}

		if shouldCollideLeft || shouldCollideRight {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes, nil
}
