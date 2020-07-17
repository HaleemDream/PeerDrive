package main

import (
	"log"
	"os"

	args "./args"
	network "./network"
)

func main() {
	networkSettings := args.Read(os.Args[1:])

	if networkSettings.ConnectionType == "SERVER" {
		network.Listen(networkSettings)
	} else if networkSettings.ConnectionType == "CLIENT" {
		network.Client(networkSettings)
	} else {
		log.Print("Invalid ConnectionType")
		return
	}
}
