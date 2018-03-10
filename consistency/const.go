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

type OP_RESULT bool

const (
	OPERATION_SUCCESS = OP_RESULT(true)
	OPERATION_FAIL    = OP_RESULT(false)
)
