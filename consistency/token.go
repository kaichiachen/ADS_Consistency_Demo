package consistency

import (
	// "log"
	"time"
)

var SERVER_PORTS = []int{20000, 20001}

var hasToken = false

func sendToken() {
	hasToken = true
	msg := NewMessage(MESSAGE_SEND_TOKEN)
	next_port := 0
	for k, v := range SERVER_PORTS {
		if v == Core.Network.Port {
			if k == len(SERVER_PORTS)-1 {
				next_port = SERVER_PORTS[0]
			} else {
				next_port = SERVER_PORTS[k+1]
			}
		}
	}
loop:
	for {
		select {
		case <-time.NewTimer(1 * time.Second).C:
			if Core.Network.SendMessage(*msg, next_port) {
				hasToken = false
			L:
				for {
					select {
					case <-Core.tokens:
					default:
						break L
					}
				}
				break loop
			} else {
				// log.Println("Resend Token...")
			}

		}
	}

}
