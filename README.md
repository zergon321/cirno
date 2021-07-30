# Cirno [![GoDoc](https://godoc.org/github.com/zergon321/cirno?status.svg)](https://pkg.go.dev/github.com/zergon321/cirno)

An easy to use collision detection and resolution library written in **Go** programming language. The library is still in development and the API might change.

## Installation

No **C** dependencies required to compile the library. Just use the following command to install it.

```bash
go get github.com/zergon321/cirno
```

## Tutorial

The tutorial series explains how to use different API methods and data types of **Cirno** to handle everything related to collisions in games.

- [Tutorial 01: Basic collisions](https://github.com/zergon321/cirno/wiki/Tutorial-01:-Basic-collisions)

## Examples

All the example programs using the library are located  in [examples](https://github.com/zergon321/cirno/tree/master/examples) directory. [Pixel](https://github.com/faiface/pixel) is required to run any of them. The most important demos are described below.

- [benchmark](https://github.com/zergon321/cirno/tree/master/examples/benchmark) - a small demo application that creates 1000 rectangles, randomly moves and rotates them and serches for collisions between them. It clearly showcases quad tree in action;

- [contacts](https://github.com/zergon321/cirno/tree/master/examples/contacts) - a small demo application that showcases finding contact points between shape outlines;

- [raycast](https://github.com/zergon321/cirno/tree/master/examples/raycast) - a small demo application that showcases raycast, just like [Physics2D.Raycast](https://docs.unity3d.com/ScriptReference/Physics2D.Raycast.html) from **Unity**;

- [sliding](https://github.com/zergon321/cirno/tree/master/examples/sliding) - a small demo that showcases shapes movement with sliding collision.

## Games

The most noteable game created with **Cirno** at the current time is a **Touhou** style [danmaku demo](https://zergon321.itch.io/touhou-game-in-go). The less noteable one is a tiny platrformer level in [examples](https://github.com/zergon321/cirno/tree/master/examples) directory.

| [Danmaku](https://zergon321.itch.io/touhou-game-in-go) | [Platformer](https://github.com/zergon321/cirno/blob/master/examples/platformer) |
| --- | --- |
| ![Danmaku](https://github.com/zergon321/cirno/blob/master/screenshots/danmaku.png) | ![Platformer](https://github.com/zergon321/cirno/blob/master/screenshots/platformer.png) |

If you have a game using **Cirno** and you want it to be present in this list, just create a new PR with your changes to the README.

## Features

- Shapes to attach to game objects to detect and resolve collisions between them:
  - circle;
  - line segment (or just line);
  - rectangle (OBB, oriented bounding box).
- Quadtree space index
- Raycast
- Contacts finding methods
- Normal computing methods
- Movement and rotation approximation
- Tag system

## Contributing

For minor and unimportant errors such as typos please just create issues instead of PRs fixing them.

Code changes adding features, optimizations, tests and bugfixes are welcome. All the contributions in [examples](https://github.com/zergon321/cirno/tree/master/examples) directory should be really small, preferrably one-file. If your game is quite big, consider placing links to it in the list of games in the README.
