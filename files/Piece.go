package file

// TODO - need better solution to filename - being used a lot

import (
	"log"
	"math"
)

// HavePiece
// DoesNotHavePiece
const (
	HavePiece        = 1
	DoesNotHavePiece = 0
)

const pieceFile = "placeholder"

var piecesByFilename map[string][]int

// InitializePieceInformation reads in file if present and init internal map
func InitializePieceInformation() {
	if Exists(pieceFile) {
		// read in serialized list
		// save it to map
	} else {
		piecesByFilename = make(map[string][]int)
	}
}

// InitializeFilePieceInformation checks if file is present
func InitializeFilePieceInformation(filename string) {
	// init array
	var piecesLength = getPieceLength(filename)

	if piecesByFilename[filename] == nil {
		piecesByFilename[filename] = make([]int, piecesLength)
	}

	if Exists(filename) {
		for i := 0; i < piecesLength; i++ {
			piecesByFilename[filename][i] = HavePiece
		}
	}
}

// InitializeFilePieceInformationExt - Test
func InitializeFilePieceInformationExt(filename string, size uint64) {
	// init array
	var piecesLength = int(math.Ceil(float64(size) / float64(chunkSize)))

	//if piecesByFilename[filename] == nil {
	piecesByFilename[filename] = make([]int, piecesLength)
	//}

	if Exists(filename) {
		for i := 0; i < piecesLength; i++ {
			piecesByFilename[filename][i] = HavePiece
		}
	}
}

// ReceivedPiece updates internal container with pieces client has
func ReceivedPiece(filename string, pieceIndex uint32) {
	if piecesByFilename[filename] == nil {
		piecesByFilename[filename] = make([]int, getPieceLength(filename))
	}

	piecesByFilename[filename][pieceIndex] = HavePiece
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

	return piecesByFilename[filename][pieceIndex] == HavePiece
}

// MissingPieces returns true if not all pieces present
func MissingPieces(filename string) bool {
	if piecesByFilename[filename] == nil {
		return true
	}

	for _, value := range piecesByFilename[filename] {
		if value != HavePiece {
			return true
		}
	}

	return false
}

// GetPieceInformation returns pieces owned by client for file
func GetPieceInformation(filename string) []int {
	return piecesByFilename[filename]
}

// getPieceLength - number pieces
func getPieceLength(filename string) int {
	return int(math.Ceil(float64(fileSize(filename)) / float64(chunkSize)))
}
