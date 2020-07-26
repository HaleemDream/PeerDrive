package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	args "../args"
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

	if err != nil {
		log.Print(err)
		return
	}

	// TODO - structure of communication need to be changed
	fmt.Println("Succesfully connected to Server!")
	fmt.Println("Sending msg..")
	con.Write(sendPieceInformationRequest(file))
	header, peerPieces := recvPieceInformationRequest(con)

	fmt.Printf("File index = %d, pieceCount = %d\n", header.FileIndex, header.PieceCount)
	fmt.Println("Done receiving msg!")

	//if files.MissingPieces(file.Name) {
	con.Write(sendPieceRequest(file, peerPieces))
	//} else {
	//	fmt.Println("not missing any")
	//}

	recvPieces(con)

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

func recvPieces(con net.Conn) {
	fmt.Println("receiving piece payload")
	header := readHeader(con)
	readIndexPayload(con, header.PieceCount)
	piecePayload := readPiecePayload(con, header.PieceCount)

	fmt.Println(string(piecePayload))
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
		//if !files.HasPiece(file.Name, int(value)) {
		pieceCount++
		pieceRequestPayload.Write(uint32ToByteArr(value))
		//}
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

// temp
func uint32ToByteArr(value uint32) []byte {
	intBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(intBuffer, value)
	return intBuffer
}
