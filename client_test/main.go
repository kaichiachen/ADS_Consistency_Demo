package main

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var mIndex = []string{"add", "new", "clear", "remove", "checkout"}
var m = map[string]string{
	"add":      "/additem",
	"new":      "/newitem",
	"clear":    "/clear",
	"remove":   "/removeitem",
	"checkout": "/checkout",
}

var client = []int{10000, 10001, 10002}

const opCount = 5

func main() {
	ops := []int{}
	for i := 0; i < opCount; i++ {
		r := rand.Intn(opCount)
		ops = append(ops, r)
	}
	var sum time.Duration = 0
	for _, v := range ops {
		switch v {
		case 0:
			port := client[rand.Intn(len(client))]
			ok, elapsed := newItemRequest("java", 100, 100, port)
			if !ok {
				fmt.Println(v, "error occur!")
			} else {
				sum += elapsed
			}
		case 1:
			port := client[rand.Intn(len(client))]
			items := itemsRequest(port)
			item := items[rand.Intn(len(items))]
			ok, elapsed := addItemRequest(item.ID, uint32(rand.Intn(int(item.Volume))), port)
			if !ok {
				fmt.Println(v, "error occur!")
			} else {
				sum += elapsed
			}
		case 2:
			port := client[rand.Intn(len(client))]
			ok, elapsed := clearCartRequest(port)
			if !ok {
				fmt.Println(v, "error occur!")
			} else {
				sum += elapsed
			}
		case 3:
			port := client[rand.Intn(len(client))]
			items := cartRequest(port)
			item := items[rand.Intn(len(items))]
			ok, elapsed := removeItemRequest(item.ID, rand.Intn(3), port)
			if !ok {
				fmt.Println(v, "error occur!")
			} else {
				sum += elapsed
			}
		case 4:
			port := client[rand.Intn(len(client))]
			ok, elapsed := checkoutRequest(port)
			if !ok {
				fmt.Println(v, "error occur!")
			} else {
				sum += elapsed
			}
		}
	}
}

func addItemRequest(id string, volume uint32, port int) (bool, time.Duration) {
	jsonStr := fmt.Sprintf(`{"id":"%s", "volume":%d}`, id, volume)
	var response common.Response
	resp, elapsed := request("POST", "/additem", []byte(jsonStr), true, port)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&response)
	return response.Succeed, elapsed
}

func newItemRequest(name string, volume uint32, price uint32, port int) (bool, time.Duration) {
	jsonStr := fmt.Sprintf(`{"name":"%s","price":%d ,"volume":%d}`, name, price, volume)
	var response common.Response
	resp, elapsed := request("POST", "/newitem", []byte(jsonStr), true, port)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&response)
	return response.Succeed, elapsed
}

func clearCartRequest(port int) (bool, time.Duration) {
	resp, elapsed := request("POST", "/clear", []byte{}, true, port)

	var response common.Response
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&response)
	return response.Succeed, elapsed
}

func removeItemRequest(id string, volume int, port int) (bool, time.Duration) {
	jsonStr := fmt.Sprintf(`{"id":"%s", "volume":%d}`, id, volume)

	var response common.Response
	resp, elapsed := request("POST", "/removeitem", []byte(jsonStr), true, port)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&response)
	return response.Succeed, elapsed
}

func checkoutRequest(port int) (bool, time.Duration) {
	resp, elapsed := request("POST", "/checkout", []byte{}, true, port)

	var response common.Response
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&response)
	return response.Succeed, elapsed
}

func itemsRequest(port int) []common.Item {
	resp, _ := request("GET", "/items", []byte{}, false, port)
	decoder := json.NewDecoder(resp.Body)
	var response common.Response
	decoder.Decode(&response)
	msgs := response.Msg.([]interface{})
	items := []common.Item{}
	for _, value := range msgs {
		v := value.(map[string]interface{})
		item := common.Item{v["name"].(string), uint32(v["volume"].(float64)), v["id"].(string), uint32(v["price"].(float64))}
		items = append(items, item)
	}
	return items
}

func cartRequest(port int) []common.Item {
	cartList := []common.Item{}
	resp, _ := request("GET", "/mycarts", []byte{}, false, port)
	decoder := json.NewDecoder(resp.Body)
	var response common.Response
	decoder.Decode(&response)
	msgs := response.Msg.([]interface{})
	for _, value := range msgs {
		v := value.(map[string]interface{})
		item := common.Item{v["name"].(string), uint32(v["volume"].(float64)), v["id"].(string), uint32(v["price"].(float64))}
		cartList = append(cartList, item)
	}
	return cartList
}

func request(method, api string, j []byte, benchmark bool, port int) (*http.Response, time.Duration) {
	var start time.Time
	if benchmark {
		start = time.Now()
	}
	url := fmt.Sprintf("localhost:%p%s", port, api)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if benchmark {
		elapsed := time.Since(start)
		fmt.Println()
		fmt.Printf("%c[%d;%d;%dm%s耗时: %s%c[0m ", 0x1B, 0, 40, 31, "", elapsed, 0x1B)
		fmt.Println()
		return resp, elapsed
	}
	return resp, 0
}
