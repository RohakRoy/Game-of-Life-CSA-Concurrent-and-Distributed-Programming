package gol

import (
	"flag"
	"net/rpc"
	"strconv"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

func makeCall(client *rpc.Client, world [][]byte, turn, width, height int) *stubs.Response{
	request := stubs.Request{World: world, Turns: turn, ImageWidth: width, ImageHeight: height}
	response := new(stubs.Response)
	client.Call(stubs.TurnsHandler, request, response)
	return response
}

func makeCallForAliveCells(client *rpc.Client, world [][]byte, turn, width, height int) *stubs.Response{
	request := stubs.Request{World: world, Turns: turn, ImageWidth: width, ImageHeight: height}
	response := new(stubs.Response)
	client.Call(stubs.GetAliveCells, request, response)
	return response
}

func makeCallForCurrentWorld(client rpc.Client) *stubs.ResWorld{
	request := stubs.ReqWorld{}
	response := new(stubs.ResWorld)
	client.Call(stubs.GetCurrentWorld, request, response)
	return response
}

func makeCallToQuitServer(client rpc.Client) {
	request := stubs.ReqQuitServer{}
	response := new(stubs.ResQuitServer)
	client.Call(stubs.QuitServer, request, response)
	return
}

var server = flag.String("server","127.0.0.1:8030","IP:port string to connect to as server")

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {

	c.ioCommand <- ioInput

	filename := strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight)
	c.ioFilename <- filename

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

	// TODO: Execute all turns of the Game of Life.

	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	turn := 0

	done := make(chan bool, 1)
	done2 := make(chan bool, 1)
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				response := makeCallForAliveCells(client, world, p.Turns, p.ImageWidth, p.ImageHeight)

				for turn < response.Turns {
					c.events <-TurnComplete{turn}
					turn++
				}
				c.events <- AliveCellsCount{response.Turns, response.AliveCells}

				if turn >= p.Turns {
					done <-true
					done2 <-true
				}
			}
		}
	}()

	response := makeCall(client, world, p.Turns, p.ImageWidth, p.ImageHeight)

	<-done2

	c.ioCommand <- ioOutput
	newFileName := filename + "x" + strconv.Itoa(response.Turns)
	c.ioFilename <- newFileName

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- response.World[y][x]
		}
	}
	c.events <- ImageOutputComplete{response.Turns, filename}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	c.events <- FinalTurnComplete{response.Turns,calculateAliveCells(p,response.World)}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
