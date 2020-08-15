package meta

// File contains information on file
type File struct {
	Name string
	Size uint64 // bytes
	ID   uint32
}

// Tracker data represents a peer drive metadata
type Tracker struct {
	Peers []Peer
	Files []File
}

// Peer represents hosts in the swarm
type Peer struct {
	Host string
	Port string
}

// Retrieve place holder request from server about a specific peer drive
func Retrieve( /*swarm hash*/ ) Tracker {
	var file File
	file.Name = "test.txt"
	file.Size = 1025
	file.ID = 0

	var peer Peer
	peer.Host = "localhost"
	peer.Port = "20000"

	var tracker Tracker
	tracker.Peers = []Peer{peer}
	tracker.Files = []File{file}

	return tracker
}

// FileInformation returns information based on file id
func FileInformation(fileID uint32) File {
	// need to figure out how often to query metadata
	swarmMetadata := Retrieve()

	for _, file := range swarmMetadata.Files {
		if file.ID == fileID {
			return file
		}
	}

	return File{}
}
