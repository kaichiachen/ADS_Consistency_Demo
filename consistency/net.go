package consistency

import (
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"
)

type NodeChannel chan *Node
type ConnectionQueue chan string
type Node struct {
	*net.TCPConn
	lastSeen  int
	connected bool
}
type Nodes map[string]*Node
type Network struct {
	Nodes
	ConnectionQueue
	Address            string
	Port               int
	ConnectionCallback NodeChannel
	BroadcastQueue     chan Message
	IncomingMessages   chan Message
	StartupMessageQueue chan Message
	StartupReplyQueue chan Message
}

func SetupNetwork(address string, port int) *Network {
	n := &Network{}
	n.BroadcastQueue, n.IncomingMessages = make(chan Message), make(chan Message)
	n.ConnectionQueue, n.ConnectionCallback = CreateConnectionQueue()
	n.StartupMessageQueue = make(chan Message)
	n.StartupReplyQueue = make(chan Message)
	n.Nodes = Nodes{}
	n.Address = fmt.Sprintf("%s:%d", address, port)
	n.Port = port
	return n
}

func (n *Network) Run() {
	log.Println("Listening in", Core.Network.Address)
	listenCb := StartListening(Core.Network.Address)
	for {
		select {
		case node := <-listenCb:
			Core.Nodes.AddNode(node)
		case node := <-n.ConnectionCallback:
			Core.Nodes.AddNode(node)
		case message := <-n.BroadcastQueue:
			go n.BroadcastMessage(message)
		case startmsg := <-n.StartupMessageQueue:
			go n.StartUpdateStatus(startmsg)
		case reply := <-n.StartupReplyQueue:
			go n.StartReplyStatus(reply)
		}
	}
}

func StartListening(address string) NodeChannel {

	cb := make(NodeChannel)
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	go func(l *net.TCPListener) {

		for {
			connection, err := l.AcceptTCP()
			if err != nil && err != io.EOF {
				log.Println(err)
			}

			cb <- &Node{connection, int(time.Now().Unix()), true}
		}
	}(listener)

	return cb
}

func (n Nodes) AddNode(node *Node) bool {
	ip, p, _ := net.SplitHostPort(node.TCPConn.RemoteAddr().String())
	port, _ := strconv.Atoi(p)
	addr := fmt.Sprintf("%s:%d", ip, port)

	if addr != Core.Network.Address { // allow restarted server to overwrite
		// log.Println("Node connected ", addr)
		n[addr] = node

		go HandleNode(node)

		return true
	}

	log.Println("Duplicate ip address")
	return false
}

func HandleNode(node *Node) {
	for {
		var bs []byte = make([]byte, 1024*1000)
		n, err := node.TCPConn.Read(bs)
		if err == io.EOF {
			log.Printf("%s： Connection Closed\n", node.TCPConn.RemoteAddr().String())
			node.TCPConn.Close()
			break
		}
		m := new(Message)
		err = m.UnMarshalBinary(bs[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}

		m.Reply = make(chan Message)

		go func(cb chan Message) {
			for {
				m, ok := <-cb

				if !ok {
					close(cb)
					break
				}

				b, _ := m.MarshalBinary()
				l := len(b)

				i := 0
				for i < l {
					a, _ := node.TCPConn.Write(b[i:])
					i += a
				}
			}
		}(m.Reply)

		Core.Network.IncomingMessages <- *m
	}
}

func CreateConnectionQueue() (ConnectionQueue, NodeChannel) {
	in := make(ConnectionQueue)
	out := make(NodeChannel)

	go func() {

		for {
			address := <-in
			if address != Core.Network.Address && Core.Nodes[address] == nil {
				log.Printf("Connecting to node: %s\n", address)
				go ConnectToNode(address, 5*time.Second, true, out)
			}
		}
	}()

	return in, out
}

func ConnectToNode(dst string, timeout time.Duration, retry bool, cb NodeChannel) {

	addrDst, err := net.ResolveTCPAddr("tcp4", dst)

	if err != nil && err != io.EOF {
		log.Println(err)
	}

	var con *net.TCPConn = nil
loop:
	for {
		breakChannel := make(chan bool)
		go func() {

			con, err = net.DialTCP("tcp", nil, addrDst)

			if con != nil {
				log.Println("Node connected ", addrDst.String())
				cb <- &Node{con, int(time.Now().Unix()), true}
				breakChannel <- true
			}
		}()

		select {
		case <-time.NewTimer(timeout).C:
			if !retry {
				break loop
			}
		case <-breakChannel:
			break loop
		}
	}
}

func (n *Network) StartUpdateStatus(message Message) {
	b, _ := message.MarshalBinary()
	for _, node := range n.Nodes {
		rip := node.TCPConn.RemoteAddr().String()
		p := rip[len(findIPAddress(rip))+1:]
		port, _ := strconv.Atoi(p);
		if port == 20000 { // 20000 is the seed
			go func(node Node) {
				_, err := node.TCPConn.Write(b)
				if err != nil {
					fmt.Println("start message send failed to ", node.TCPConn.RemoteAddr())
				} else {
					fmt.Println("start message send success!!")
				}
			}(*node)
		}
	}
}

func (n *Network) BroadcastMessage(message Message) {
	b, _ := message.MarshalBinary()
	for _, node := range n.Nodes {
		rip := node.TCPConn.RemoteAddr().String()
		p := rip[len(findIPAddress(rip))+1:]
		port, _ := strconv.Atoi(p)
		if port >= 20000 && port <= 20002 {
			fmt.Println("Broadcasting...", rip)
			go func(node Node) {
				_, err := node.TCPConn.Write(b)
				if err != nil {
					fmt.Println("Error broadcast to", node.TCPConn.RemoteAddr())
				} else {
					fmt.Println("start message send success!!")
				}
			}(*node)
		}
	}
}

func (n *Network) StartReplyStatus(message Message) {
        b, _ := message.MarshalBinary()
        for _, node := range n.Nodes {
                rip := node.TCPConn.RemoteAddr().String()
                p := rip[len(findIPAddress(rip))+1:]
                port, _ := strconv.Atoi(p);
                fmt.Println("p = ", p, " port = ", port)
                if Comuport == 20000 { // 20000 is the seed
                        fmt.Println("Going to reply to ", node.TCPConn.RemoteAddr())
                        go func(node Node) {
                                _, err := node.TCPConn.Write(b)
                                if err != nil {
                                        fmt.Println("start message reply failed to ", node.TCPConn.RemoteAddr())
                                } else {
                                        fmt.Println("start message reply success!!")
                                }
                        }(*node)
                }
        }
}

func (n *Network) SendMessage(message Message, sendPort int) bool {
	b, _ := message.MarshalBinary()
	send := false
	for k, node := range n.Nodes {
		p := k[len(findIPAddress(k))+1:]
		port, _ := strconv.Atoi(p)
		if port == sendPort {
			// fmt.Println("Send Token...", k)
			go func() {
				_, err := node.TCPConn.Write(b)
				if err != nil {
					fmt.Println("Error broadcast to", node.TCPConn.RemoteAddr())
				}
			}()
			send = true
			break
		}
	}
	return send
}

func findIPAddress(input string) string {
	validIpAddressRegex := "([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})"
	re := regexp.MustCompile(validIpAddressRegex)
	return re.FindString(input)
}
