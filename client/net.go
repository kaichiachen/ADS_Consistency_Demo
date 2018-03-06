package main

import (
	"common"
	"net"
)

type Network struct {
	ServerAddress string
	*net.TCPConn
}

func NewNetwork(address string) {

}

func (n *Network) RefreshShoppingCart() {
	common.NewMessage(common.REFRESH)
}

func (n *Network) ClearShoppingCart() {
	common.NewMessage(common.CLEAR)
}

func (n *Network) SettleShoppingCart() {
	m := common.NewMessage(common.SETTLE)
	b, _ := m.MarshalBinary()
	n.TCPConn.Write(b)
}

func clearShoppingCart() {

}
