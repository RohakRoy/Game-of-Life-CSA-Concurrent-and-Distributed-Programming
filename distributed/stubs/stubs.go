package stubs

var TurnsHandler = "GameOfLifeOperations.ProcessTurns"
var GetAliveCells = "GameOfLifeOperations.AliveCellsFinder"
var GetCurrentWorld = "GameOfLifeOperations.CurrentWorldFinder"
var QuitServer = "GameOfLifeOperations.EndServer"

type Response struct {
	World [][]byte
	Turns int
	AliveCells int
}

type Request struct {
	World [][]byte
	Turns       int
	ImageWidth  int
	ImageHeight int
}

type ReqWorld struct {}

type ResWorld struct {
	World [][]byte
	Turn int
}

type ReqQuitServer struct {}

type ResQuitServer struct {}