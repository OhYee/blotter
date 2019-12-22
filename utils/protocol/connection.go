package proto

import (
	"time"
)

// Connection from a sub-server
type Connection struct {
	Description string    // Description of the connection
	APIList     []API     // APIList list of the sub-server api
	KeepAlive   time.Time // KeepAlive the last connect time of the sub-server
}

// API of the sub-server
type API struct {
	URL         string           // URL of the api
	Description string           // Description of the api
	Input       map[string]Value // Input arguments types
	Output      map[string]Value // Output data types
}

// Value of the input/output
type Value struct {
	Type        string // Type of the value
	Description string // Description of the value
}
