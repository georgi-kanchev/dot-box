package dots

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nsf/termbox-go"
)

var x, y = 0, 0

func Main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(100, 100, "")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{})
		rl.EndDrawing()
		processEvents(eventQueue)

		termbox.Clear(0, 0)
		termbox.SetChar(0, 0, 'H')
		termbox.SetChar(1, 0, 'e')
		termbox.SetChar(2, 0, 'l')
		termbox.SetChar(3, 0, 'l')
		termbox.SetChar(4, 0, 'o')
		termbox.SetChar(x, y, '#')
		termbox.Flush()
	}
	termbox.Close()
}

func processEvents(eventQueue chan termbox.Event) {
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
				rl.CloseWindow()
			}
			if ev.Type == termbox.EventResize {
			}
			if ev.Type == termbox.EventMouse {
				termbox.SetCursor(ev.MouseX, ev.MouseY)
			}
		default:
			return
		}
	}
}
