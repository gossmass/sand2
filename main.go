package main

import (
	"sand2/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "Sand V2")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	worldDebugger := world.NewDebugger()

	world := world.New()
	world.AddObserver(worldDebugger)

	// Background stuff
	bgImg := rl.GenImageChecked(32, 32, 16, 16, rl.Blank, rl.DarkGray)
	bgTex := rl.LoadTextureFromImage(bgImg)
	rl.SetTextureWrap(bgTex, rl.WrapRepeat)
	rl.UnloadImage(bgImg)
	bgRect := rl.NewRectangle(0, 0, 1280, 720)
	bgPos := rl.NewVector2(0, 0)
	defer rl.UnloadTexture(bgTex)

	// Main loop
	for !rl.WindowShouldClose() {
		world.Update()

		rl.BeginDrawing()
		rl.BeginMode2D(world.Camera)
		rl.ClearBackground(rl.Black)
		rl.DrawTextureRec(bgTex, bgRect, bgPos, rl.RayWhite)
		world.Render()
		rl.EndMode2D()
		rl.EndDrawing()
	}
}
