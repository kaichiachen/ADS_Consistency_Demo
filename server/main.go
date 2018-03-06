package main

import (
	"common"
	"consistency"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"regexp"
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
	var addItem common.AddItem
	err := decoder.Decode(&addItem)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addItem)

}

func refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	consistency.RefreshShoppingCart()
	resp := common.Response{Succeed: true}
	jData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func clear(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	consistency.ClearShoppingCart()
}

func settle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	consistency.SettleShoppingCart()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	consistency.Start(getIPAddress(), comuport, nodes)

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/additem", additem)
	router.POST("/refrsh", refresh)
	router.POST("/settle", settle)
	router.POST("/clear", clear)
	fmt.Println(fmt.Sprintf("localhost:%d", restport))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", restport), router))
}

func getIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func findIPAddress(input string) string {
	validIpAddressRegex := "([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})"
	re := regexp.MustCompile(validIpAddressRegex)
	return re.FindString(input)
}
