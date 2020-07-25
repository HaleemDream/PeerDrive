package main

import (
	"log"
	"os"

	args "./args"
	files "./files"
	meta "./meta"
	network "./network"
)

func main() {
	networkSettings := args.Read(os.Args[1:])

	// init file map data
	files.InitializePieceInformation()

	// TODO - remove
	// init files
	swarmMetadata := meta.Retrieve()

	for _, file := range swarmMetadata.Files {
		files.InitializeFilePieceInformation(file.Name)
	}

	if networkSettings.ConnectionType == "SERVER" {
		network.Listen(networkSettings)
	} else if networkSettings.ConnectionType == "CLIENT" {
		network.Client(networkSettings)
	} else {
		log.Print("Invalid ConnectionType")
		return
	}
}
