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

var port int

func init() {
	flag.IntVar(&port, "port", common.SERVER_DEFAULT_PORT, "server port")
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
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/additem", additem)
	router.POST("/refrsh", refresh)
	router.POST("/settle", settle)
	router.POST("/clear", clear)
	fmt.Println(fmt.Sprintf("localhost:%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
