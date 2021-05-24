package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth, windowHeight int = 800, 600
const speed float32 = 5

type color struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

type velocity struct {
	Vx, Vy float32
}

type animal struct {
	position
	velocity
	size int
	color
}

func (animal *animal) draw(pixels []byte) {
	for y := -animal.size; y < animal.size; y++ {
		for x := -animal.size; x < animal.size; x++ {
			setPixle(int(animal.x)+x, int(animal.y)+y, animal.color, pixels)
		}
	}
}

func (animal *animal) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		animal.y -= speed
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		animal.y += speed
	}
	if keyState[sdl.SCANCODE_LEFT] != 0 {
		animal.x -= speed
	}
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		animal.x += speed
	}
}

func setPixle(x, y int, c color, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clearScreen(pixels []byte) {
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			setPixle(x, y, color{0, 0, 0}, pixels)
		}
	}
}

func main() {
	//Setup
	//go get -v github.com/veandco/go-sdl2/sdl@master
	//go mod init
	//go mod tidy
	//go mod vendor
	//go run stl2.go

	//Added for mecosx issue
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("stl2 PONG Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Print(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING,
		int32(windowWidth), int32(windowHeight))
	if err != nil {
		fmt.Print(err)
	}
	defer texture.Destroy()

	pixels := make([]byte, windowHeight*windowWidth*4)
	keyState := sdl.GetKeyboardState()

	animal := animal{position{float32(windowWidth / 2), float32(windowHeight / 2)}, velocity{0, 0}, 10, color{255, 255, 255}}

	//go build -o red.exe && red.exe
	//Gameloop
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clearScreen(pixels)

		animal.draw(pixels)
		animal.update(keyState)

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}
