package consistency

func RefreshShoppingCart() {
	mes := NewMessage(MESSAGE_SEND_RED)
	Core.Network.BroadcastQueue <- *mes
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
