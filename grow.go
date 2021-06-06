package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const plantNameFile = "plant.json"
const leafShapesFile = "leafShapes.json"
const partShapesFile = "partShapes.json"

const windowWidth, windowHeight int = 340, 340

const cellSize int = 10
const cellsX int = 32
const cellsY int = 32
const lightAmmount float32 = 500 //default should be 100
const waterAmmount float32 = 500 //default should be 100
const plantStatingWater int = 100
const plantStantingEnergy int = 100
const plantMaxWater int = 1000
const plantMaxEnergy int = 1000
const cellStartingPlants = 10

//export this to a json file at some stage, or csv ect
var plantIndexCounter Counter = Counter{0}
var plantColor Color = Color{0, 128, 0}

var leafShapes [][]Point
var partShapes ShapeSet

// var plantBodyTemplate = plantBody {1,
// 	[]plantPart{
// 		{at{0, 0},clr{0,255,0}}, {at{0, -1},clr{0,255,0}}, {at{0, -2},clr{0,255,0}}, {at{1, -2},clr{0,255,0}}, {at{0, -3},clr{0,255,0}}, {at{0, -4},clr{0,255,0}}, {at{-1, -4},clr{0,255,0}},
// 	},
// }

// var leafTemplate = plantBody{1,
// 	[]Position{
// 		{0, 0}, {0, -1}, {0, -2}, {1, -2}, {0, -3}, {0, -4}, {-1, -4}, {0, -5}, {1, -5}, {0, -6},
// 		{0, -7}, {-1, -7}, {1, -8}, {2, -6}, {-1, -8}, {-2, -4}, {-3, -3}, {2, -3}, {-3, -5}, {-1, -9},
// 		{-2, -10}, {2, -9}, {3, -7}, {3, -10}, {0, -10}, {-3, -11}, {-1, -1}, {-1, -2}, {-1, -3}, {-1, -5},
// 		{-1, -6}, {-2, -12}, {-4, -12}, {-2, -9}, {-3, -10}, {1, -7}, {2, -8}, {3, -9}, {2, -11}, {3, -11},
// 		{5, -12}, {6, -13}, {1, -12}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 		{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
// 	},
// }

//export this to a json file at some stage, or csv ect
var growthTemplate = growth{1,
	[]int{
		1, 4, 8, 12, 16, 21, 28, 35, 42, 50,
		59, 67, 75, 84, 93, 103, 113, 125, 136, 148,
		162, 185, 199, 214, 229, 244, 260, 276, 293, 310,
		328, 346, 364, 383, 402, 422, 442, 463, 484, 505,
		527, 549, 573, 597, 10000,
	},
}

//var sides = []direction{up, down, left, right}

//------------------------------------end Global variables and constants
type Counter struct {
	Count int
}

type Color struct {
	R, G, B byte
}

type Position struct {
	X int
	Y int
}

type Plant struct {
	Index int
	Position
	Mass   int
	Energy Energy
	Water  Water
	GrowAt int
	Parts  []PlantPart //a plant will be made up a bunch of different parts
}

//plant part will take water from plant
type PlantPart struct {
	Mass int //determined by shape
	Position
	End   Position //fill this in when section grows//based on last postion of shape number
	Shape int      //there will be a set of different shapes
	Requirement
	//need variable to determine how much water it can transfer
	Leafs []Leaf
}

//leafs will take water from the parent plant part
type Leaf struct {
	Mass int
	Position
	GrowAt int //keep track of how much ener
	Shape  int //what groth template to use
}

type Point struct {
	Color Color
	Position
}

type Shape struct {
	Points []Point
}

type ShapeGroup struct {
	Shapes []Shape
}

type ShapeSet struct {
	Groups []ShapeGroup
}

type growth struct {
	pattern      int
	growthStages []int
}

type Vegetation struct {
	Plants  []Plant
	Density int
}

type Ground struct {
	Humidity float32
	//catagory int
	//fertility float32
	//tempreture int
}

type Requirement struct {
	Waterneeded  int //negative number ?
	EnergyNeeded int
}

