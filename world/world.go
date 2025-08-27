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

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func ToPosition(vec2 rl.Vector2) Position {
	return NewPosition(int(vec2.X), int(vec2.Y))
}

func ToVector2(pos Position) rl.Vector2 {
	return rl.NewVector2(float32(pos.X), float32(pos.Y))
}

type World struct {
	screenBuff rl.RenderTexture2D
	background rl.Texture2D
	chunks     map[Position]*Chunk
	Camera     rl.Camera2D

	observer.Observerable
}

func floor(a, b int) int {
	if b == 0 {
		panic("Division by 0")
	}

	if a >= 0 {
		return a / b
	}

	return ((a + 1) / b) - 1
}

func New() *World {
	w := &World{
		chunks: make(map[Position]*Chunk),
		Camera: rl.NewCamera2D(rl.NewVector2(100, 100), rl.NewVector2(0, 0), 0, 1),
	}

	// Screen buffer to draw stuff
	screenW := rl.GetScreenWidth()
	screenH := rl.GetScreenHeight()

	w.screenBuff = rl.LoadRenderTexture(int32(screenW), int32(screenH))

	// Background stuff
	bgImg := rl.GenImageChecked(32, 32, 16, 16, rl.Blank, rl.DarkGray)
	w.background = rl.LoadTextureFromImage(bgImg)
	rl.UnloadImage(bgImg)
	rl.SetTextureWrap(w.background, rl.WrapRepeat)

	return w
}

func (w *World) Clean() {
	rl.UnloadTexture(w.background)
	rl.UnloadRenderTexture(w.screenBuff)
}

func (w *World) MouseToWorldSpace() rl.Vector2 {
	mpos := rl.GetMousePosition()
	return rl.GetScreenToWorld2D(mpos, w.Camera)
}

func (w *World) GetChunkGridPos(pos Position) Position {
	return NewPosition(floor(pos.X, CHUNK_SIZE), floor(pos.Y, CHUNK_SIZE))
}

// Set particles at global position
func (w *World) Set(globalPos Position, mat Material) {
	pos := w.GetChunkGridPos(globalPos)

	chunk, ok := w.chunks[pos]

	if !ok {
		chunk = NewChunk(pos)
		w.chunks[pos] = chunk
	}

	localPos := chunk.ToLocalSpace(globalPos)

	chunk.Set(localPos, mat.color)
}

func (w *World) Update() {
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		w.Camera.Offset = rl.Vector2Add(w.Camera.Offset, rl.GetMouseDelta())
	}

	// if wheel := rl.GetMouseWheelMove(); wheel != 0 {
	// 	mpos := rl.GetMousePosition()
	// 	globalPos := rl.GetScreenToWorld2D(mpos, w.Camera)
	//
	// 	if wheel < 0 {
	//
	// 		w.Camera.Zoom = max(w.Camera.Zoom-0.2, 0.2)
	// 	} else {
	// 		w.Camera.Zoom = min(w.Camera.Zoom+0.2, 4)
	// 	}
	// 	w.Camera.Offset = mpos
	// 	w.Camera.Target = globalPos
	// }

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		mpos := w.MouseToWorldSpace()

		w.Set(ToPosition(mpos), Material{color: rl.Yellow})
	}
}

func (w *World) getVisibleChunks() []*Chunk {
	result := make([]*Chunk, 0, 1)

	screenW := rl.GetScreenWidth()
	screenH := rl.GetScreenHeight()

	topLeft := rl.GetScreenToWorld2D(rl.Vector2Zero(), w.Camera)
	bottomRight := rl.GetScreenToWorld2D(rl.NewVector2(float32(screenW+CHUNK_SIZE*4), float32(screenH-CHUNK_SIZE)), w.Camera)

	chukMinPos := w.GetChunkGridPos(ToPosition(topLeft))
	chukMaxPos := w.GetChunkGridPos(ToPosition(bottomRight))

	for y := chukMaxPos.Y; y >= chukMinPos.Y; y-- {
		for x := chukMinPos.X; x < chukMaxPos.X; x++ {
			if chunk, ok := w.chunks[NewPosition(x, y)]; ok {
				result = append(result, chunk)
			}
		}
	}

	return result
}

func (w *World) Render() {
	// screenW := rl.GetScreenWidth()
	// screenH := rl.GetScreenHeight()

	// chunks := w.getVisibleChunks()

	rl.BeginTextureMode(w.screenBuff)
	rl.ClearBackground(rl.NewColor(0, 0, 0, 1))
	rl.EndTextureMode()
	rl.BeginMode2D(w.Camera)
	// for _, chunk := range chunks {
	for _, chunk := range w.chunks {
		gridPos := chunk.GetWorldPosition()
		pos := rl.GetWorldToScreen2D(ToVector2(gridPos), w.Camera)
		rec := rl.NewRectangle(float32(pos.X), float32(pos.Y), float32(CHUNK_SIZE), float32(CHUNK_SIZE))

		// text := fmt.Sprintf("%v", rec)
		// rl.DrawText(text, int32(gridPos.X), int32(gridPos.Y), 10, rl.Gray)

		chunk.Render(&w.screenBuff.Texture, rec)
	}
	// rl.DrawTextureRec(w.background, rl.NewRectangle(-w.Camera.Offset.X, w.Camera.Offset.Y, float32(screenW), float32(screenH)), rl.Vector2Zero(), rl.RayWhite)

	rl.DrawTexture(w.screenBuff.Texture, 0, 0, rl.White)
	rl.EndMode2D()
	rl.DrawText(fmt.Sprintf("Chunks to draw: %d", len(w.chunks)), 5, 20, 20, rl.Blue)
}
