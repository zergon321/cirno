package cirno

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
)

// Raycast casts a ray in the space and returns the hit shape closest
// to the origin of the ray.
//
// Ray cannot hit against the shape within which it's located.
func (space *Space) Raycast(origin, direction Vector, distance float64, mask int32) (Shape, Vector, error) {
	if Sign(direction.X) == 0 && Sign(direction.Y) == 0 {
		return nil, Zero(), fmt.Errorf("the direction vector is Zero()")
	}

	if distance <= 0 {
		distance = Distance(space.min, space.max)
	}

	normDir, err := direction.Normalize()

	if err != nil {
		return nil, Zero(), err
	}

	ray, err := NewLine(origin, origin.Add(
		normDir.MultiplyByScalar(distance)))

	if err != nil {
		return nil, Zero(), err
	}

	ray.SetMask(mask)
	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	minExists := false

	var (
		minSquaredDistance float64
		hitShape           Shape
		hit                Vector
	)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)
		overlapped, err := node.boundary.collidesLine(ray)

		if err != nil {
			return nil, Zero(), err
		}

		if !overlapped {
			continue
		}

		// If the node is a leaf.
		if node.northWest == nil {
			for shape := range node.shapes {
				raycastHit, err := ResolveCollision(ray, shape, space.useTags)

				if err != nil {
					return nil, Zero(), err
				}

				if raycastHit && !shape.ContainsPoint(ray.p) {
					contacts, err := Contact(ray, shape)

					if err != nil {
						return nil, Zero(), err
					}

					for _, contact := range contacts {
						sqrDistance := SquaredDistance(ray.p, contact)

						if !minExists {
							hit = contact
							hitShape = shape
							minSquaredDistance = sqrDistance

							minExists = true
						} else if sqrDistance < minSquaredDistance {
							hit = contact
							hitShape = shape
							minSquaredDistance = sqrDistance
						}
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

	return hitShape, hit, nil
}

// Boxcast casts a box in the space and returns all the
// shapes overlapped by this box.
func (space *Space) Boxcast(rect *Rectangle) (Shapes, error) {
	if rect == nil {
		return nil, fmt.Errorf("the rectangle is nil")
	}

	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	shapes := make(Shapes, 0)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)
		overlapped, err := node.boundary.collidesRectangle(rect)

		if err != nil {
			return nil, err
		}

		if !overlapped {
			continue
		}

		if node.northWest == nil {
			for shape := range node.shapes {
				boxcastHit, err := ResolveCollision(rect, shape, space.useTags)

				if err != nil {
					return nil, err
				}

				if boxcastHit {
					shapes.Insert(shape)
				}
			}
		}
	}

	return shapes, nil
}

// Circlecast casts a circle in the space and returns all the
// shapes overlapped by the circle.
func (space *Space) Circlecast(circle *Circle) (Shapes, error) {
	if circle == nil {
		return nil, fmt.Errorf("the circle is nil")
	}

	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)
	shapes := make(Shapes, 0)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)
		overlapped, err := node.boundary.collidesCircle(circle)

		if err != nil {
			return nil, err
		}

		if !overlapped {
			continue
		}

		if node.northWest == nil {
			for shape := range node.shapes {
				circlecastHit, err := ResolveCollision(circle, shape, space.useTags)

				if err != nil {
					return nil, err
				}

				if circlecastHit {
					shapes.Insert(shape)
				}
			}
		}
	}

	return shapes, nil
}
