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
	header, payload := recvPieceInformationRequest(con)

	fmt.Printf("File index = %d, pieceCount = %d\n", header.FileIndex, header.PieceCount)
	fmt.Println("Payload...")
	fmt.Println(payload)
	fmt.Println("Done receiving msg!")

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
	buffer.Write(int32ToByteArr(file.ID))           // specifiy file
	buffer.Write(int32ToByteArr(0))                 // zero fill

	return buffer.Bytes()
}

// temp
func int32ToByteArr(value uint32) []byte {
	intBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(intBuffer, value)
	return intBuffer
}