type Water struct {
	ObtainingMultiplyer float32
	Ammount             float32
	Maximum             int
}

type Energy struct {
	ObtainingMultiplyer float32
	Ammount             float32
	Maximum             int
}

type Sun struct {
	Intensity float32
	Energy    float32
	//heat int
}

// type dot struct {
// 	color color
// 	checked bool `default:false`
// }

type Cell struct {
	Position
	Ground
	Vegetation
	Sun Sun
	//dots [10][10]dot
}

//------------------------------------end structs

func loadPartShapes() {

	file, err := ioutil.ReadFile(partShapesFile)
	if err != nil {
		fmt.Printf("You need this file :/ %v\n", partShapesFile)
		createPlantShapes(&partShapes, 20)
		file, _ := json.MarshalIndent(partShapes, "", "\t")
		_ = ioutil.WriteFile(partShapesFile, file, 0644)
	} else {
		fmt.Printf("Loading %v...\n", partShapesFile)
		_ = json.Unmarshal([]byte(file), &partShapes)
	}
}
func loadLeafShapes(leafShapes *[][]Point) {

	file, err := ioutil.ReadFile(leafShapesFile)
	if err != nil {
		fmt.Printf("You need this file :/ %v\n", leafShapesFile)
		shapes := make([][]Point, 2)
		shapes[0] = []Point{{Color{0, 0, 0}, Position{0, 0}}}
		shapes[1] = []Point{{Color{0, 0, 0}, Position{0, 0}}}
		file, _ := json.MarshalIndent(shapes, "", "\t")
		_ = ioutil.WriteFile(leafShapesFile, file, 0644)
	} else {
		fmt.Printf("Loading %v...\n", leafShapesFile)
		_ = json.Unmarshal([]byte(file), &leafShapes)
	}
}

func loadOrCreatePlant(plant *Plant) {
	file, err := ioutil.ReadFile(plantNameFile)
	if err != nil {
		fmt.Printf("No Plant found, creating %v\n", plantNameFile)
		*plant = getRandomPlant()
		file, _ := json.MarshalIndent(plant, "", "\t")
		_ = ioutil.WriteFile(plantNameFile, file, 0644)
	} else {
		fmt.Printf("Loading %v...\n", plantNameFile)
		_ = json.Unmarshal([]byte(file), &plant)
	}
}

//------------------------------------end file io functions

func setUpCells(cells *[]Cell) {
	var wg sync.WaitGroup

	wg.Add(cellsY * cellsX)
	var cellNumber = 0
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			var posX, posY = (x * cellSize) + cellSize, (y * cellSize) + cellSize
			go addCellToCell(posX, posY, cells, &wg)
			fmt.Printf("setting up cell %v of %v}\n", cellNumber, cellsX*cellsY)
			cellNumber++
		}
	}
	fmt.Printf("Finalizing...\n")
	wg.Wait()
}

func addCellToCell(posX, posY int, cells *[]Cell, wg *sync.WaitGroup) {
	cell := getRandomCell(posX, posY)
	*cells = append(*cells, cell)
	wg.Done()
}

func getRandomCell(posX, posY int) Cell {
	return Cell{Position{posX, posY}, Ground{waterAmmount}, getRandomVegetaion(), Sun{lightAmmount, lightAmmount}}
}

func getRandomVegetaion() Vegetation {

	ammount := getRandomInt(cellStartingPlants, 1)

	var plants []Plant
	for i := 0; i < ammount; i++ {
		//todo add sanity check for if plan is on pixel already and regenerate
		plants = append(plants, getRandomPlant())
	}
	//fmt.Printf("Creating vegetaion with %d plants \n%v\n", ammount, plants)

	return Vegetation{plants, 1}
}

func getRandomPlant() Plant {
	//NOTE not threadsafe due to index tracker
	randomSize := getRandomInt(15, 1)
	water := Water{1, float32(getRandomInt(plantStatingWater, 0)), plantMaxWater}
	energy := Energy{1, float32(getRandomInt(plantStantingEnergy, 0)), plantMaxEnergy}
	return Plant{plantIndexCounter.next(), Position{getRandomInt(cellSize-1, 0), getRandomInt(cellSize-1, 0)},
		randomSize, energy, water, 1, []PlantPart{}}
}

