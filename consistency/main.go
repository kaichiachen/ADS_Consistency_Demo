package consistency

import (
	"log"
	"fmt"
)

var Core = struct {
	*Network
	OperationSlice
	tokens chan string
}{}

var RedNum uint32 = uint32(0)
var SendRequest bool = false  // don't send reply unless you receive one

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
	case MESSAGE_START_UPDATE:
		fmt.Println("Receive MESSAGE_START_UPDATE")
		if Comuport != 20000 {
			break
		}
		startmsg := StartupMsg{}
		startmsg.UnMarshalInt(msg.Data)
		//fmt.Println("startmsg.RedNum = ", startmsg.RedNum, "RedNum", RedNum)
		if startmsg.RedNum < RedNum {
			ReplyCurStatus()
		} else {
			ReplyStatusIsNew()
		}
		break
	case MESSAGE_START_UPDATE_REPLY:
		fmt.Println("Receive MESSAGE_START_UPDATE_REPLY")
		fmt.Println("SendRequest is ", SendRequest)
		if SendRequest {
			fmt.Println("======== RedNum before ", RedNum)
			RedNum = UnMarshalCart(msg.Data)
			fmt.Println("======== RedNum after ", RedNum)
			SendRequest = false
		}
		break
	case MESSAGE_STATUS_IS_NEW:
		fmt.Println("Receive MESSAGE_STATUS_IS_NEW")
		fmt.Println("SendRequest is ", SendRequest)
		if SendRequest {
			SendRequest = false
		}
	}
}
