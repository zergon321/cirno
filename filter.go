package cirno

// ОШИБКА ОБРАЩЕНИЯ К ДАННЫМ: ни в одной из функций данного файла
// не происходит проверка аргумента на nil (пустой указатель).

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
func (shapes Shapes) FilterByCollisionRight(shape Shape) Shapes {
	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		if shape.ShouldCollide(item) {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes
}

// FilterByCollisionLeft returns all the shapes that should
// collide the given shape.
func (shapes Shapes) FilterByCollisionLeft(shape Shape) Shapes {
	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		if item.ShouldCollide(shape) {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes
}

// FilterByCollision returns all the shapes that should collide
// or get collided by the given shape.
func (shapes Shapes) FilterByCollision(shape Shape) Shapes {
	filteredShapes := make(Shapes, 0)

	for item := range shapes {
		if item.ShouldCollide(shape) || shape.ShouldCollide(item) {
			filteredShapes.Insert(item)
		}
	}

	return filteredShapes
}
