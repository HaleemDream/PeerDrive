package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	args "../args"
)

// Client attempts to connect to tcp server
func Client(settings args.NetworkConfig) {
	serverHostPort := fmt.Sprintf("%s:%s", settings.Host, settings.Port)
	con, err := net.Dial("tcp", serverHostPort)

	if err != nil {
		log.Print(err)
		return
	}

	fmt.Println("Succesfully connected to Server!")
	con.Write(sendPieceInformationRequest())
	con.Close()
}

func sendPieceInformationRequest() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(RequestPieceInformation))
	buffer.Write(int32ToByteArr(0))
	buffer.Write(int32ToByteArr(0))

	return buffer.Bytes()
}

// temp
func int32ToByteArr(value int32) []byte {
	intBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32intBuffer, uint32(value))
	return intBuffer
}
