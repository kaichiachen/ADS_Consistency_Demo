package common

type AddCartItem struct {
	Name   string
	Volume int
}

type NewItem struct {
	Name   string
	Volume uint32
	Price  uint32
}

type RemoveCartItem struct {
	ID string
}

type Response struct {
	Succeed bool
	Msg     interface{}
}

func NewResponse() chan Response {
	return make(chan Response)
}
