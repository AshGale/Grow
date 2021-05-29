package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth, windowHeight int = 800, 600

const cellSize int = 100
const light float32 = 100 //default should be 100
const water float32 = 100 //default should be 100
const plantStatingWater int = 100
const plantStantingEnergy int = 100

var leafTemplate = leaves{1,
	[]position{
		{0, 0}, {0, -1}, {0, -2}, {1, -2}, {0, -3}, {0, -4}, {-1, -4}, {0, -5}, {1, -5}, {0, -6},
		{0, -7}, {-1, -7}, {1, -8}, {2, -6}, {-1, -8}, {-2, -4}, {-3, -3}, {2, -3}, {-3, -5}, {-1, -9},
		{-2, -10}, {2, -9}, {3, -7}, {3, -10}, {0, -10}, {-3, -11}, {-1, -1}, {-1, -2}, {-1, -3}, {-1, -5},
		{-1, -6}, {-2, -12}, {-4, -12}, {-2, -9}, {-3, -10}, {1, -7}, {2, -8}, {3, -9}, {2, -11}, {3, -11},
		{5, -12}, {6, -13}, {1, -12}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}},
}

var up = direction{0, -1}
var down = direction{0, +1}
var left = direction{-1, 0}
var right = direction{+1, 0}

//var sides = []direction{up, down, left, right}

//------------------------------------end Global variables and constants

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

type leaves struct {
	efficiency float32
	leaf       []position
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

func getRandomCell(posX, posY int) cell {
	return cell{position{posX, posY}, ground{water}, getRandomVegetaion(), sun{light, light}}
}

func getRandomVegetaion() vegetation {

	ammount := getRandomInt(cellSize/2, 2)

	var plants []plant
	for i := 0; i < ammount; i++ {
		//todo add sanity check for if plan is on pixel already and regenerate
		plants = append(plants, getRandomPlant())
	}
	fmt.Printf("Creating vegetaion with %d plants \n%v\n", ammount, plants)

	return vegetation{plants, 1}
}

func getRandomPlant() plant {
	return plant{position{getRandomInt(cellSize-1, 0), getRandomInt(cellSize-1, 0)},
		getRandomInt(40, 1), float32(getRandomInt(plantStantingEnergy, 0)), float32(getRandomInt(plantStatingWater, 0))}
}

func getRandomInt(max, min int) int {
	time.Sleep(time.Microsecond) //was generating the same number
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//------------------------------------end random helper functions

func (cell *cell) calculateDesity() {
	var density int
	var plants = cell.plants
	for p := 0; p < len(plants); p++ {
		density += plants[p].size
	}
	cell.density = density
}

//------------------------------------end calculation functions

func (plant *plant) draw(cell *cell, pixels []byte, wg *sync.WaitGroup) {

	for i := 0; i < plant.size; i++ {
		var x, y = cell.x + plant.x + leafTemplate.leaf[i].x, cell.y + plant.y + leafTemplate.leaf[i].y
		setPixle(x, y, color{0, 255, 0}, pixels)
	}
	wg.Done()
}

func (cell *cell) draw(pixels []byte) {
	var wg sync.WaitGroup
	wg.Add(len(cell.plants))

	for p := 0; p < len(cell.plants); p++ {
		cell.plants[p].draw(cell, pixels, &wg)
	}
	wg.Wait()
	//fmt.Printf(" - vegetaion %v %v\n", cell.density, cell.plants) //pring status of the plants
}

//------------------------------------end draw functions

func (plant *plant) update(density int, light float32, humidity *float32, wg *sync.WaitGroup) {

	//determine the % of resources this plant gets based on size
	share := float32(plant.size) / float32(density)

	//calculate the ammount of sunshime recived
	sunShine := light * float32(share) //not * 100 due to larger light value

	//calculate the ammount of water taken from the cell,
	absorbedWater := share * *humidity
	*humidity -= absorbedWater
	plant.water += absorbedWater

	// fmt.Printf("  > plant%v share:%v energy:%v sunShine:%v water:%v huumidity:%v\n",
	// plant.size, share, plant.energy, sunShine, plant.water, *humidity)
	if plant.water > sunShine {
		plant.energy += sunShine
		plant.water -= sunShine
	} else {
		plant.energy += plant.water
		plant.water = 0
	}

	//here todo, add in the abbility to grow in size<<<<<<<<<<<<<<<<<<<

	//calculate plant upkeep for being alive// seed, flower, growth costs
	plant.energy -= float32(plant.size)
	if plant.energy <= 0 {
		plant.size--
		plant.energy = 0
	}
	wg.Done()
}

func (cell *cell) update() {

	cell.calculateDesity()
	cell.humidity += water //todo have sepearte water array for the cells

	var wg sync.WaitGroup
	wg.Add(len(cell.plants))

	for p := 0; p < len(cell.plants); p++ {
		cell.plants[p].update(cell.density, light, &cell.humidity, &wg)
	}
	wg.Wait()
	//loop throughall to see if died//could have alive flag instead
	for p := 0; p < len(cell.plants); p++ {
		if cell.plants[p].size <= 0 {
			if p == (len(cell.plants) - 1) {
				cell.plants = append(cell.plants[:p], nil...)
			} else {
				cell.plants = append(cell.plants[:p], cell.plants[p+1:]...)
			}
			p--
		}
	}

}

//------------------------------------end update functions

//------------------------------------start window interaction functions

func setPixle(x, y int, c color, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

// func getPixle(x, y int, pixels []byte) color {
// 	index := (y*windowWidth + x) * 4

// 	color := color{0, 0, 0}

// 	if index < len(pixels)-4 && index >= 0 {
// 		color.r = pixels[index]
// 		color.g = pixels[index+1]
// 		color.b = pixels[index+2]
// 	}
// 	return color
// }

func clearScreen(pixels []byte) {
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			setPixle(x, y, color{0, 0, 0}, pixels)
		}
	}
}

//------------------------------------end window interaction functions

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

	window, err := sdl.CreateWindow("stl2 Grow Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

		cell.update()
		cell.draw(pixels)

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		//sdl.Delay(1024) //wait 1 second
		sdl.Delay(16)
	}
}
