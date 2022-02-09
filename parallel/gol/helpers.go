package gol

import (
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

const alive = 255
const dead = 0

func calculateNextState(p Params, world [][]byte, start int, finish int, out chan<- [][]uint8) {

	height := finish - start
	//fmt.Printf("Calculating state %d - %d with height %d\n", start, finish, height)
	width  := p.ImageWidth

	newWorld := make([][]byte, height)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	var state byte
	var neighbours int

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			state = world[i + start][j]
			neighbours = checkNeighbours(p, world, i + start, j)

			if state == alive && neighbours < 2 {
				newWorld[i][j] = dead
			}
			if (state == alive && neighbours == 2) || (state == alive && neighbours == 3) {
				newWorld[i][j] = alive
			}
			if state == alive && neighbours > 3 {
				newWorld[i][j] = dead
			}
			if state == dead && neighbours == 3 {
				newWorld[i][j] = alive
			}
		}
	}
	out <- newWorld
}

func calculateAliveCells(p Params, world [][]byte) []util.Cell {

	var activeCells []util.Cell

	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {

			if world[i][j] == alive {
				newCell := util.Cell{
					X: j,
					Y: i,
				}
				activeCells = append(activeCells, newCell)
			}
		}
	}
	return activeCells
}

func checkNeighbours(p Params, arr [][]byte, x, y int) int {

	var aliveNeighbours = 0
	var row, col int

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {

			if i == 0 && j == 0 {
				continue
			}
			row = (x + i + p.ImageHeight) % p.ImageHeight
			col = (y + j + p.ImageWidth) % p.ImageWidth

			if arr[row][col] == alive {
				aliveNeighbours++
			}
		}
	}
	return aliveNeighbours
}

func worker(p Params, world [][]byte, start int, finish int, out chan<- [][]uint8) {
	calculateNextState(p, world, start, finish, out)
}

func runTicker(c distributorChannels, aliveCells int, turn *int, ticker *time.Ticker) {

	go func() {
		for {
			select {
			case <-ticker.C:
				c.events <- AliveCellsCount{*turn, aliveCells}
			}
		}
	}()
}