package dots

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nsf/termbox-go"
)

var brailleMap = [4][2]int{{0x1, 0x8}, {0x2, 0x10}, {0x4, 0x20}, {0x40, 0x80}}
var termW, termH = 0, 0
var zBuffer []float64

func Main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

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

	resize(termbox.Size())
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()

		ev := <-eventQueue
		if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
			return
		}
		if ev.Type == termbox.EventResize {
			resize(ev.Width, ev.Height)
		}
		termbox.Clear(0, 0)
		termbox.SetChar(0, 0, 'H')
		termbox.Flush()
	}

	// err := termbox.Init()
	// if err != nil {
	// 	panic(err)
	// }
	// defer termbox.Close()

	// eventQueue := make(chan termbox.Event)
	// go func() {
	// 	for {
	// 		eventQueue <- termbox.PollEvent()
	// 	}
	// }()

	// angle := 0.0
	// ticker := time.NewTicker(time.Second / 30)
	// defer ticker.Stop()

	// resize(termbox.Size())
	// for {
	// 	select {
	// 	case ev := <-eventQueue:
	// 		if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
	// 			return
	// 		}
	// 		if ev.Type == termbox.EventResize {
	// 			resize(ev.Width, ev.Height)
	// 		}
	// 	case <-ticker.C:
	// 		for i := range zBuffer {
	// 			zBuffer[i] = math.Inf(1)
	// 		}

	// 		angle += 0.01
	// 		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// 		renderRect(-angle*1.5, 0, 0, 10, 50, 0, termbox.ColorCyan, true)
	// 		renderRect(angle, 0, 0, 60, 20, 1, termbox.ColorGreen, true)
	// 		renderRect(angle*2.0, 10, 10, 10, 10, -1, termbox.ColorMagenta, false)
	// 		renderImage(GetMario(), 12, 16, math.Pi/2, -10, 0, -1)

	// 		termbox.Flush()
	// 	}
	// }
}

func resize(w, h int) {
	termW, termH = w, h

	bufferSize := termW * termH
	if len(zBuffer) != bufferSize {
		zBuffer = make([]float64, bufferSize)
	}
}

func renderRect(angle, rx, ry, rw, rh, z float64, color termbox.Attribute, ellipse bool) {
	w, h := termW, termH
	cosA, sinA := math.Cos(-angle), math.Sin(-angle)
	rx, ry = rx+float64(termW)/2, ry+float64(termH)/2
	rw, rh = rw/2, rh/2
	minX, minY, maxX, maxY := float64(w), float64(h), 0.0, 0.0
	corners := [4][2]float64{
		{-rw / 2, -rh / 2}, {rw / 2, -rh / 2},
		{rw / 2, rh / 2}, {-rw / 2, rh / 2},
	}

	for _, c := range corners {
		rx := (c[0]*cosA-c[1]*sinA)/0.5 + rx
		ry := (c[0]*sinA + c[1]*cosA) + ry

		if rx < minX {
			minX = rx
		}
		if rx > maxX {
			maxX = rx
		}
		if ry < minY {
			minY = ry
		}
		if ry > maxY {
			maxY = ry
		}
	}

	startX, startY := int(math.Max(0, math.Floor(minX-1))), int(math.Max(0, math.Floor(minY-1)))
	endX, endY := int(math.Min(float64(w-1), math.Ceil(maxX+1))), int(math.Min(float64(h-1), math.Ceil(maxY+1)))
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			mask := 0
			for iy := range 4 {
				for ix := range 2 {
					sampleX, sampleY := float64(x)+float64(ix)*0.5, float64(y)+float64(iy)*0.25

					if !ellipse && isPointInRotatedRect(sampleX, sampleY, rx, ry, rw, rh, cosA, sinA) {
						mask |= brailleMap[iy][ix]
					}
					if ellipse && isPointInRotatedEllipse(sampleX, sampleY, rx, ry, rw, rh, cosA, sinA) {
						mask |= brailleMap[iy][ix]
					}
				}
			}

			idx := y*w + x
			if mask == 0 {
				continue
			}
			if z >= zBuffer[idx] {
				if mask == 0xFF && termbox.GetCell(x, y).Bg == termbox.ColorDefault {
					termbox.SetBg(x, y, color)
				}
				continue
			}

			termbox.SetChar(x, y, rune(0x2800+mask))
			termbox.SetFg(x, y, color|termbox.AttrBold)
			if mask == 0xFF {
				termbox.SetBg(x, y, color)
			}
			zBuffer[idx] = z
		}
	}
}

