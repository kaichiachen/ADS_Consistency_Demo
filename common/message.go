package common

type AddItem struct {
	Name string
}

type Response struct {
	Succeed bool
	Msg     interface{}
}

func NewResponse() chan Response {
	return make(chan Response)
}
