package netutils

import "encoding/json"

// Node holds a single UUID -> Host touple
type Node struct {
	UUID string
	Host string
	Port int
	IP   string
}

// JSON returns the json representation of a Node
func (n *Node) JSON() []byte {
	b, _ := json.Marshal(n)
	return b
}
