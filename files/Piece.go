package file

// TODO - need better solution to filename - being used a lot

import (
	"log"
	"math"

	network "../network"
)

const (
	hasPiece         = 1
	doesNotHavePiece = 0
)

const pieceFile = "placeholder"

var piecesByFilename map[string][]int

// InitializePieceInformation reads in file if present and init internal map
func InitializePieceInformation() {
	if fileExists(pieceFile) {
		// read in serialized list
		// save it to map
	} else {
		piecesByFilename = make(map[string][]int)
	}
}

// ReceivedPiece updates internal container with pieces client has
func ReceivedPiece(filename string, pieceIndex int) {
	if piecesByFilename[filename] == nil {
		piecesByFilename[filename] = make([]int, getPieceLength(filename))
	}

	piecesByFilename[filename][pieceIndex] = hasPiece
}

// HasPiece return if client has data for pieceIndex
func HasPiece(filename string, pieceIndex int) bool {
	if piecesByFilename[filename] == nil {
		return false
	}

	if pieceIndex >= len(piecesByFilename[filename]) {
		log.Print("piece index exceeds number of pieces")
		return false
	}

	return piecesByFilename[filename][pieceIndex] == hasPiece
}

func getPieceInformation(filename string) []int {
	// init array
	var piecesLength = getPieceLength(filename)

	if piecesByFilename[filename] == nil {
		piecesByFilename[filename] = make([]int, piecesLength)
	}

	if fileExists(filename) {
		for i := 0; i < piecesLength; i++ {
			piecesByFilename[filename][i] = hasPiece
		}
	}

	return piecesByFilename[filename]
}

// getPieceLength - number pieces
func getPieceLength(filename string) int {
	return int(math.Ceil(float64(fileSize(filename)) / float64(network.ChunkSize)))
}