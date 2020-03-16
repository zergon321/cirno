package cirno

import (
	"github.com/golang-collections/collections/queue"
)

// Raycast casts a ray in the space and returns the hit shape closest
// to the origin of the ray.
//
// Ray cannot hit against the shape within which it's located.
func (space *Space) Raycast(origin Vector, direction Vector, distance float64, mask int32) Shape {
	if distance <= 0 {
		distance = Distance(space.min, space.max)
	}

	ray := NewLine(origin, origin.Add(direction.Normalize().MultiplyByScalar(distance)))
	ray.SetMask(mask)
	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	minExists := false

	var (
		minSquaredDistance float64
		hitShape           Shape
	)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		if !IntersectionLineToRectangle(ray, node.boundary) {
			continue
		}

		// If the node is a leaf.
		if node.northWest == nil {
			for shape := range node.shapes {
				if ResolveCollision(ray, shape, space.useTags) && !shape.ContainsPoint(ray.p) {
					hitDistance := SquaredDistance(ray.p, shape.Center())

					if !minExists {
						minSquaredDistance = hitDistance
						hitShape = shape

						minExists = true
					} else if hitDistance < minSquaredDistance {
						minSquaredDistance = hitDistance
						hitShape = shape
					}
				}
			}
		} else {
			nodeQueue.Enqueue(node.northEast)
			nodeQueue.Enqueue(node.northWest)
			nodeQueue.Enqueue(node.southEast)
			nodeQueue.Enqueue(node.southWest)
		}
	}

	return hitShape
}

// Boxcast casts a box in the space and returns all the
// shapes overlapped by this box.
func (space *Space) Boxcast(rect *Rectangle) Shapes {
	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	shapes := make(Shapes, 0)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		if !CollisionRectangleToRectangle(rect, node.boundary) {
			continue
		}

		if node.northWest == nil {
			for shape := range node.shapes {
				if ResolveCollision(rect, shape, space.useTags) {
					shapes.Insert(shape)
				}
			}
		}
	}

	return shapes
}

// Circlecast casts a circle in the space and returns all the
// shapes overlapped by the circle.
func (space *Space) Circlecast(circle *Circle) Shapes {
	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	shapes := make(Shapes, 0)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		if !CollisionRectangleToCircle(node.boundary, circle) {
			continue
		}

		if node.northWest == nil {
			for shape := range node.shapes {
				if ResolveCollision(circle, shape, space.useTags) {
					shapes.Insert(shape)
				}
			}
		}
	}

	return shapes
}