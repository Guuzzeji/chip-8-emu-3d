package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Create chip 8 cimputer
	var chip Cpu
	chip.Init()
	chip.LoadRom("./rom/Pong.ch8")
	chip.Display.Scale = 10
	chip.Speed = 10

	// Setting raylib window
	rl.InitWindow(640, 480, "Golang 3D Chip 8 emu")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// 3D stuff
	var camera rl.Camera3D
	camera.Position = rl.Vector3{X: 10.0, Y: 25.0, Z: 80.0}
	camera.Target = rl.Vector3{X: 10.0, Y: 10.0, Z: 0.0}
	camera.Up = rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0}
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	// Update loop
	isFreeCamera := false

	for !rl.WindowShouldClose() {

		// Raylib update and begin drawing
		if rl.IsKeyPressed(96) {
			isFreeCamera = !isFreeCamera
		}

		if isFreeCamera {
			rl.UpdateCamera(&camera, rl.CameraFree)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.BeginMode3D(camera)

		//fmt.Println("Keypressed = ", rl.GetKeyPressed())
		//fmt.Println(rl.GetCharPressed())

		chip.Cycle() // Update Chip 8 cpu

		// Render each pixel from chip 8 computer
		for y := 0; y < chip.Display.Height; y++ {
			for x := 0; x < chip.Display.Width; x++ {
				// offsetX := (x * 1) - (1 / 2)
				// offsetY := (y * 1) + (1 / 2)
				offsetX := (x * 1)
				offsetY := (y * 1)

				// fmt.Println(x, y)
				positon := rl.Vector3{X: float32(offsetX) - 25, Y: -float32(offsetY) + 25, Z: 0.0}

				if chip.Display.Canvas[y][x] > 0 {
					rl.DrawCube(positon, 1.0, 1.0, 1.0, rl.Black)
					rl.DrawCubeWires(positon, 1.0, 1.0, 1.0, rl.White)
				}
			}
		}

		rl.DrawGrid(10.0, 10.0)

		rl.EndMode3D()
		rl.EndDrawing()

		// Text debug
		//chip.Display.Render()
	}

}
