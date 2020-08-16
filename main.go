package main

import (
	"log"
	"os"

	args "./args"
	files "./files"
	network "./network"
)

func main() {
	networkSettings := args.Read(os.Args[1:])

	// init file map data
	files.InitializePieceInformation()

	if networkSettings.ConnectionType == "SERVER" {
		network.Listen(networkSettings)
	} else if networkSettings.ConnectionType == "CLIENT" {
		network.Client(networkSettings)
	} else {
		log.Print("Invalid ConnectionType")
		return
	}
}
