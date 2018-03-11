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

var mIndex = []string{"新增商品", "放进购物车", "清除", "移除商品", "结账"}
var timeSlice = []time.Duration{0, 0, 0, 0, 0}
var countSlice = []int{0, 0, 0, 0, 0}
var averageSlice = []time.Duration{0, 0, 0, 0, 0}

// var m = map[string]string{
// 	"add":      "/additem",
// 	"new":      "/newitem",
// 	"clear":    "/clear",
// 	"remove":   "/removeitem",
// 	"checkout": "/checkout",
// }

var client = []int{10000, 10001, 10002}

const opCount = 100
const logError = true

func main() {
	ops := []int{}
	for i := 0; i < opCount; i++ {
		r := rand.Intn(5)
		ops = append(ops, r)
	}
loop:
	for _, v := range ops {
		fmt.Println(mIndex[v])
		switch v {
		case 0:
			port := client[rand.Intn(len(client))]
			ok, elapsed := newItemRequest(generateBook(5), 100, 100, port)
			if !ok && logError {
				fmt.Println(v, "error occur!")
			} else {
				timeSlice[v] += elapsed
				countSlice[v] += 1
			}
		case 1:
			port := client[rand.Intn(len(client))]
			items := itemsRequest(port)
			item := items[rand.Intn(len(items))]
			ok, elapsed := addItemRequest(item.ID, uint32(rand.Intn(int(item.Volume))), port)
			if !ok && logError {
				fmt.Println(v, "error occur!")
			} else {
				timeSlice[v] += elapsed
				countSlice[v] += 1
			}
		case 2:
			port := client[rand.Intn(len(client))]
			ok, elapsed := clearCartRequest(port)
			if !ok && logError {
				fmt.Println(v, "error occur!")
			} else {
				timeSlice[v] += elapsed
				countSlice[v] += 1
			}
		case 3:
			port := client[rand.Intn(len(client))]
			items := cartRequest(port)
			if len(items) != 0 {
				item := items[rand.Intn(len(items))]
				if item.Volume != 0 {
					num := rand.Intn(int(item.Volume))
					ok, elapsed := removeItemRequest(item.ID, num, port)
					if !ok && logError {
						fmt.Println(items)
						fmt.Println(item)
						fmt.Println(num)
						fmt.Println(v, "error occur!")
						fmt.Println(port)
						fmt.Println(cartRequest(port))
						break loop
					} else {
						timeSlice[v] += elapsed
						countSlice[v] += 1
					}
				}
			}
		case 4:
			port := client[rand.Intn(len(client))]
			ok, elapsed := checkoutRequest(port)
			if !ok && logError {
				fmt.Println(v, "error occur!")
			} else {
				timeSlice[v] += elapsed
				countSlice[v] += 1
			}
		}
	}
	for i := 0; i < len(timeSlice); i++ {
		if countSlice[i] != 0 {
			averageSlice[i] = timeSlice[i] / time.Duration(countSlice[i])
		}
	}
	for i := 0; i < len(averageSlice); i++ {
		if common.TypeMap[i] == common.RED {
			fmt.Printf("%c[%d;%d;%dm%s: %0.50s%c[0m ", 0x1B, 0, 40, 31, mIndex[i]+"latency", averageSlice[i], 0x1B)
			fmt.Println()
		} else {
			fmt.Printf("%c[%d;%d;%dm%s: %0.50s%c[0m ", 0x1B, 0, 40, 36, mIndex[i]+"latency", averageSlice[i], 0x1B)
			fmt.Println()
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
	url := fmt.Sprintf("http://localhost:%d%s", port, api)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if benchmark {
		elapsed := time.Since(start)
		return resp, elapsed
	}
	return resp, 0
}

func generateBook(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
