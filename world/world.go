package world

import (
	"fmt"
	"sand2/observer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const CHUNK_SIZE int = 64

type Position struct {
	X int
	Y int
}

type World struct {
	chunks map[Position]*Chunk
	Camera rl.Camera2D

	observer.Observerable
}

func New() *World {
	return &World{
		chunks: make(map[Position]*Chunk),
		Camera: rl.NewCamera2D(rl.NewVector2(100, 100), rl.NewVector2(0, 0), 0, 1),
	}
}

func (w *World) MouseToWorldSpace() rl.Vector2 {
	mpos := rl.GetMousePosition()
	return rl.GetScreenToWorld2D(mpos, w.Camera)
}

// Set particles at global position
func (w *World) Set(globalX, globalY int, mat Material) {
	oX := globalX
	oY := globalY

	if globalX < 0 {
		globalX -= CHUNK_SIZE
	}

	if globalY < 0 {
		globalY -= CHUNK_SIZE
	}

	pos := Position{
		X: globalX / CHUNK_SIZE,
		Y: globalY / CHUNK_SIZE,
	}

	chunk, ok := w.chunks[pos]

	if !ok {
		chunk = NewChunk(pos)
		w.chunks[pos] = chunk
	}

	// p := chunk.ToLocalSpace(globalX, globalY)
	localX := oX - (pos.X * CHUNK_SIZE)
	localY := oY - (pos.Y * CHUNK_SIZE)

	chunk.Set(localX, localY, mat.color)
}

func (w *World) Update() {
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		w.Camera.Offset = rl.Vector2Add(w.Camera.Offset, rl.GetMouseDelta())
	}

	if wheel := rl.GetMouseWheelMove(); wheel != 0 {
		mpos := rl.GetMousePosition()
		globalPos := rl.GetScreenToWorld2D(mpos, w.Camera)

		if wheel < 0 {

			w.Camera.Zoom = max(w.Camera.Zoom-0.2, 0.2)
		} else {
			w.Camera.Zoom = min(w.Camera.Zoom+0.2, 4)
		}
		w.Camera.Offset = mpos
		w.Camera.Target = globalPos
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		mpos := w.MouseToWorldSpace()

		w.Set(int(mpos.X), int(mpos.Y), Material{color: rl.Yellow})
	}
}

func (w *World) Render() {
	// Cameraman

	for pos, _ := range w.chunks {
		x := int32(pos.X * CHUNK_SIZE)
		y := int32(pos.Y * CHUNK_SIZE)
		text := fmt.Sprintf("%v", pos)

		rl.DrawRectangleLines(x, y, int32(CHUNK_SIZE), int32(CHUNK_SIZE), rl.Green)
		rl.DrawText(text, x, y, 10, rl.RayWhite)
	}

	// mpos := w.MouseToWorldSpace()
	//
	// hx := int(mpos.X)
	// hy := int(mpos.Y)
	// ox := hx
	//
	// if hx < 0 {
	// 	hx -= CHUNK_SIZE
	// }
	//
	// if hy < 0 {
	// 	hy -= CHUNK_SIZE
	// }

	// lx := ox - ((hx / CHUNK_SIZE) * CHUNK_SIZE)
	// ly := 0
	//
	// rl.DrawText(fmt.Sprintf("%v, %v\nDIV: %v, %v\nPERC: %v, %v", hx, hy, hx/CHUNK_SIZE, hy/CHUNK_SIZE, lx, ly), 10, 10, 24, rl.Yellow)
}
