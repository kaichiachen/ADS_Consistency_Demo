package consistency

var Core = struct {
	*Network
}{}

func Start(address string, port int, nodes []string) {
	Core.Network = SetupNetwork(address, port)
	go Core.Network.Run()
	for _, n := range nodes {
		Core.Network.ConnectionQueue <- n
	}

	go func() {
		for {
			select {
			case msg := <-Core.Network.IncomingMessages:
				HandleIncomingMessage(msg)
			}
		}
	}()
}

func HandleIncomingMessage(msg Message) {
	switch msg.Identifier {
	case MESSAGE_SEND_RED:
		op := new(Operation)
		op.UnMarshalBinary(msg.Data)
	case MESSAGE_SEND_TOKEN:
	case MESSAGE_SEND_REPLY:
	}
}
