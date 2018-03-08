package main

import (
	"bufio"
	"bytes"
	"common"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	url = fmt.Sprintf("http://localhost:%d", port)

	initData()

	fmt.Println("欢迎光临线上购物书城")
	for {
		fmt.Println()
		fmt.Println("购物车:")
		printCartList()
		fmt.Println()
		fmt.Println("请问您要：")
		fmt.Println("1. 新增商品")   // red
		fmt.Println("2. 看商品")    // blue
		fmt.Println("3. 移除某项商品") // blue
		fmt.Println("4. 清空购物车")  // blue
		fmt.Println("5. 结算购物车")  //red
		fmt.Print("请选择：")
		input := <-readStdin()
		switch input {
		case "1":
			fmt.Println()
			fmt.Print("商品名称: ")
			name := <-readStdin()
			fmt.Print("商品售价: ")
			input = <-readStdin()
			price, _ := strconv.Atoi(input)
			fmt.Print("商品数量: ")
			input = <-readStdin()
			volume, _ := strconv.Atoi(input)

			jsonStr := fmt.Sprintf(`{"name":"%s","price":%d ,"volume":%d}`, name, price, volume)
			resp := request("POST", "/newitem", []byte(jsonStr))
			decoder := json.NewDecoder(resp.Body)
			var response common.Response
			err := decoder.Decode(&response)
			if err != nil {
				fmt.Println(err)
			}
			if response.Succeed {
				fmt.Println("新增商品成功！")
			}
		case "2":
			resp := request("GET", "/items", []byte{})
			decoder := json.NewDecoder(resp.Body)
			var response common.Response
			err := decoder.Decode(&response)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println()
			msgs := response.Msg.([]interface{})
			items := []common.Item{}
			for k, value := range msgs {
				v := value.(map[string]interface{})
				item := common.Item{v["name"].(string), uint32(v["volume"].(float64)), v["id"].(string), uint32(v["price"].(float64))}
				fmt.Print(k + 1)
				fmt.Print(". ", item.Name, " - ")
				fmt.Print(item.Price)
				fmt.Print("元 - ")
				fmt.Println("剩余", item.Volume, "个")
				items = append(items, item)
			}
			fmt.Print("请选择：")
			input = <-readStdin()
			index, _ := strconv.Atoi(input)
			num := -1
			for uint32(num) > items[index-1].Volume || num == -1 {
				fmt.Print("请选择数量：")
				input = <-readStdin()
				num, _ = strconv.Atoi(input)
			}

			jsonStr := fmt.Sprintf(`{"id":"%s", "volume":%d}`, items[index-1].ID, num)
			resp = request("POST", "/additem", []byte(jsonStr))
			decoder = json.NewDecoder(resp.Body)
			err = decoder.Decode(&response)
			if err != nil {
				fmt.Println(err)
			}
			if response.Succeed {
				fmt.Println("加入购物车成功！")
			}
		case "3":
			// id(string(10))
			fmt.Println("你想移除哪项商品商品?")
			input := <-readStdin()
			fmt.Println(input)
			var jsonStr = []byte(`{"id":"l3k4l1n3x1m3"`)
			resp := request("POST", "/removeitem", jsonStr)
			fmt.Println(resp)
		case "4":
			// id(string,volume(string)
			var jsonStr = []byte(`{"id":"l3k4l1n3x1m3,34dsd214dsd,23fsdfd123", "volume":"2,3,1"}`)
			resp := request("POST", "/clear", jsonStr)
			fmt.Println(resp)
		case "5":
			// id(string,volume(string)
			var jsonStr = []byte(`{"id":"l3k4l1n3x1m3,34dsd214dsd,23fsdfd123", "volume":"2,3,1"}`)
			resp := request("POST", "/checkout", jsonStr)
			fmt.Println(resp)
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

func printCartList() {
	resp := request("GET", "/mycarts", []byte{})
	decoder := json.NewDecoder(resp.Body)
	var response common.Response
	err := decoder.Decode(&response)
	if err != nil {
		fmt.Println(err)
	}
	msgs := response.Msg.([]interface{})
	sum := 0
	for k, value := range msgs {
		v := value.(map[string]interface{})
		item := common.Item{v["name"].(string), uint32(v["volume"].(float64)), v["id"].(string), uint32(v["price"].(float64))}
		fmt.Print(k + 1)
		fmt.Print(". ", item.Name, " - ")
		fmt.Print(item.Volume)
		fmt.Print(" X ")
		fmt.Println(item.Price, "=", item.Volume*item.Price, "元")
		sum += int(item.Volume * item.Price)
	}
	fmt.Println("总共: ", sum, "元")
}
