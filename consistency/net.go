package consistency

import (
	"fmt"
	"net"
)

type NodeChannel chan *Node
type Node struct {
	*net.TCPConn
	lastSeen int
}
type Nodes map[string]*Node
type Network struct {
	Nodes
	Address            string
	ConnectionCallback NodeChannel
	BroadcastQueue     chan Message
	IncomingMessages   chan Message
}

func SetupNetwork(address string, port int) *Network {
	n := &Network{}
	n.BroadcastQueue, n.IncomingMessages = make(chan Message), make(chan Message)
	n.Nodes = Nodes{}
	n.Address = fmt.Sprintf("%s:%d", address, port)
	return n
}

func (n *Network) Run() {
	for {
		select {
		case message := <-n.BroadcastQueue:
			go n.BroadcastMessage(message)
		}
	}
}

func (n *Network) BroadcastMessage(message Message) {
	b, _ := message.MarshalBinary()
	for k, node := range n.Nodes {
		fmt.Println("Broadcasting...", k)
		go func() {
			_, err := node.TCPConn.Write(b)
			if err != nil {
				fmt.Println("Error bcing to", node.TCPConn.RemoteAddr())
			}
		}()
	}
}
