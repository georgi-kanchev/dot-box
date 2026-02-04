package dots

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nsf/termbox-go"
)

var eventQueue chan termbox.Event

var screen rl.RenderTexture2D

var brailleMap = [4][2]rune{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

func Main() {
	initTermbox()
	initRaylib()

	// var texture = rl.LoadTexture("picture.png")
	// var a float32

	dist := float32(10.0)

	camera := rl.Camera3D{
		// To get the -45 Z rotation, X and Z must be equal
		Position:   rl.NewVector3(dist, 8.165, dist),
		Target:     rl.NewVector3(0.0, 0.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       10.0, // Adjust this to make it "bigger"
		Projection: rl.CameraOrthographic,
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.BeginTextureMode(screen)
		rl.ClearBackground(rl.NewColor(0, 0, 0, 0))

		rl.BeginMode3D(camera)

		rl.DrawCube(rl.NewVector3(0.5, 0.8165/2, 0.5), 1, 0.8165, 1, rl.White)
		rl.DrawGrid(10, 1)
		rl.EndMode3D()
		// var w, h = float32(texture.Width), float32(texture.Height)
		// rl.DrawTextureRec(texture, rl.Rectangle{Width: w, Height: h}, rl.Vector2{}, rl.White)
		// rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 0, 140, 32, rl.White)
		// a += rl.GetFrameTime() * 10
		// rl.DrawRectanglePro(rl.Rectangle{X: 130, Y: 60, Width: 50, Height: 50}, rl.Vector2{X: 25, Y: 25}, a, rl.Red)

		rl.EndTextureMode()

		var w, h = float32(screen.Texture.Width), float32(screen.Texture.Height)
		rl.SetWindowSize(int(w)*8, int(h)*8)
		rl.ClearBackground(rl.Black)
		rl.DrawTexturePro(screen.Texture, rl.Rectangle{X: 0, Y: 0, Width: w, Height: -h}, rl.Rectangle{X: 0, Y: 0, Width: w * 8, Height: -h * 8}, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()

		updateTerminal()
		processEvents()
	}
}

func updateTerminal() {
	var img = rl.LoadImageFromTexture(screen.Texture)
	var pixels = rl.LoadImageColors(img)
	var w, h = int(img.Width), int(img.Height)
	var termW, termH = termbox.Size()
	defer rl.UnloadImage(img)
	defer rl.UnloadImageColors(pixels)

	for y := range termH {
		for x := range termW {
			var mask rune = 0
			for dy := range 4 {
				for dx := range 2 {
					var px, py = x*2 + dx, (h - 1) - (y*4 + dy)
					if px < w && py >= 0 && py < h {
						if pixels[py*w+px].A > 0 {
							mask |= brailleMap[dy][dx]
						}
					}
				}
			}

			termbox.SetChar(x, y, 0x2800+mask)
			termbox.SetFg(x, y, termbox.RGBToAttribute(255, 255, 255))
		}
	}
	termbox.Flush()
}

func initTermbox() {
	termbox.Init()
	termbox.SetOutputMode(termbox.OutputRGB)
	// termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetInputMode(termbox.InputAlt)

	eventQueue = make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
}
func initRaylib() {
	rl.SetTraceLogLevel(rl.LogNone)
	// rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(400, 400, "")
	rl.SetTargetFPS(60)
	rl.SetBlendFactors(rl.OneMinusDstColor, rl.Zero, rl.FuncAdd)
	rl.BeginBlendMode(rl.BlendCustom)

	var w, h = termbox.Size()
	screen = rl.LoadRenderTexture(int32(w*2), int32(h*4))
	rl.SetTextureFilter(screen.Texture, rl.TextureFilterNearest)
}

func processEvents() {
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
				rl.CloseWindow()
			}
			if ev.Type == termbox.EventResize {
				rl.UnloadRenderTexture(screen)
				screen = rl.LoadRenderTexture(int32(ev.Width*2), int32(ev.Height*4))
				rl.SetTextureFilter(screen.Texture, rl.TextureFilterNearest)
			}
			if ev.Type == termbox.EventMouse {
				termbox.SetCursor(ev.MouseX, ev.MouseY)
			}
		default:
			return
		}
	}
}
