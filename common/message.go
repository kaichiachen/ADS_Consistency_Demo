package common

type AddCartItem struct {
	ID     string `json:"id"`
	Volume int    `json:"volume"`
}

type NewItem struct {
	Name   string `json:"name"`
	Volume uint32 `json:"volume"`
	Price  uint32 `json:"price"`
}

type RemoveCartItem struct {
	ID     string `json:"id"`
	Volume int    `json:"volume"`
}

type Response struct {
	Succeed bool        `json:"succeed"`
	Msg     interface{} `json:"msg"`
}

func NewResponse() chan Response {
	return make(chan Response)
}
