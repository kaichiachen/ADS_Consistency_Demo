package consistency

import "common"

func RefreshShoppingCart() chan common.Response {
	resp := common.NewResponse()
	mes := NewMessage(MESSAGE_SEND_RED)
	Core.Network.BroadcastQueue <- *mes

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
