package consistency

import (
	"log"
)

var Core = struct {
	*Network
	OperationSlice
	tokens chan string
}{}

func Start(address string, port int, nodes []string) {
	initData()
	Core.tokens = make(chan string, 5)
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
		//log.Println(msg.Data)
		ops.UnMarshalBinary(msg.Data)

		log.Println("Op sequence len: ", ops.Len())
		// for _, s := range ops {
		// 	log.Println(s.PayloadLength)
		// }
		ops.HandleOperations()
	case MESSAGE_SEND_TOKEN:
		// log.Println("Receieve token")
		for i := 0; i < 5; i++ {
			Core.tokens <- ""
		}
		go sendToken()
	}
}
