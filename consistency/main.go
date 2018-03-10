package consistency

import (
	"log"
)

var Core = struct {
	*Network
	OperationSlice
}{}

func Start(address string, port int, nodes []string) {
	initData()
	Core.Network = SetupNetwork(address, port)
	go Core.Network.Run()
	if port == SERVER_PORTS[0] {
		hasToken = true
		go sendToken()
	}
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
		ops := OperationSlice{}
		ops.UnMarshalBinary(msg.Data)
		log.Println("Op sequence len: ", ops.Len())
		ops.HandleOperations()
	case MESSAGE_SEND_TOKEN:
		// log.Println("Receieve token")
		go sendToken()
	}
}
