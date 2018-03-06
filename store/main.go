package store

var Core = struct {
	*Network
}{}

func Start(address string, port int) {
	Core.Network = SetupNetwork(address, port)
	go Core.Network.Run()
}
