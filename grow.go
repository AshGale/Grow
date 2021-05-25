package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth, windowHeight int = 800, 600
const light float32 = 3
const water float32 = 3

var up = direction{0, -1}
var down = direction{0, +1}
var left = direction{-1, 0}
var right = direction{+1, 0}
var sides = []direction{up, down, left, right}

type direction struct {
	x, y int
}

type color struct {
	r, g, b byte
}

type position struct {
	x, y int
}

type plant struct {
	position
	size   int
	energy float32
	water  float32
	//growth int
}

type vegetation struct {
	plants  []plant
	density int
}

type ground struct {
	humidity float32
	//catagory int
	//fertility float32
	//tempreture int
}

type sun struct {
	intensity float32
	energy    float32
	//heat int
}

// type dot struct {
// 	color color
// 	checked bool `default:false`
// }

type cell struct {
	position
	ground
	vegetation
	sun sun
	//dots [10][10]dot
}

func (plant *plant) draw(cell *cell, pixels []byte) {

	setPixle(int(cell.x+plant.x), int(cell.y+plant.y), color{0, 255, 0}, pixels)

	if plant.size == 1 {
		return
	}

	// if cell.dots[0][0].checked {
	// }

	//while notFinished -> draw plant
	//almost a pathfinding algorithm
	fmt.Printf("drawing plant of size %d at %d,%d\n", plant.size, cell.x+plant.x, cell.x+plant.y)

	var targetX, targetY = cell.x + plant.x, cell.y + plant.y

	for i := 0; i < plant.size; {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(3-0) + 0

		side := sides[index]
		var x, y = targetX + side.x, targetY + side.y

		fmt.Printf("testing pixel at %d,%d = %v\n", x, y, getPixle(x, y, pixels))
		if getPixle(x, y, pixels).g != 0 {
			fmt.Printf("pixel had plant there: %d, %d \n", x, y)
			targetX = x
			targetY = y
		} else {
			i++
			setPixle(x, y, color{0, 255, 0}, pixels)
		}
	}
	//rand.Seed(time.Now().UnixNano())
	//index := rand.Intn(3-0) + 0

	//figure out how to draw a plant that is bigger than 1
	//directions :=

	//what i want to do here is to have the central pixel drawn, then draw around it randomly

	//setPixle(int(plant.x), int(plant.y), color{0, 255, 0}, pixels)
}

func (cell *cell) draw(pixels []byte) {
	for p := 0; p < len(cell.plants); p++ {
		cell.plants[p].draw(cell, pixels)
	}

}

func (cell *cell) update() {

	//light and water constants for now

	// figure out density
	//for each plant -> get resources -> photosynthesis -> upkeep

}

func setPixle(x, y int, c color, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func getPixle(x, y int, pixels []byte) color {
	index := (y*windowWidth + x) * 4

	color := color{0, 0, 0}

	if index < len(pixels)-4 && index >= 0 {
		color.r = pixels[index]
		color.g = pixels[index+1]
		color.b = pixels[index+2]
	}
	return color
}

func clearScreen(pixels []byte) {
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			setPixle(x, y, color{0, 0, 0}, pixels)
		}
	}
}

func main() {
	//go build -o grow.exe && grow.exe
	fmt.Print("Started Grow\n")
	//Setup //https://stackoverflow.com/questions/38948418/cross-compiling-hello-world-on-mac-for-android
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

	plants := []plant{{position{4, 4}, 12, 100, 100}, {position{8, 8}, 8, 10, 10}}
	cell := cell{position{(windowWidth / 2), windowHeight / 2}, ground{water}, vegetation{plants, 1}, sun{light, light}}

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

		cell.draw(pixels)
		cell.update()

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(1024)
	}
}
