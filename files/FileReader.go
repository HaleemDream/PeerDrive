package file

import (
	"bufio"
	"log"
	"os"
)

// ByteChunks - Reads a file and returns 2d array of bytes
// size of last array will be <= chunkSize
func ByteChunks(fileName string, chunkSize int) [][]byte {
	fileInfo, err := os.Stat(fileName)

	if err != nil {
		log.Print(err)
	}

	// file size in bytes
	var fileSize = fileInfo.Size()
	var pieceCount = int(fileSize / int64(chunkSize))
	var lastByteArraySize = fileSize % int64(chunkSize)

	if lastByteArraySize != 0 {
		pieceCount++
	}

	var fileBytes = make([][]byte, pieceCount)

	for i := range fileBytes {
		// resize last array to be size of remaining bytes
		if lastByteArraySize != 0 && i == pieceCount-1 {
			chunkSize = int(lastByteArraySize)
		}

		fileBytes[i] = make([]byte, chunkSize)
	}

	f, err := os.Open(fileName)
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
