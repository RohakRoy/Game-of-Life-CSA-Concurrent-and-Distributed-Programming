package gol

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {

	c.ioCommand <- ioInput

	filename := strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight)
	c.ioFilename <- filename

	factor := p.ImageHeight / p.Threads
	multiplier := 1
	heightSlice := make([]int, p.Threads)

	for i := 0; i < p.Threads; i++ {
		heightSlice[i] = factor * multiplier
		multiplier++
	}

	worldSlice := make([]chan [][]uint8, p.Threads)
	for i := range worldSlice {
		worldSlice[i] = make(chan [][]uint8, factor*p.ImageWidth)
	}

	// TODO: Create a 2D slice to store the world.

	world := make([][]byte, p.ImageHeight) //created an empty 2D world
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	for y := 0; y < p.ImageHeight; y++ { //copied the image into my 2D world
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.ioInput
			world[y][x] = val
		}
	}

	turn := 0
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				newCell := util.Cell{
					X: x,
					Y: y,
				}
				c.events <- CellFlipped{turn,newCell}
			}
		}
	}

	pause := make(chan bool, 1)
	var check bool
	check = true

	go func() {
		for {
			select {
			case val := <-keyPresses:
				switch val {
				case 's':
					c.ioCommand <- ioOutput

					newFileName := filename + "x" + strconv.Itoa(turn)
					c.ioFilename <- newFileName

					for y := 0; y < p.ImageHeight; y++ {
						for x := 0; x < p.ImageWidth; x++ {
							c.ioOutput <- world[y][x]
						}
					}

					c.events <- ImageOutputComplete{turn, filename}

				case 'q':
					c.ioCommand <- ioOutput

					newFileName := filename + "x" + strconv.Itoa(turn)
					c.ioFilename <- newFileName

					for y := 0; y < p.ImageHeight; y++ {
						for x := 0; x < p.ImageWidth; x++ {
							c.ioOutput <- world[y][x]
						}
					}
					c.events <- ImageOutputComplete{turn, filename}

					c.events <- FinalTurnComplete{turn, calculateAliveCells(p, world)}

				case 'p':
					if check {
						fmt.Printf("Paused at turn %d \n", turn + 1)
					} else {
						fmt.Println("Continuing")
						pause <- true
					}
					check = !check
				}
			}
		}
	}()

	var finalState FinalTurnComplete
	var m sync.Mutex
	var m2 sync.Mutex

	// TODO: Execute all turns of the Game of Life.

	done := make(chan bool)
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				{
					m.Lock()
					m2.Lock()
					c.events <- AliveCellsCount{turn, len(calculateAliveCells(p, world))}
					m2.Unlock()
					m.Unlock()
				}
			}
		}
	}()

	if p.Turns == 0 {
		finalState.CompletedTurns = turn
		finalState.Alive = calculateAliveCells(p, world)

	} else {
		for turn < p.Turns {

			if check == true {
				pause <- true
			}

			<- pause

			for k := 0; k < p.Threads; k++ {
				if k == 0 {
					go calculateNextState(p, world, 0, heightSlice[k], worldSlice[k])
				} else {
					if k == p.Threads-1 {
						go calculateNextState(p, world, heightSlice[k-1], p.ImageHeight, worldSlice[k])
					} else {
						go calculateNextState(p, world, heightSlice[k-1], heightSlice[k], worldSlice[k])
					}
				}
			}

			var updatedWorld [][]uint8

			for i := 0; i < p.Threads; i++ {
				matrix := <-worldSlice[i]
				updatedWorld = append(updatedWorld, matrix...)
			}

			for y := 0; y < p.ImageHeight; y++ {
				for x := 0; x < p.ImageWidth; x++ {
					if updatedWorld[y][x] != world[y][x] {
						newCell := util.Cell{
							X: x,
							Y: y,
						}
						c.events <- CellFlipped{turn,newCell}
					}
				}
			}

			m.Lock()
			world = updatedWorld
			m.Unlock()

			c.events <- TurnComplete{turn}

			m2.Lock()
			turn++
			m2.Unlock()

		}
		finalState.CompletedTurns = turn
		finalState.Alive = calculateAliveCells(p, world)
	}
	// TODO: Report the final state using FinalTurnCompleteEvent.

	c.ioCommand <- ioOutput

	newFileName := filename + "x" + strconv.Itoa(turn)
	c.ioFilename <- newFileName

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- world[y][x]
		}
	}
	c.events <- ImageOutputComplete{turn, filename}

	done <- true
	c.events <- finalState

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
