package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

const alive = 255
const dead  = 0

func calculateNextState(height, width int, world [][]byte) [][]byte {


	//fmt.Printf("Calculating state %d - %d with height %d\n", start, finish, height)


	newWorld := make([][]byte, height)
	for i := range newWorld {
		newWorld[i] = make([]byte, width)
	}

	var state byte
	var neighbours int

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			state = world[i][j]
			neighbours = checkNeighbours(height, width, world, i, j)

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

func checkNeighbours(height, width int, arr [][]byte, x, y int) int {

	var aliveNeighbours = 0
	var row, col int

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {

			if i == 0 && j == 0 {
				continue
			}
			row = (x + i + height) % height
			col = (y + j + width) % width

			if arr[row][col] == alive {
				aliveNeighbours++
			}
		}
	}
	return aliveNeighbours
}

func calculateAliveCells(height, width int, world [][]byte) []util.Cell {

	var activeCells []util.Cell

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

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

var aliveCells int
var turns int
var currentWorld [][]byte
var world sync.Mutex
var turnLock sync.Mutex
var alock sync.Mutex

type GameOfLifeOperations struct {}

func (s *GameOfLifeOperations) EndServer(req stubs.Request, res *stubs.Response) (err error) {

	fmt.Println("SHUTTING DOWN SERVER")
	os.Exit(0)
	return
}

func (s *GameOfLifeOperations) CurrentWorldFinder(req stubs.ReqWorld, res *stubs.ResWorld) (err error) {
	res.World = currentWorld
	res.Turn = turns
	return
}

func (s *GameOfLifeOperations) AliveCellsFinder(req stubs.Request, res *stubs.Response) (err error) {

	alock.Lock()
	res.AliveCells = aliveCells
	alock.Unlock()

	turnLock.Lock()
	res.Turns = turns
	turnLock.Unlock()

	return
}

func (s *GameOfLifeOperations) ProcessTurns(req stubs.Request, res *stubs.Response) (err error) {

	turn := 0

		for turn < req.Turns {

			world.Lock()
			req.World = calculateNextState(req.ImageHeight, req.ImageWidth, req.World)
			world.Unlock()

			currentWorld = req.World

			fmt.Printf("Finished turn %d \n", turn)

			world.Lock()
			alock.Lock()
			aliveCells = len(calculateAliveCells(req.ImageHeight, req.ImageWidth, req.World))
			alock.Unlock()
			world.Unlock()


			turn = turn + 1

			turnLock.Lock()
			turns = turn
			turnLock.Unlock()
		}
	res.Turns = turn
	res.World = req.World
	return
}

func main(){
	pAddr := flag.String("port","8030","Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLifeOperations{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