func getRandomInt(max, min int) int {
	time.Sleep(time.Microsecond) //was generating the same number
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//------------------------------------end helper functions

func (counter *Counter) currentValue() int {
	return counter.Count
}

func (counter *Counter) increment() {
	counter.Count += 1
}

func (counter *Counter) next() int {
	counter.increment()
	return counter.currentValue()
}

func (cell *Cell) calculateDesity() {
	var density int
	var plants = cell.Plants
	for p := 0; p < len(plants); p++ {
		density += plants[p].Mass
	}
	cell.Density = density
}

//------------------------------------end calculation functions

func (plant *Plant) draw(cell *Cell, pixels *[]byte) {
	var plantGroup sync.WaitGroup
	plantGroup.Add(len(plant.Parts))

	for p := 0; p < len(plant.Parts); p++ {
		go plant.Parts[p].draw(pixels, &plantGroup)
	}
	plantGroup.Wait()
	//fmt.Printf(" - vegetaion %v %v\n", cell.density, cell.plants) //pring status of the plants

}

func (part *PlantPart) draw(pixels *[]byte, plantGroup *sync.WaitGroup) {
	var partGroup sync.WaitGroup
	partGroup.Add(len(part.Leafs))

	//draw the body of the part

	//drawe the all the leaves on the part
	for p := 0; p < len(part.Leafs); p++ {
		go part.Leafs[p].draw(pixels, &partGroup)
	}
	partGroup.Wait()
	//fmt.Printf(" - vegetaion %v %v\n", cell.density, cell.plants) //pring status of the plants
	plantGroup.Done()
}

func (leaf *Leaf) draw(pixels *[]byte, partGroup *sync.WaitGroup) {

	// might be lots of overlap if each stage of the leaf is a full set of points ?

	// for i := 0; i < leaf.Size; i++ {
	// 	var x, y = leaf.X + leafShapes[leaf.Shape].point[i].X, leaf.Y + leafShapes[leaf.Shape].point[i].Y
	// 	setPixle(x, y, template[i].color, *pixels)
	// }
	partGroup.Done()
}

// func (cell *Cell) draw(pixels []byte, mainLoop *sync.WaitGroup) {

// }

//------------------------------------end draw functions

func (plant *Plant) update() {

	// //determine the % of resources this plant gets based on size
	// share := float32(plant.Size) / float32(density)

	// //calculate the ammount of sunshime recived
	// sunShine := *light * float32(share) * plant.Energy.ObtainingMultiplyer //not * 100 due to larger light value

	// //calculate the ammount of water taken from the cell,
	// absorbedWater := share * *humidity * plant.Water.ObtainingMultiplyer

	// //logic to apply the water to the plant
	// if (plant.Water.Ammount + absorbedWater) > float32(plant.Water.Maximum) {
	// 	//only take enough to max out plant water
	// 	*humidity -= (float32(plant.Water.Maximum) - plant.Water.Ammount)
	// 	plant.Water.Ammount = float32(plant.Water.Maximum)
	// } else {
	// 	*humidity -= absorbedWater
	// 	plant.Water.Ammount += absorbedWater
	// }

	// //logic to apply the sunshine and energy to the plant
	// if plant.Water.Ammount > sunShine {
	// 	if plant.Energy.Ammount+sunShine > float32(plant.Energy.Maximum) {
	// 		plant.Water.Ammount -= float32(plant.Energy.Maximum - int(plant.Energy.Ammount))
	// 		plant.Energy.Ammount = float32(plant.Energy.Maximum)
	// 	} else {
	// 		plant.Energy.Ammount += sunShine
	// 		plant.Water.Ammount -= sunShine
	// 	}

	// } else {
	// 	plant.Energy.Ammount += plant.Water.Ammount
	// 	plant.Water.Ammount = 0
	// }

	// //calculate plant upkeep for being alive// seed, flower, growth costs
	// plant.Energy.Ammount -= float32(plant.Size)
	// if plant.Energy.Ammount <= 0 {
	// 	plant.Size--
	// 	plant.Energy.Ammount = 0
	// } else if plant.Energy.Ammount >= float32(plant.GrowAt) {
	// 	plant.Energy.Ammount -= float32(plant.GrowAt)
	// 	plant.Size++
	// 	plant.GrowAt = growthTemplate.growthStages[plant.Size] //todo reference growthPattern
	// }
	// wg.Done()
}

func (cell *Cell) update(mainLoop *sync.WaitGroup) {

	// cell.calculateDesity()
	// cell.Humidity += float32(getRandomInt(int(waterAmmount), 0)) //todo have sepearte water array for the cells
	// cell.Sun.Energy = float32(getRandomInt(int(lightAmmount), 0))

	// var wg sync.WaitGroup
	// wg.Add(len(cell.Plants))

	// for p := 0; p < len(cell.Plants); p++ {
	// 	go cell.Plants[p].update(cell.Density, &cell.Sun.Energy, &cell.Humidity, &wg)
	// }
	// wg.Wait()
	// //loop throughall to see if died//could have alive flag instead
	// for p := 0; p < len(cell.Plants); p++ {
	// 	if cell.Plants[p].Size <= 0 {
	// 		if p == (len(cell.Plants) - 1) {
	// 			cell.Plants = append(cell.Plants[:p], nil...)
	// 		} else {
	// 			cell.Plants = append(cell.Plants[:p], cell.Plants[p+1:]...)
	// 		}
	// 		p--
	// 	}
	// }
	// mainLoop.Done()
}

//------------------------------------end update functions

//------------------------------------start window interaction functions

func setPixle(x, y int, c Color, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
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

func clearScreen(pixels *[]byte, color Color) {
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			setPixle(x, y, color, *pixels)
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

	//
	var plant Plant
	loadOrCreatePlant(&plant)
	fmt.Printf("Your Plant: %+v\n", plant)
	//
	loadLeafShapes(&leafShapes)
	fmt.Printf("Number of LeafShapes %v\n", len(leafShapes))

	loadPartShapes()
	fmt.Printf("Number of PartShapes %v\n", len(partShapes.Groups))

	//
	fmt.Println("Setup Done ... \nStart game loop ")

	//cd grow && doskey /listsize=0 && go build -o grow.exe && grow.exe
	//------------------------------------Game loop
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clearScreen(&pixels, Color{25, 12, 8})

		//plant.update()
		//plant.draw(&pixels)

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(2024) //wait 1 second
		//sdl.Delay(16)
	}
}

func createPlantShapes(set *ShapeSet, length int) {
	var up = Position{0, -1}
	var upRight = Position{+1, -1}
	var right = Position{+1, 0}
	var downRight = Position{+1, +1}
	var down = Position{0, +1}
	var downLeft = Position{-1, +1}
	var left = Position{-1, 0}
	var upLeft = Position{-1, -1}

	directions := []Position{up, upRight, right, downRight, down, downLeft, left, upLeft}
	numberDirections := len(directions)

	*set = ShapeSet{} //make([]ShapeGroup, length-2)
	set.Groups = make([]ShapeGroup, length)
	for g := 2; g < len(set.Groups); g++ {
		set.Groups[g].Shapes = make([]Shape, numberDirections)
	}

	for l := 2; l < length; l++ { //for the ammount of points in a shape

		for d := 0; d < numberDirections; d++ { //8, for each direct
			//fmt.Printf("l=%v d=%v\n", l, d)
			var shape Shape
			for s := 1; s <= l; s++ {
				point := Point{Color{0, byte(200 - l), 0}, Position{s * directions[d].X, s * directions[d].Y}}
				//fmt.Printf("\tPoint: %v -> direction %v s=%v\n", point, directions[d], s)
				shape.Points = append(shape.Points, point)
			}
			fmt.Printf("%v\n", shape)

			set.Groups[l].Shapes[d] = shape
		}
	}
	//fmt.Printf("%+v\n", shapes)
}
