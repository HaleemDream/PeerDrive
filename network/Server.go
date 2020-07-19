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

	files.InitializePieceInformation()

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
	fmt.Println("Receieved Request Piece Information...")
	file := meta.FileInformation(header.FileIndex)

	// TODO
	// init file data
	// should be done elsewhere
	fmt.Println("Init file piece information...")
	files.InitializeFilePieceInformation(file.Name)

	piecesOwned := files.GetPieceInformation(file.Name)

	var pieceCount uint32
	piecesIndexPayload := new(bytes.Buffer)
	for pieceIndex, piece := range piecesOwned {
		if piece == files.HavePiece {
			pieceCount++
			piecesIndexPayload.Write(int32ToByteArr(uint32(pieceIndex)))
		}
	}

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(SendingPieceInformation)) // request piece information
	buffer.Write(int32ToByteArr(file.ID))           // specifiy file id
	buffer.Write(int32ToByteArr(pieceCount))        // specifiy number of pieces
	buffer.Write(piecesIndexPayload.Bytes())

	fmt.Println("Sending bytes...")
	client.Write(buffer.Bytes())
	fmt.Println("Bytes receieved by client!")
}

func sendRequestedPieces(client net.Conn, header Header) {
	fmt.Println("received RequestPieces")
}
