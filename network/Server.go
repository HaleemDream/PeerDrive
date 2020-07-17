package network

// Listen(args)
// onClient(client)

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	args "../args"
	files "../files"
	meta "../meta"
)

// Listen creates TCP server
// handles incoming messages
func Listen(settings args.NetworkConfig) {
	fmt.Printf("Server Listening on %s\n", settings.Port)
	port := ":" + settings.Port
	server, err := net.Listen("tcp", port)

	if err != nil {
		log.Print(err)
		return
	}

	defer server.Close()

	for {
		client, err := server.Accept()

		if err != nil {
			log.Print(err)
			return
		}

		// TODO limit number of goroutines
		go onClient(client)
	}
}

func onClient(client net.Conn) {
	var header Header

	if err := binary.Read(client, binary.BigEndian, &header.MessageType); err != nil {
		log.Print(err)
	}

	if err := binary.Read(client, binary.BigEndian, &header.PieceCount); err != nil {
		log.Print(err)
	}

	if err := binary.Read(client, binary.BigEndian, &header.FileIndex); err != nil {
		log.Print(err)
	}

	switch header.MessageType {
	case RequestPieceInformation:
		sendPieceInformation(client, header)
	case RequestPieces:
		sendRequestedPieces(client, header)
	}
}

func sendPieceInformation(client net.Conn, header Header) {
	file := meta.FileInformation(header.FileIndex)
	piecesOwned := files.GetPieceInformation(file.Name)
	pieceLength := len(piecesOwned)

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(SendingPieceInformation)) // request piece information
	buffer.WriteByte(byte(file.ID))                 // specifiy file id
	buffer.WriteByte(byte(pieceLength))             // specifiy number of pieces

	for _, index := range piecesOwned {
		buffer.Write(int32ToByteArr(uint32(index)))
	}

	client.Write(buffer.Bytes())
}

func sendRequestedPieces(client net.Conn, header Header) {
	fmt.Println("received RequestPieces")
}
