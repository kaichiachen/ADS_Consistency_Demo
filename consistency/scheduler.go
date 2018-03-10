package consistency

import (
	"common"
	"log"
	"math/rand"
)

var needBroadcast = false

func NewItem(item common.NewItem) chan common.Response {
	resp := common.NewResponse()
	op := NewOperation(OP_ADDITEM)
	newItem := common.Item{ID: generateID(common.ITEM_ID_LENGTH), Name: item.Name, Volume: item.Volume, Price: item.Price}
	op.Payload, _ = newItem.MarshalBinary()

	op.generator()
	Core.OperationSlice = Core.OperationSlice.AddOperation(op)

	if hasToken {
		broadcastOperations()
	} else {
		needBroadcast = true
	}

	return resp
}

func AddItemToCart(addeditem common.AddCartItem) (chan common.Response, bool) {
	resp := common.NewResponse()
	op := NewOperation(OP_ADDCART)
	item := common.Item{ItemIDMap[addeditem.ID].Name, uint32(addeditem.Volume), addeditem.ID, ItemIDMap[addeditem.ID].Price}
	op.Payload, _ = item.MarshalBinary()

	OpResult := op.generator()
	Core.OperationSlice = Core.OperationSlice.AddOperation(op)

	return resp, OpResult
}

func RemoveItemFromCart(rmitem common.RemoveCartItem) (chan common.Response, bool) {
	resp := common.NewResponse()
	op := NewOperation(OP_REMOVE)
	item := common.Item{ItemIDMap[rmitem.ID].Name, uint32(rmitem.Volume), rmitem.ID, ItemIDMap[rmitem.ID].Price}
	op.Payload, _ = item.MarshalBinary()

	OpResult := op.generator()
	Core.OperationSlice = Core.OperationSlice.AddOperation(op)

	return resp, OpResult
}

func ClearShoppingCart() (chan common.Response, bool) {
	resp := common.NewResponse()
	op := NewOperation(OP_CLEAR)

	OpResult := op.generator()
	Core.OperationSlice = Core.OperationSlice.AddOperation(op)

	return resp, OpResult
}

func CheckoutShoppingCart() {
	mes := NewMessage(MESSAGE_SEND_RED)
	Core.Network.BroadcastQueue <- *mes
}

func broadcastOperations() {
	mes := NewMessage(MESSAGE_SEND_RED)
	data, err := Core.OperationSlice.MarshalBinary()
	if err != nil {
		log.Println(err)
	}
	mes.Data = data
	Core.OperationSlice = Core.OperationSlice.ClearOperation()
	Core.Network.BroadcastQueue <- *mes
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
