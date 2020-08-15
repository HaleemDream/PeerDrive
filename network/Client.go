package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	args "../args"
	files "../files"
	meta "../meta"
)

// Client attempts to connect to tcp server
func Client(settings args.NetworkConfig) {

	// metadata (request from server?)
	swarmMetadata := meta.Retrieve()

	// onClick (peer clicked on a file to download)
	// (create file)
	peer := swarmMetadata.Peers[0]
	file := swarmMetadata.Files[0]

	// download event triggerd :
	// do the following
	serverHostPort := fmt.Sprintf("%s:%s", peer.Host, peer.Port)
	con, err := net.Dial("tcp", serverHostPort)

	fmt.Println("Succesfully connected to Server!")

	if err != nil {
		log.Print(err)
		return
	}

	// init internal map
	if !files.Exists(file.Name) {
		files.InitializeFilePieceInformationExt(file.Name, file.Size)
	}

	// create file if not present
	f, err := os.OpenFile(file.Name, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		log.Print(err)
	}

	// request pieces we don't have
	for files.MissingPieces(file.Name) {
		// find out what pieces peer has
		// TODO - request periodically instead of constantly
		con.Write(sendPieceInformationRequest(file))
		_, peerPieces := recvPieceInformationRequest(con)

		// retrieve payload
		con.Write(sendPieceRequest(file, peerPieces))
		handleReceievedPieces(con, f)
	}

	// terminate connection
	con.Write(sendTerminationRequest())

	con.Close()
}

func readHeader(server net.Conn) Header {
	var header Header

	if err := binary.Read(server, binary.BigEndian, &header.MessageType); err != nil {
		log.Print(err)
	}

	if err := binary.Read(server, binary.BigEndian, &header.FileIndex); err != nil {
		log.Print(err)
	}

	if err := binary.Read(server, binary.BigEndian, &header.PieceCount); err != nil {
		log.Print(err)
	}

	return header
}

func readIndexPayload(server net.Conn, pieceLength uint32) []uint32 {
	indexPayload := make([]uint32, pieceLength)

	for i := 0; i < int(pieceLength); i++ {
		if err := binary.Read(server, binary.BigEndian, &indexPayload[i]); err != nil {
			log.Print(err)
		}
	}

	return indexPayload
}

func readPiecePayload(server net.Conn, pieceLength uint32) []byte {
	buffer := make([]byte, pieceLength*ChunkSize)

	server.Read(buffer)
	return buffer
}

func recvPieceInformationRequest(con net.Conn) (Header, []uint32) {
	header := readHeader(con)
	indexPayload := readIndexPayload(con, header.PieceCount)

	return header, indexPayload
}

func handleReceievedPieces(con net.Conn, f *os.File) {
	fmt.Println("receiving piece payload")
	header := readHeader(con)
	indexPayload := readIndexPayload(con, header.PieceCount)
	piecePayload := readPiecePayload(con, header.PieceCount)

	for index, i := range indexPayload {
		// maintain information on what pieces client now maintains
		files.ReceivedPiece(f.Name(), i)

		// write pieces
		// TODO - write chunk size pieces?
		_, err := f.WriteAt(piecePayload[i*ChunkSize:i*ChunkSize+ChunkSize], int64(index*ChunkSize))

		if err != nil {
			log.Print(err)
		}
	}
}

func sendPieceInformationRequest(file meta.File) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(RequestPieceInformation)) // request piece information
	buffer.Write(uint32ToByteArr(file.ID))          // specifiy file
	buffer.Write(uint32ToByteArr(0))                // zero fill

	return buffer.Bytes()
}

func sendPieceRequest(file meta.File, payload []uint32) []byte {
	var pieceCount uint32
	pieceRequestPayload := new(bytes.Buffer)

	for _, value := range payload {
		if !files.HasPiece(file.Name, int(value)) {
			pieceCount++
			pieceRequestPayload.Write(uint32ToByteArr(value))
		}
	}

	fmt.Printf("payload size = %d\n", pieceCount)
	fmt.Println(pieceRequestPayload)

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(RequestPieces))
	buffer.Write(uint32ToByteArr(file.ID))
	buffer.Write(uint32ToByteArr(pieceCount))
	buffer.Write(pieceRequestPayload.Bytes())

	return buffer.Bytes()
}

func sendTerminationRequest() []byte {
	fmt.Printf("Sending termination request\n")

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(TerminateConnection)) // request piece information
	buffer.Write(uint32ToByteArr(0))            // zero fill
	buffer.Write(uint32ToByteArr(0))            // zero fill

	return buffer.Bytes()
}

// temp
func uint32ToByteArr(value uint32) []byte {
	intBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(intBuffer, value)
	return intBuffer
}
