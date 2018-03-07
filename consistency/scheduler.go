package consistency

import (
	"common"
	"log"
	"math/rand"
)

func NewItem(item common.NewItem) chan common.Response {
	resp := common.NewResponse()
	mes := NewMessage(MESSAGE_SEND_RED)
	op := NewOperation(OP_ADDITEM)
	newItem := common.Item{ID: generateID(common.ITEM_ID_LENGTH), Name: item.Name, Volume: item.Volume}

	op.Payload, _ = newItem.MarshalBinary()
	Core.OperationSlice = Core.OperationSlice.AddOperation(op)

	data, err := Core.OperationSlice.MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	mes.Data = data
	Core.OperationSlice.ClearOperation()
	Core.Network.BroadcastQueue <- *mes

	return resp
}

func AddItemToCart(item common.AddCartItem) chan common.Response {
	resp := common.NewResponse()
	op := NewOperation(OP_ADDITEM)
	Core.OperationSlice.AddOperation(op)

	return resp
}

func ClearShoppingCart() {
	mes := NewMessage(MESSAGE_SEND_RED)
	Core.Network.BroadcastQueue <- *mes
}

func SettleShoppingCart() {
	ClearShoppingCart()
	mes := NewMessage(MESSAGE_SEND_RED)
	Core.Network.BroadcastQueue <- *mes
}

func generateID(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
