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
	con.Write(sendPieceInformationRequest(file))
	con.Close()
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
