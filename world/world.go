package world

import (
	"fmt"
	"sand2/observer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const CHUNK_SIZE int = 64

// TODO: Fix resolution
const SCREEN_W = 1280
const SCREEN_H = 720

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

type World struct {
	screenBuff rl.RenderTexture2D
	background rl.Texture2D
	chunks map[Position]*Chunk
	Camera rl.Camera2D

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

func mod(a, b int) int {
	m := a % b
	
	if m < 0 {
		m += b
	}
	
	return m
}

func New() *World {
	w := &World{
		chunks: make(map[Position]*Chunk),
		Camera: rl.NewCamera2D(rl.NewVector2(100, 100), rl.NewVector2(0, 0), 0, 1),
	}
	
	// Screen buffer to draw stuff
	w.screenBuff = rl.LoadRenderTexture(int32(SCREEN_W), int32(SCREEN_W))
	
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

// Set particles at global position
func (w *World) Set(globalX, globalY int, mat Material) {
	pos := NewPosition(floor(globalX, CHUNK_SIZE), floor(globalY, CHUNK_SIZE))

	chunk, ok := w.chunks[pos]

	if !ok {
		chunk = NewChunk(pos)
		w.chunks[pos] = chunk
	}

	localX := mod(globalX, CHUNK_SIZE)
	localY := mod(globalY, CHUNK_SIZE)

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

func (w *World) getVisibleChunks() []*Chunk {
	result := make([]*Chunk, 0, 1)
	
	topLeft := rl.GetScreenToWorld2D(rl.Vector2Zero(), w.Camera)
	bottomRight := rl.GetScreenToWorld2D(rl.NewVector2(float32(SCREEN_W + CHUNK_SIZE), float32(SCREEN_H)), w.Camera)
	
	chunkMinX := floor(int(topLeft.X), CHUNK_SIZE)
	chunkMinY := floor(int(topLeft.Y), CHUNK_SIZE)
	chunkMaxX := floor(int(bottomRight.X), CHUNK_SIZE)
	chunkMaxY := floor(int(bottomRight.Y), CHUNK_SIZE)
	
	for y := chunkMaxY; y >= chunkMinY; y-- {
		for x := chunkMinX; x < chunkMaxX; x++ {
			if chunk, ok := w.chunks[NewPosition(x, y)]; ok {
				result = append(result, chunk)
			}
		}
	}
	
	return result
}

func (w *World) Render() {
	rl.BeginMode2D(w.Camera)
	rl.DrawTextureRec(w.background, rl.NewRectangle(0, 0, float32(SCREEN_W), float32(SCREEN_H)), rl.Vector2Zero(), rl.RayWhite)
	
	chunks := w.getVisibleChunks()
	
	// UpdateTextureRec(texture Texture2D, rec Rectangle, pixels []color.RGBA)
	
	for _, chunk := range chunks {
		chunk.Render(&w.screenBuff.Texture)
	}
	
	rl.DrawTexture(w.screenBuff.Texture, 0, 0, rl.White)
	
	rl.EndMode2D()
	
	rl.DrawText(fmt.Sprintf("Chunks to draw: %d", len(chunks)), 5, 20, 20, rl.Blue)
}
