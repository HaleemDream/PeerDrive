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

	con.Close()
}

func recvPieceInformationRequest(con net.Conn) (Header, []uint32) {
	var header Header

	if err := binary.Read(con, binary.BigEndian, &header.MessageType); err != nil {
		log.Print(err)
	}

	if err := binary.Read(con, binary.BigEndian, &header.FileIndex); err != nil {
		log.Print(err)
	}

	if err := binary.Read(con, binary.BigEndian, &header.PieceCount); err != nil {
		log.Print(err)
	}

	indexPayload := make([]uint32, header.PieceCount)

	for i := 0; i < int(header.PieceCount); i++ {
		if err := binary.Read(con, binary.BigEndian, &indexPayload[i]); err != nil {
			log.Print(err)
		}
	}

	return header, indexPayload
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
