package file

// TODO - need better solution to filename - being used a lot

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"math"
	"os"
)

// HavePiece
// DoesNotHavePiece
const (
	HavePiece        = 1
	DoesNotHavePiece = 0
)

const pieceFile = "peerdrive.data"

var piecesByFilename map[string][]int

// InitializePieceInformation reads in file if present and init internal map
func InitializePieceInformation() {
	if Exists(pieceFile) {
		piecesByFilename = ReadMap()
	} else {
		piecesByFilename = make(map[string][]int)
	}
}

// InitializeFilePieceInformation checks if file is present
func InitializeFilePieceInformation(filename string) {
	if piecesByFilename[filename] == nil {
		// init array
		var piecesLength = getPieceLength(filename)

		piecesByFilename[filename] = make([]int, piecesLength)

		// if file is present, then all pieces are present
		if Exists(filename) {
			for i := 0; i < piecesLength; i++ {
				piecesByFilename[filename][i] = HavePiece
			}
		}
		// else, the file is not present
		// possible that a in-progress file is present
		// attempt to read and capture what pieces we have
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

// SaveMap serialize map to disk
func SaveMap() {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(piecesByFilename); err != nil {
		log.Print(err)
	}

	if err := ioutil.WriteFile(pieceFile, buffer.Bytes(), 0644); err != nil {
		log.Print(err)
	}
}

// ReadMap deserializes map from disk
func ReadMap() map[string][]int {
	f, err := os.Open(pieceFile)
	defer f.Close()

	if err != nil {
		log.Print(err)
	}

	var decodedMap map[string][]int
	decoder := gob.NewDecoder(f)

	if err := decoder.Decode(&decodedMap); err != nil {
		log.Print(err)
	}

	return decodedMap
}

// GetPieceInformation returns pieces owned by client for file
func GetPieceInformation(filename string) []int {
	return piecesByFilename[filename]
}

// getPieceLength - number pieces
func getPieceLength(filename string) int {
	return int(math.Ceil(float64(fileSize(filename)) / float64(chunkSize)))
}
