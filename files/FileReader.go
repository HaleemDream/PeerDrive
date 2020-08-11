package file

import (
	"bufio"
	"log"
	"os"
)

const chunkSize = 128

// ByteChunks - Reads a file and returns 2d array of bytes
// size of last array will be <= chunkSize
func ByteChunks(filename string) [][]byte {

	// file size in bytes
	var fileSize = fileSize(filename)
	var pieceCount = int(fileSize / int64(chunkSize))
	var lastByteArraySize = fileSize % int64(chunkSize)

	if lastByteArraySize != 0 {
		pieceCount++
	}

	var fileBytes = make([][]byte, pieceCount)
	var chunkSize uint32 = chunkSize

	for i := range fileBytes {
		// resize last array to be size of remaining bytes
		if lastByteArraySize != 0 && i == pieceCount-1 {
			chunkSize = uint32(lastByteArraySize)
		}

		fileBytes[i] = make([]byte, chunkSize)
	}

	f, err := os.Open(filename)
	reader := bufio.NewReader(f)

	defer f.Close()

	if err != nil {
		log.Print(err)
	}

	for i := 0; i < pieceCount; i++ {
		_, err := reader.Read(fileBytes[i])

		if err != nil {
			log.Print(err)
		}
	}

	return fileBytes
}

// GetPiece - returns bytes of pieceIndex
func GetPiece(filename string, pieceIndex uint32) []byte {
	// TODO - len of last piece is == ChunkSize but should be <=
	// TODO - err handle pieceIndex > len(pieces)

	bytes := make([]byte, chunkSize)
	f, err := os.Open(filename)

	defer f.Close()

	if err != nil {
		log.Print(err)
	}

	var offset int64 = int64(pieceIndex * chunkSize)

	f.ReadAt(bytes, offset)

	return bytes
}

// GetPieces - returns bytes of pieces
func GetPieces(filename string, pieceIndexes []uint32) [][]byte {
	bytes := make([][]byte, len(pieceIndexes)*chunkSize)

	for i, pieceIndex := range pieceIndexes {
		bytes[i] = GetPiece(filename, pieceIndex)
	}

	return bytes
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func fileSize(filename string) int64 {

	if !fileExists(filename) {
		return 0
	}

	fileInfo, err := os.Stat(filename)

	if err != nil {
		log.Print(err)
	}

	// file size in bytes
	return fileInfo.Size()
}
