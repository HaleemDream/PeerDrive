package args

import (
	"log"
	"strconv"
)

// NetworkConfig - struct to hold network args
type NetworkConfig struct {
	ConnectionType string
	Port           int
	Host           string
}

func getDefaultNetworkConfig() NetworkConfig {
	var networkConfig NetworkConfig

	networkConfig.ConnectionType = "SERVER"
	networkConfig.Port = 20000
	networkConfig.Host = "localHost"

	return networkConfig
}

// TODO - default peer drive dir, read in user select dir

// Read -[SERVER/CLIENT],[Port], [HostNAME]
func Read(args []string) NetworkConfig {
	var defaultNetworkSettings = getDefaultNetworkConfig()

	if len(args) > 0 {
		defaultNetworkSettings.ConnectionType = args[0]
	}

	if len(args) > 1 {
		var err error
		defaultNetworkSettings.Port, err = strconv.Atoi(args[1])

		if err != nil {
			log.Print(err)
		}
	}

	if len(args) > 2 {
		defaultNetworkSettings.Host = args[2]
	}

	return defaultNetworkSettings
}
