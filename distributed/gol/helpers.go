package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

const alive = 255
const dead = 0

func calculateNextState(p Params, world [][]byte) [][]byte {


	//fmt.Printf("Calculating state %d - %d with height %d\n", start, finish, height)


	newWorld := make([][]byte, p.ImageHeight)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	var state byte
	var neighbours int

	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {

			state = world[i][j]
			neighbours = checkNeighbours(p, world, i, j)

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
	return newWorld
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
