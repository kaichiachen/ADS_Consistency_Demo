package common

type AddCartItem struct {
	Name   string
	Volume int
}

type NewItem struct {
	Name   string
	Volume uint32
}

type Response struct {
	Succeed bool
	Msg     interface{}
}

func NewResponse() chan Response {
	return make(chan Response)
}
