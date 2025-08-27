package world

import (
	"fmt"
	"image/color"
	"sand2/events"
	"sand2/observer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

func (c *Chunk) index(localX, localY int) int {
	return localY*CHUNK_SIZE + localX
}

func (c *Chunk) Set(localX, localY int, col color.RGBA) {
	index := c.index(localX, localY)
	c.buff[index] = col
	c.next[index] = col
}

func (c *Chunk) SetP(localPos Position, col color.RGBA) {
	c.Set(localPos.X, localPos.Y, col)
}

func (c *Chunk) Get(localX, localY int) color.RGBA {
	return c.buff[c.index(localX, localY)]
}

func (c *Chunk) GetP(localPos Position) color.RGBA {
	return c.Get(localPos.X, localPos.Y)
}

func (c *Chunk) Swap() {
	copy(c.next, c.buff)
	c.buff, c.next = c.next, c.buff
}

func (c *Chunk) Render(output *rl.Texture2D) {
	x := c.position.X*CHUNK_SIZE
	y := c.position.Y*CHUNK_SIZE
	
	rec := rl.NewRectangle(float32(x), float32(y), float32(CHUNK_SIZE), float32(CHUNK_SIZE))
	rl.DrawRectangleLines(int32(x), int32(y), int32(CHUNK_SIZE), int32(CHUNK_SIZE), rl.Green)
	rl.UpdateTextureRec(*output, rec, c.buff)
}

func (c *Chunk) Update(updateStep func(x, y int)) {
	for x := c.dirtyMinX; x < c.dirtyMaxX; x++ {
		for y := c.dirtyMinY; y < c.dirtyMaxY; y++ {
			c.Nofify(events.EVENT_CHUNK_UPDATE, x, y)
			updateStep(x, y)
		}
	}
}

func (c *Chunk) ToLocalSpace(globalX, globalY int) Position {
	x := globalX - c.position.X*CHUNK_SIZE
	y := globalY - c.position.Y*CHUNK_SIZE
	if x < 0 || x >= CHUNK_SIZE || y < 0 || y >= CHUNK_SIZE {
		panic(fmt.Sprintf("Inccorect conversion to local space! global: (%d, %d), local: (%d, %d), chunk: (%d, %d)", globalX, globalY, x, y, c.position.X, c.position.Y))
	}

	return Position{
		X: x,
		Y: y,
	}
}

func (c *Chunk) GetPosition() Position {
	return c.position
}
