package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var port int

func init() {
	flag.IntVar(&port, "port", 33, "server port")
}

func usage() {
	flag.PrintDefaults()
}

var url string

func main() {
	url = fmt.Sprintf("http://localhost:%d", port)
	fmt.Println("欢迎光临线上购物书城")
	for {
		fmt.Println()
		fmt.Println("购物车:")
		fmt.Println()
		fmt.Println("请问您要：")
		fmt.Println("1. 新增商品")   // red
		fmt.Println("2. 看商品")    // blue
		fmt.Println("3. 移除某项商品") // blue
		fmt.Println("4. 清空购物车")  // blue
		fmt.Println("5. 结算购物车")  //red
		input := <-readStdin()
		switch input {
		case "1":
			var jsonStr = []byte(`{"name":"book"}`)
			resp := request("POST", "/refresh", jsonStr)
			fmt.Println(resp)
		case "2":
			var jsonStr = []byte(`{"name":"book"}`)
			resp := request("POST", "/additem", jsonStr)
			fmt.Println(resp)
		case "3":
			fmt.Println("你想移除哪项商品商品?")
			input := <-readStdin()
			fmt.Println(input)
		case "4":
			// store.ClearShoppingCart()
		case "5":
			// store.SettleShoppingCart()
		default:
			break
		}
	}
}

func request(method, api string, j []byte) *http.Response {
	req, err := http.NewRequest(method, url+api, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func readStdin() chan string {
	cb := make(chan string)
	input := bufio.NewScanner(os.Stdin)
	go func() {
		if input.Scan() {
			cb <- input.Text()
		}
	}()
	return cb
}
