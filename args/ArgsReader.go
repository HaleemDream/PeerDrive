package args

// NetworkConfig - struct to hold network args
type NetworkConfig struct {
	ConnectionType string
	Port           string
	Host           string
}

func getDefaultNetworkConfig() NetworkConfig {
	var networkConfig NetworkConfig

	networkConfig.ConnectionType = "SERVER"
	networkConfig.Port = "20000"
	networkConfig.Host = "localhost"

	return networkConfig
}

// TODO - default peer drive dir, read in user select dir

// Read -[SERVER/CLIENT],[PORT], [HOSTNAME]
func Read(args []string) NetworkConfig {
	var defaultNetworkSettings = getDefaultNetworkConfig()

	if len(args) > 0 {
		defaultNetworkSettings.ConnectionType = args[0]
	}

	if len(args) > 1 {
		defaultNetworkSettings.Port = args[1]
	}

	if len(args) > 2 {
		defaultNetworkSettings.Host = args[2]
	}

	return defaultNetworkSettings
}
