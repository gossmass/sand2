package world

import (
	"image/color"
	"sand2/events"
	"sand2/observer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func mod(a, b int) int {
	m := a % b

	if m < 0 {
		m += b
	}

	return m
}

type Chunk struct {
	position  Position // position in the world (grid coordinats e.g. (0,0) (0,1) etc.)
	active    bool
	partciles int // amount of filled cells
	buff      []color.RGBA
	next      []color.RGBA
	dirtyMinX int
	dirtyMaxX int
	dirtyMinY int
	dirtyMaxY int

	observer.Observerable
}

func NewChunk(pos Position) *Chunk {
	return &Chunk{
		position:  pos,
		active:    false,
		buff:      make([]color.RGBA, CHUNK_SIZE*CHUNK_SIZE),
		next:      make([]color.RGBA, CHUNK_SIZE*CHUNK_SIZE),
		dirtyMaxX: CHUNK_SIZE,
		dirtyMaxY: CHUNK_SIZE,
	}
}

func (c *Chunk) index(localPos Position) int {
	return localPos.Y*CHUNK_SIZE + localPos.X
}

func (c *Chunk) Set(localPos Position, col color.RGBA) {
	index := c.index(localPos)
	c.buff[index] = col
	c.next[index] = col
	c.active = true
}

func (c *Chunk) Get(localPos Position) color.RGBA {
	return c.buff[c.index(localPos)]
}

func (c *Chunk) Swap() {
	copy(c.next, c.buff)
	c.buff, c.next = c.next, c.buff
}

func (c *Chunk) Render2() {
	// rl.DrawTexture()
}

func (c *Chunk) Render(output *rl.Texture2D, rec rl.Rectangle) {
	if !c.active {
		return
	}

	// pos := c.GetWorldPosition()
	// rec := rl.NewRectangle(float32(pos.X), float32(pos.Y), float32(CHUNK_SIZE), float32(CHUNK_SIZE))
	// rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(CHUNK_SIZE), int32(CHUNK_SIZE), rl.Green)
	rl.UpdateTextureRec(*output, rec, c.buff)
}

func (c *Chunk) Update(updateStep func(x, y int)) {
	for x := c.dirtyMinX; x < c.dirtyMaxX; x++ {
		for y := c.dirtyMinY; y < c.dirtyMaxY; y++ {
			c.Nofify(events.EVENT_CHUNK_UPDATE, x, y)
			updateStep(x, y)
		}
	}

	c.active = false
}

func (c *Chunk) GetWorldPosition() Position {
	return NewPosition(c.position.X*CHUNK_SIZE, c.position.Y*CHUNK_SIZE)
}

func (c *Chunk) ToLocalSpace(globalPos Position) Position {
	return NewPosition(mod(globalPos.X, CHUNK_SIZE), mod(globalPos.Y, CHUNK_SIZE))
}

func (c *Chunk) GetPosition() Position {
	return c.position
}