func isPointInRotatedRect(x, y, rx, ry, rw, rh, cos, sin float64) bool {
	dx, dy := (x-rx)*0.5, y-ry
	rotX, rotY := dx*cos-dy*sin, dx*sin+dy*cos
	return math.Abs(rotX) < rw/2.0 && math.Abs(rotY) < rh/2.0
}
func isPointInRotatedEllipse(x, y, rx, ry, rw, rh, cos, sin float64) bool {
	dx, dy := (x-rx)*0.5, y-ry
	rotX, rotY := dx*cos-dy*sin, dx*sin+dy*cos
	a, b := rw/2.0, rh/2.0
	return (rotX*rotX)/(a*a)+(rotY*rotY)/(b*b) <= 1.0
}

func renderImage(pixels []termbox.Attribute, imgW, imgH int, angle, rx, ry, z float64) {
	w, h := termW, termH
	cosA, sinA := math.Cos(-angle), math.Sin(-angle)
	rx, ry = rx+float64(termW)/2, ry+float64(termH)/2

	// Calculate bounding box based on rotated corners
	halfW, halfH := float64(imgW)/2.0, float64(imgH)/2.0
	corners := [4][2]float64{
		{-halfW, -halfH}, {halfW, -halfH},
		{halfW, halfH}, {-halfW, halfH},
	}

	minX, minY, maxX, maxY := float64(w), float64(h), 0.0, 0.0
	for _, c := range corners {
		tx := (c[0]*cosA-c[1]*sinA)/0.5 + rx
		ty := (c[0]*sinA + c[1]*cosA) + ry
		minX, minY = math.Min(minX, tx), math.Min(minY, ty)
		maxX, maxY = math.Max(maxX, tx), math.Max(maxY, ty)
	}

	startX, startY := int(math.Max(0, math.Floor(minX-1))), int(math.Max(0, math.Floor(minY-1)))
	endX, endY := int(math.Min(float64(w-1), math.Ceil(maxX+1))), int(math.Min(float64(h-1), math.Ceil(maxY+1)))

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			mask := 0
			var lastColor termbox.Attribute

			for iy := range 4 {
				for ix := range 2 {
					// 1. Transform screen sub-pixel back to image-space coordinates
					sx, sy := float64(x)+float64(ix)*0.5, float64(y)+float64(iy)*0.25
					dx, dy := (sx-rx)*0.5, sy-ry
					rotX, rotY := dx*cosA-dy*sinA, dx*sinA+dy*cosA

					// 2. Map back to 0 -> imgW/imgH range
					imgX := int(rotX + halfW)
					imgY := int(rotY + halfH)

					// 3. Check if we are inside image bounds
					if imgX >= 0 && imgX < imgW && imgY >= 0 && imgY < imgH {
						mask |= brailleMap[iy][ix]
						lastColor = pixels[imgY*imgW+imgX]
					}
					if lastColor == termbox.ColorDefault {
						mask = 0
					}
				}
			}

			idx := y*w + x
			if mask == 0 {
				continue
			}
			if z >= zBuffer[idx] || lastColor == termbox.ColorDefault {
				if mask == 0xFF && termbox.GetCell(x, y).Bg == termbox.ColorDefault {
					termbox.SetBg(x, y, lastColor)
				}
				continue
			}

			termbox.SetChar(x, y, rune(0x2800+mask))
			termbox.SetFg(x, y, lastColor|termbox.AttrBold)
			if mask == 0xFF {
				termbox.SetBg(x, y, lastColor)
			}
			zBuffer[idx] = z
		}
	}
}
func GetMario() []termbox.Attribute {
	R := termbox.ColorRed
	B := termbox.ColorBlue
	Y := termbox.ColorYellow  // Skin tone
	K := termbox.ColorBlack   // Hair/Eyes
	i := termbox.ColorDefault // Transparent

	return []termbox.Attribute{
		i, i, i, R, R, R, R, R, i, i, i, i,
		i, i, R, R, R, R, R, R, R, R, R, i,
		i, i, K, K, K, Y, Y, K, Y, i, i, i,
		i, K, Y, K, Y, Y, Y, K, Y, Y, Y, i,
		i, K, Y, K, K, Y, Y, Y, K, Y, Y, Y,
		i, K, K, Y, Y, Y, Y, K, K, K, K, i,
		i, i, i, Y, Y, Y, Y, Y, Y, Y, i, i,
		i, i, R, R, B, R, R, R, i, i, i, i,
		i, R, R, R, B, R, R, B, R, R, R, i,
		R, R, R, R, B, B, B, B, R, R, R, R,
		Y, Y, R, B, Y, B, B, Y, B, R, Y, Y,
		Y, Y, Y, B, B, B, B, B, B, Y, Y, Y,
		Y, Y, B, B, B, B, B, B, B, B, Y, Y,
		i, i, B, B, B, i, i, B, B, B, i, i,
		i, K, K, K, i, i, i, i, K, K, K, i,
		K, K, K, K, i, i, i, i, K, K, K, K,
	}
}
