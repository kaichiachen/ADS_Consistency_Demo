package main

import (
	"fmt"

	"github.com/stvp/tempredis"
)

func main() {
	for i := 0; i < 3; i++ {
                //server, err := tempredis.Start(tempredis.Config{"databases": "5", "port" : fmt.Sprintf("%d", 6666 + i)})
                server, err := tempredis.Start(tempredis.Config{"port" : fmt.Sprintf("%d", 6666 + i)})
                if err != nil {
                        panic(err)
                }
		fmt.Println(server.Socket())
        }
}
