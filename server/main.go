package main

import (
	"common"
	"consistency"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type arrayFlags []string

var restport int
var comuport int
var nodes arrayFlags

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.IntVar(&restport, "restport", common.SERVER_REST_PORT, "server restful port")
	flag.IntVar(&comuport, "comuport", common.SERVER_COMUNICATION_PORT, "server  communication port")
	flag.Var(&nodes, "addr", "Other Server Address")
}

func usage() {
	flag.PrintDefaults()
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func additem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var addCartItem common.AddCartItem
	err := decoder.Decode(&addCartItem)
	if err != nil {
		fmt.Println(err)
	}

	_, OpResult := consistency.AddItemToCart(addCartItem)

	resp := common.Response{Succeed: true}
	if !OpResult {
		resp = common.Response{Succeed: false}
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func newitem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var newItem common.NewItem
	err := decoder.Decode(&newItem)
	if err != nil {
		fmt.Println(err)
	}

	consistency.NewItem(newItem)

	resp := common.Response{Succeed: true}
	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func removeCartItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var removeItem common.RemoveCartItem
	err := decoder.Decode(&removeItem)
	if err != nil {
		fmt.Println(err)
	}

	_, OpResult := consistency.RemoveItemFromCart(removeItem)

	resp := common.Response{Succeed: true}
	if !OpResult {
		resp = common.Response{Succeed: false}
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func clear(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, OpResult := consistency.ClearShoppingCart()

	resp := common.Response{Succeed: true}
	if !OpResult {
		resp = common.Response{Succeed: false}
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	consistency.CheckoutShoppingCart()

	resp := common.Response{Succeed: true}

	checkoutSuccess := consistency.CheckoutForServer()
	if 0 == checkoutSuccess {
		resp = common.Response{Succeed: false}
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func items(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := common.Response{Succeed: true}
	resp.Msg = consistency.GetClientItemIDMap()
	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func mycarts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := common.Response{Succeed: true}
	resp.Msg = consistency.GetItemIDMapFromCart()
	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	consistency.MutexInit()
	consistency.Start(getIPAddress(), comuport, nodes)

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/additem", additem)
	router.POST("/newitem", newitem)
	router.POST("/removeitem", removeCartItem)
	router.POST("/checkout", checkout)
	router.POST("/clear", clear)
	router.GET("/mycarts", mycarts)
	router.GET("/items", items)
	fmt.Println(fmt.Sprintf("localhost:%d", restport))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", restport), router))
}

func getIPAddress() string {
	// conn, err := net.Dial("udp", "8.8.8.8:80")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()

	// localAddr := conn.LocalAddr().(*net.UDPAddr)

	// return localAddr.IP.String()
	return "localhost"
}
