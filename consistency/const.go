package consistency

const (
	RED = iota
	BLUE
)

const (
	OP_ADDITEM = iota + 2
	OP_ADDCART
	OP_REMOVE
	OP_CLEAR
	OP_CHECKOUT
)

const (
	MESSAGE_SEND_RED = iota + 20
	MESSAGE_SEND_TOKEN
)

const (
	OPERATION_SUCCESS = iota + 1
	OPERATION_FAIL
)
