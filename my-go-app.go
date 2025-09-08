package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
)

var runPeriod = 300

const numRows = 25
const numCols = 25
const cellSize = 20

var curUpRow = 0
var curUpCol = 0

type gameState struct {
	canvas       js.Value
	topButton    js.Value
	bottomButton js.Value
	leftButton   js.Value
	rightButton  js.Value
	// Game of life grid; true means the cell is alive.
	grid [numRows][numCols]bool

	running bool
}

func main() {
	fmt.Println("Go Web Assembly")
	doc := js.Global().Get("document")

	state := gameState{
		canvas:       doc.Call("getElementById", "gameCanvas"),
		leftButton:   doc.Call("getElementById", "leftButton"),
		rightButton:  doc.Call("getElementById", "rightButton"),
		topButton:    doc.Call("getElementById", "topButton"),
		bottomButton: doc.Call("getElementById", "bottomButton"),
	}
	state.initializeGrid()

	state.canvas.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			x := args[0].Get("offsetX").Int()
			y := args[0].Get("offsetY").Int()
			state.onCanvasClicked(x, y)
			return nil
		}))
	state.topButton.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			state.onTopClicked()
			return nil
		}))
	state.bottomButton.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			state.onBottomClicked()
			return nil
		}))
	state.leftButton.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			state.onLeftClicked()
			return nil
		}))
	state.rightButton.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			state.onRightClicked()
			return nil
		}))

	state.draw()
	state.start()
	select {}
}

func (state *gameState) initializeGrid() {
	state.grid[rand.Intn(numRows)][rand.Intn(numCols)] = true
}

func (state *gameState) onCanvasClicked(x, y int) {
	defer state.draw()
	r := y / cellSize
	c := x / cellSize

	state.grid[r][c] = !state.grid[r][c]
	state.draw()
}

func (state *gameState) onTopClicked() {
	defer state.draw()
	curUpRow = 1
	curUpCol = 0
	state.advance(curUpRow, curUpCol)
}

func (state *gameState) onBottomClicked() {
	defer state.draw()
	curUpRow = -1
	curUpCol = 0
	state.advance(curUpRow, curUpCol)
}

func (state *gameState) onLeftClicked() {
	defer state.draw()
	curUpRow = 0
	curUpCol = -1
	state.advance(curUpRow, curUpCol)
}

func (state *gameState) onRightClicked() {
	defer state.draw()
	curUpRow = 0
	curUpCol = 1
	state.advance(curUpRow, curUpCol)
}

func (state *gameState) start() {
	defer state.draw()

	state.running = true

	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
		state.runningStep()
		return nil
	}), runPeriod)
}

func (state *gameState) runningStep() {
	if state.running {
		state.advance(curUpRow, curUpCol)
		state.draw()
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
			state.runningStep()
			return nil
		}), runPeriod)
	}
}

func (state *gameState) draw() {
	cctx := state.canvas.Call("getContext", "2d")
	// Cells
	for r := range numRows {
		for c := range numCols {
			if state.grid[r][c] {
				cctx.Set("fillStyle", "#baba18")
			} else {
				cctx.Set("fillStyle", "#f4f4f4")
			}

			cctx.Call("fillRect", c*cellSize, r*cellSize, cellSize, cellSize)
		}
	}

	// Horizontal lines
	cctx.Set("strokeStyle", "#000000")
	for r := range numRows {
		cctx.Call("beginPath")
		cctx.Call("moveTo", 0, (r+1)*cellSize-1)
		cctx.Call("lineTo", numCols*cellSize, (r+1)*cellSize-1)
		cctx.Call("stroke")
	}

	// Vertical lines
	for c := range numCols {
		cctx.Call("beginPath")
		cctx.Call("moveTo", (c+1)*cellSize-1, 0)
		cctx.Call("lineTo", (c+1)*cellSize-1, numRows*cellSize)
		cctx.Call("stroke")
	}
}

// advance updates the game state to the next generation
func (state *gameState) advance(upRow int, upCol int) {

	newGrid := [numRows][numCols]bool{}
	for r := range numRows {
		for c := range numCols {
			if state.grid[r][c] {
				newGrid[r][c] = state.grid[r-upRow][c-upCol]
			} else {
				newGrid[r][c] = state.grid[r-upRow][c-upCol]
			}
		}
	}

	state.grid = newGrid
}
