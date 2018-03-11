package common

const SERVER_REST_PORT = 10000
const SERVER_COMUNICATION_PORT = 10001

const (
	REFRESH = iota
	ADD
	REMOVE
	CLEAR
	SETTLE
)

const (
	RED = iota
	BLUE
)

var TypeMap = []int{
	RED,
	BLUE,
	BLUE,
	BLUE,
	RED,
}
