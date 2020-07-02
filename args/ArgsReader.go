package args

import (
	"log"
	"strconv"
)

// NetworkConfig - struct to hold network args
type NetworkConfig struct {
	connectionType string
	port           int
	host           string
}

func getDefaultNetworkConfig() NetworkConfig {
	var networkConfig NetworkConfig

	networkConfig.connectionType = "SERVER"
	networkConfig.port = 20000
	networkConfig.host = "localhost"

	return networkConfig
}

// Read -[SERVER/CLIENT],[PORT], [HOSTNAME]
func Read(args []string) NetworkConfig {
	var defaultNetworkSettings = getDefaultNetworkConfig()

	if len(args) > 0 {
		defaultNetworkSettings.connectionType = args[0]
	}

	if len(args) > 1 {
		var err error
		defaultNetworkSettings.port, err = strconv.Atoi(args[1])

		if err != nil {
			log.Print(err)
		}
	}

	if len(args) > 2 {
		defaultNetworkSettings.host = args[2]
	}

	return defaultNetworkSettings
}
