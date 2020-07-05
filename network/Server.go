package network

// Listen(args)
// onClient(client)

import (
	"encoding/binary"
	"log"
	"net"

	args "../args"
)

// Listen creates TCP server
// handles incoming messages
func Listen(settings args.NetworkConfig) {
	port := ":" + string(settings.Port)
	server, err := net.Listen("tcp4", port)

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
		// sendPieceInformation(client, header)
	case RequestPieces:
		// sendRequestedPieces(client, header)
	}

}
