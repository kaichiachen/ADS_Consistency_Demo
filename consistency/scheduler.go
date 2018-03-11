package consistency

import (
	"common"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var needBroadcast = false

func NewItem(item common.NewItem) chan common.Response {
	fmt.Println("新增商品")
	resp := make(chan common.Response)
	op := NewOperation(OP_ADDITEM)
	newItem := common.Item{ID: generateID(common.ITEM_ID_LENGTH), Name: item.Name, Volume: item.Volume, Price: item.Price}
	op.Payload, _ = newItem.MarshalBinary()
	op.PayloadLength = uint32(len(op.Payload))
	go execOpAndBroadcast(op, resp)
	return resp
}

func AddItemToCart(addeditem common.AddCartItem) chan common.Response {
	fmt.Println("放商品进购物车")
	resp := make(chan common.Response)
	op := NewOperation(OP_ADDCART)
	item := common.Item{ItemIDMap[addeditem.ID].Name, uint32(addeditem.Volume), addeditem.ID, ItemIDMap[addeditem.ID].Price}
	op.Payload, _ = item.MarshalBinary()
	op.PayloadLength = uint32(len(op.Payload))
	go execOpAndBroadcast(op, resp)
	return resp
}

func RemoveItemFromCart(rmitem common.RemoveCartItem) chan common.Response {
	fmt.Println("从购物车移除商品")
	resp := make(chan common.Response)
	op := NewOperation(OP_REMOVE)
	item := common.Item{ItemIDMap[rmitem.ID].Name, uint32(rmitem.Volume), rmitem.ID, ItemIDMap[rmitem.ID].Price}
	op.Payload, _ = item.MarshalBinary()
	op.PayloadLength = uint32(len(op.Payload))
	go execOpAndBroadcast(op, resp)
	return resp
}

func ClearShoppingCart() chan common.Response {
	fmt.Println("清除购物车")
	resp := make(chan common.Response)
	op := NewOperation(OP_CLEAR)
	go execOpAndBroadcast(op, resp)
	return resp
}

func CheckoutShoppingCart() chan common.Response {
	fmt.Println("结账")
	resp := make(chan common.Response, 1)
	op := NewOperation(OP_CHECKOUT)
	if !CheckItemVolume() {
		resp <- common.Response{Succeed: false}
	}
	op.Payload = ArchiveCartItems()
	op.PayloadLength = uint32(len(op.Payload))
	go execOpAndBroadcast(op, resp)
	return resp
}

func execOpAndBroadcast(op *Operation, resp chan common.Response) OP_RESULT {
	OpResult := op.generator()
	if OpResult == OPERATION_SUCCESS {
		Core.OperationSlice = Core.OperationSlice.AddOperation(op)
		red, blue := Core.OperationSlice.Count()
		fmt.Printf("%c[%d;%d;%dm%s %s: %d%c[0m ", 0x1B, 0, 40, 31, "", "red: ", red, 0x1B)
		fmt.Print("\t")
		fmt.Printf("%c[%d;%d;%dm%s %s: %d%c[0m ", 0x1B, 0, 40, 36, "", "blue: ", blue, 0x1B)
		fmt.Println()
		// if hasToken && op.Optype == RED {
		// 	broadcastOperations(resp)
		// } else if op.Optype == RED {
		// 	select {
		// 	case <-Core.tokens:
		// 		broadcastOperations(resp)
		// 		break
		// 	case <-time.NewTimer(5 * time.Second).C:
		// 		resp <- common.Response{Succeed: false}
		// 		break
		// 	}
		// } else {
		// 	resp <- common.Response{Succeed: true}
		// }
		if op.Optype == RED {
			mutex.Lock()
			broadcastOperations(resp)
			mutex.Unlock()
		} else {
			resp <- common.Response{Succeed: true}
		}
	} else {
		resp <- common.Response{Succeed: false}
	}
	return OpResult
}

func broadcastOperations(resp chan common.Response) {
	mes := NewMessage(MESSAGE_SEND_RED)
	data, err := Core.OperationSlice.MarshalBinary()
	if err != nil {
		log.Println(err)
	}
	mes.Data = data
	Core.OperationSlice = Core.OperationSlice.ClearOperation()
	time.Sleep(100 * time.Millisecond)
	Core.Network.BroadcastQueue <- *mes
	resp <- common.Response{Succeed: true}
	needBroadcast = false
}

func generateID(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
