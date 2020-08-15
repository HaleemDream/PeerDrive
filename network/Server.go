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
	fmt.Println("Connected to client...")

	defer client.Close()

	// TODO handle client reads better
	for {
		header := readHeader(client)

		switch header.MessageType {
		case RequestPieceInformation:
			sendPieceInformation(client, header)
		case RequestPieces:
			sendRequestedPieces(client, header)
		case TerminateConnection:
			return
		}
	}
}

func sendPieceInformation(client net.Conn, header Header) {
	fmt.Println("Receieved Request Piece Information...")
	file := meta.FileInformation(header.FileIndex)

	piecesOwned := files.GetPieceInformation(file.Name)

	var pieceCount uint32
	piecesIndexPayload := new(bytes.Buffer)
	for pieceIndex, piece := range piecesOwned {
		if piece == files.HavePiece {
			pieceCount++
			piecesIndexPayload.Write(uint32ToByteArr(uint32(pieceIndex)))
		}
	}

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(SendingPieceInformation)) // request piece information
	buffer.Write(uint32ToByteArr(file.ID))          // specifiy file id
	buffer.Write(uint32ToByteArr(pieceCount))       // specifiy number of pieces
	buffer.Write(piecesIndexPayload.Bytes())

	client.Write(buffer.Bytes())
}

func sendRequestedPieces(client net.Conn, header Header) {
	fmt.Println("received RequestPieces")

	file := meta.FileInformation(header.FileIndex)

	requestPieces := make([]uint32, header.PieceCount)
	payload := new(bytes.Buffer)

	for i := 0; i < int(header.PieceCount); i++ {
		if err := binary.Read(client, binary.BigEndian, &requestPieces[i]); err != nil {
			log.Print(err)
		}

		payload.Write(files.GetPiece(file.Name, requestPieces[i]))
	}

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(SendingPieces))
	buffer.Write(uint32ToByteArr(header.FileIndex))
	buffer.Write(uint32ToByteArr(header.PieceCount))
	binary.Write(buffer, binary.BigEndian, requestPieces)
	buffer.Write(payload.Bytes())

	client.Write(buffer.Bytes())
}
