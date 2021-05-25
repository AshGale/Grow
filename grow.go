package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth, windowHeight int = 800, 600

//const cellSize int = 10
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

//------------------------------------end structs

//------------------------------------start random helper functions

func getRandomVegetaion() vegetation {

	ammount := getRandomInt(7, 1)

	var plants []plant
	for i := 0; i < ammount; i++ {
		plants = append(plants, getRandomPlant())
	}
	fmt.Printf("Creating vegetaion with %d plants \n%v\n", ammount, plants)

	return vegetation{plants, 1}
}

func getRandomPlant() plant {
	return plant{position{getRandomInt(9, 0), getRandomInt(9, 0)},
		getRandomInt(12, 1), float32(getRandomInt(128, 16)), float32(getRandomInt(128, 16))}
}

func getRandomInt(max, min int) int {
	//sdl.Delay(1) //was generating the same number
	time.Sleep(time.Microsecond)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//------------------------------------end random helper functions

func (plant *plant) draw(cell *cell, pixels []byte) {

	setPixle(int(cell.x+plant.x), int(cell.y+plant.y), color{255, 255, 255}, pixels)

	if plant.size == 1 {
		return
	}

	var targetX, targetY = cell.x + plant.x, cell.y + plant.y

	for i := 0; i < plant.size; {
		index := getRandomInt(3, 0)

		side := sides[index]
		//todo ensure x , cell.x + cellsize yy''
		var x, y = targetX + side.x, targetY + side.y

		//fmt.Printf("testing pixel at %d,%d = %v\n", x, y, getPixle(x, y, pixels))
		if getPixle(x, y, pixels).g != 0 {
			//fmt.Printf("pixel had plant there: %d, %d \n", x, y)
			targetX = x
			targetY = y
		} else {
			i++
			setPixle(x, y, color{0, 255, 0}, pixels)
		}
	}
}

func (cell *cell) draw(pixels []byte) {
	for p := 0; p < len(cell.plants); p++ {
		cell.plants[p].draw(cell, pixels)
	}

}

func (cell *cell) update() {

	//here

	//light and water constants for now

	// figure out density
	//for each plant -> get resources -> photosynthesis -> upkeep

}

//------------------------------------start window interaction functions

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

	rand.Seed(time.Now().UnixNano())
	pixels := make([]byte, windowHeight*windowWidth*4)

	//------------------------------------end stl2 setup
	//------------------------------------intitilize variables

	cell := getRandomCell(windowWidth/2, windowHeight/2)

	//cd grow && doskey /listsize=0 && go build -o grow.exe && grow.exe
	//------------------------------------Game loop
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
		//sdl.Delay(1024)
		sdl.Delay(16)
	}
}

func getRandomCell(posX, posY int) cell {
	return cell{position{posX, posY}, ground{water}, getRandomVegetaion(), sun{light, light}}
}
