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

	// Main loop
	for !rl.WindowShouldClose() {
		world.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		world.Render()
		
		rl.DrawFPS(5, 5)
		rl.EndDrawing()
	}
	
	world.Clean()
}
