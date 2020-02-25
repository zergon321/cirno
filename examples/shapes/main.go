package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"reflect"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tracer8086/cirno"
	"golang.org/x/image/colornames"
)

const (
	width     = 1280
	height    = 720
	speed     = 700
	intensity = 100
)

var (
	controlledShape string
)

type object struct {
	shape     cirno.Shape
	sprite    *pixel.Sprite
	transform pixel.Matrix
}

func loadPicture(pic string) (pixel.Picture, error) {
	file, err := os.Open(pic)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func parseFlags() {
	flag.StringVar(&controlledShape, "shape", "circle",
		"The shape controlled during execution of the demo.")

	flag.Parse()
}

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno test",
		Bounds: pixel.R(0, 0, width, height),
		//VSync:  true,
		//Undecorated: true,
		//Monitor: pixelgl.PrimaryMonitor(),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	// Setup physics.
	space, err := cirno.NewSpace(5, 20, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height))
	handleError(err)

	// Create borders.
	borderWest := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(0, height))
	borderSouth := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(width, 0))
	borderNorth := cirno.NewLine(cirno.NewVector(0, height), cirno.NewVector(width, height))
	borderEast := cirno.NewLine(cirno.NewVector(width, 0), cirno.NewVector(width, height))

	// Add the borders in the space.
	err = space.Add(borderNorth)
	handleError(err)
	err = space.Add(borderWest)
	handleError(err)
	err = space.Add(borderSouth)
	handleError(err)
	err = space.Add(borderEast)
	handleError(err)

	// Setup graphics sprites.
	circlePic, err := loadPicture("round_particle.png")
	handleError(err)
	cubePic, err := loadPicture("cube.png")
	handleError(err)
	rectPic, err := loadPicture("rect.png")
	handleError(err)

	// Setup objects.
	circle := &object{
		shape:  cirno.NewCircle(cirno.NewVector(1024, 256), 30),
		sprite: pixel.NewSprite(circlePic, pixel.R(0, 0, 45, 45)),
		transform: pixel.IM.Scaled(pixel.ZV, 60.0/45.0).
			Moved(pixel.V(1024, 256)),
	}
	otherCircle := &object{
		shape:  cirno.NewCircle(cirno.NewVector(420, 380), 50),
		sprite: pixel.NewSprite(circlePic, pixel.R(0, 0, 45, 45)),
		transform: pixel.IM.Scaled(pixel.ZV, 100.0/45.0).
			Moved(pixel.V(420, 380)),
	}
	cube := &object{
		shape:  cirno.NewRectangle(cirno.NewVector(128, 256), 100, 100, 60),
		sprite: pixel.NewSprite(cubePic, pixel.R(0, 0, 32, 32)),
		transform: pixel.IM.Scaled(pixel.ZV, 100.0/32.0).
			Rotated(pixel.ZV, 60*cirno.DegToRad).
			Moved(pixel.V(128, 256)),
	}
	otherCube := &object{
		shape:  cirno.NewRectangle(cirno.NewVector(1024, 512), 100, 100, 0),
		sprite: pixel.NewSprite(cubePic, pixel.R(0, 0, 32, 32)),
		transform: pixel.IM.Scaled(pixel.ZV, 100.0/32.0).
			Moved(pixel.V(1024, 512)),
	}
	rect := &object{
		shape:  cirno.NewRectangle(cirno.NewVector(640, 520), 150, 50, 30.0),
		sprite: pixel.NewSprite(rectPic, pixel.R(0, 0, 32, 44)),
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(150.0/32.0, 50.0/44.0)).
			Rotated(pixel.ZV, 30*cirno.DegToRad).
			Moved(pixel.V(640, 520)),
	}
	line := &object{
		shape: cirno.NewLine(cirno.NewVector(480, 520), cirno.NewVector(480, 370)),
	}
	otherLine := &object{
		shape: cirno.NewLine(cirno.NewVector(720, 450), cirno.NewVector(720, 280)),
	}

	line.shape.(*cirno.Line).Rotate(-90)
	//otherLine.shape.(*cirno.Line).Rotate(90)
	//line.shape.(*cirno.Line).Rotate(60)
	//otherLine.shape.(*cirno.Line).Rotate(-30)

	err = space.Add(circle.shape)
	handleError(err)
	err = space.Add(cube.shape)
	handleError(err)
	err = space.Add(line.shape)
	handleError(err)
	err = space.Add(rect.shape)
	handleError(err)
	err = space.Add(otherCircle.shape)
	handleError(err)
	err = space.Add(otherCube.shape)
	handleError(err)
	err = space.Add(otherLine.shape)
	handleError(err)

	// Determine which object should
	// be controlled.
	var obj *object

	switch controlledShape {
	case "circle":
		obj = circle

	case "rectangle":
		obj = rect

	case "cube":
		obj = cube

	case "line":
		obj = line

	default:
		panic(fmt.Errorf("Undefined shape"))
	}

	// Setup metrics.
	last := time.Now()
	fps := 0
	cps := 0
	perSecond := time.Tick(time.Second)

	// IMDraw to draw lines.
	imd := imdraw.New(nil)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.White)

		// Motion control.
		movement := cirno.Zero

		if win.Pressed(pixelgl.KeyUp) {
			movement = movement.Add(cirno.Up)
		}

		if win.Pressed(pixelgl.KeyDown) {
			movement = movement.Add(cirno.Down)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			movement = movement.Add(cirno.Left)
		}

		if win.Pressed(pixelgl.KeyRight) {
			movement = movement.Add(cirno.Right)
		}

		// Move the controllable shape.
		if movement != cirno.Zero {
			movement = movement.MultiplyByScalar(speed * deltaTime)
			shapes, err := space.WouldBeColliding(obj.shape,
				movement.MultiplyByScalar(1.5))
			handleError(err)

			// If a collision occurres, the shape
			// won't move.
			pos := obj.shape.Center()

			if len(shapes) > 0 {
				cps += len(shapes)
				pos, err = cirno.Approximate(obj.shape, movement,
					shapes, intensity)
				handleError(err)
				movement = pos.Subtract(obj.shape.Center())
			}

			// Consistent code block.
			obj.shape.Move(movement)
			space.AdjustShapePosition(obj.shape)
			obj.transform = obj.transform.
				Moved(pixel.V(movement.X, movement.Y))
			_, err = space.Update(obj.shape)
			handleError(err)

			/*shapes, err := space.CollidingWith(obj.shape)
			handleError(err)

			// If a collision occurres, approximate
			// the shape position.
			if len(shapes) > 0 {
				obj.shape.Move(movement.MultiplyByScalar(-1))
				obj.transform = obj.transform.
					Moved(pixel.V(-movement.X, -movement.Y))

				pos, err := cirno.Approximate(obj.shape, movement,
					shapes, intensity)
				handleError(err)
				movement = pos.Subtract(obj.shape.Center())
				obj.shape.Move(movement)
				space.AdjustShapePosition(obj.shape)
				obj.transform = obj.transform.
					Moved(pixel.V(movement.X, movement.Y))
				_, err = space.Update(obj.shape)
				handleError(err)
			}*/

			for shape := range shapes {
				t := reflect.TypeOf(shape).Elem()
				fmt.Println(t.Name(), shape.Center())
			}

			fmt.Println("Movement:", movement)
			fmt.Println("Position:", obj.shape.Center())
		}

		// Rendering.
		cube.sprite.Draw(win, cube.transform)
		circle.sprite.Draw(win, circle.transform)
		rect.sprite.Draw(win, rect.transform)
		otherCircle.sprite.Draw(win, otherCircle.transform)
		otherCube.sprite.Draw(win, otherCube.transform)

		lineShape := line.shape.(*cirno.Line)
		imd.Clear()
		imd.Color = colornames.Blue
		imd.Push(pixel.V(lineShape.P().X, lineShape.P().Y))
		imd.Color = colornames.Red
		imd.Push(pixel.V(lineShape.Q().X, lineShape.Q().Y))
		imd.Line(1)
		imd.Draw(win)

		otherLineShape := otherLine.shape.(*cirno.Line)
		imd.Clear()
		imd.Color = colornames.Blue
		imd.Push(pixel.V(otherLineShape.P().X, otherLineShape.P().Y))
		imd.Color = colornames.Red
		imd.Push(pixel.V(otherLineShape.Q().X, otherLineShape.Q().Y))
		imd.Line(1)
		imd.Draw(win)

		win.Update()

		// Show FPS in the window title.
		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d, CPS: %d", cfg.Title, fps, cps))
			fps = 0
			cps = 0

		default:
		}
	}
}

func main() {
	parseFlags()
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
