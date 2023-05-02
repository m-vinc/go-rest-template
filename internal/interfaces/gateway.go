package interfaces

import "github.com/gorilla/websocket"

type IGatewayListener func(client IGatewayClient, data map[string]string)

type IGateway interface {
	RegisterClient(conn *websocket.Conn)
	On(event string, listener IGatewayListener) error

	SendChannel(name string, msg []byte)
	Broadcast(msg []byte)
}

type IGatewayClient interface {
	Send(msg []byte)
}
