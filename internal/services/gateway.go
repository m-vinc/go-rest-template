package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mpj/internal/interfaces"

	"github.com/gorilla/websocket"
	"go.uber.org/zap/zapcore"
)

type Gateway struct {
	logger interfaces.ILoggerService

	pongWait   time.Duration
	writeWait  time.Duration
	pingPeriod time.Duration
	maxSize    int

	clients   map[*Client]bool
	channels  map[string]map[*Client]bool
	listeners map[string][]interfaces.IGatewayListener
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func (client *Client) Send(msg []byte) {
	client.send <- msg
}

func NewGatewayService(logger interfaces.ILoggerService) *Gateway {
	return &Gateway{
		logger: logger,

		pongWait:   time.Second * 2,
		pingPeriod: 1 * time.Second,
		writeWait:  10 * time.Second,
		maxSize:    10 * 1000 * 1000, // 10MB
		clients:    map[*Client]bool{},
		channels:   map[string]map[*Client]bool{},
		listeners:  map[string][]interfaces.IGatewayListener{},
	}
}

func (gw *Gateway) joinChannel(name string, client *Client) {
	_, ok := gw.channels[name]
	if !ok {
		gw.channels[name] = map[*Client]bool{}
	}

	gw.channels[name][client] = true
	for _, listener := range gw.listeners["channel:join"] {
		listener(client, map[string]string{
			"channel": name,
		})
	}
}

func (gw *Gateway) leaveChannel(name string, client *Client) {
	_, ok := gw.channels[name]
	if !ok {
		return
	}

	if _, ok := gw.channels[name]; ok {
		delete(gw.channels[name], client)
		for _, listener := range gw.listeners["channel:leave"] {
			listener(client, map[string]string{
				"channel": name,
			})
		}
	}
}

type Message struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`

	Body json.RawMessage
}

func (gw *Gateway) RegisterClient(conn *websocket.Conn) {
	client := &Client{conn: conn, send: make(chan []byte)}

	gw.clients[client] = true
	go gw.readPump(client)
	go gw.writePump(client)
}

func (gw *Gateway) unregisterClient(client *Client) {
	if _, ok := gw.clients[client]; ok {
		delete(gw.clients, client)
		close(client.send)
	}

	for chanName, clients := range gw.channels {
		delete(clients, client)
		if len(clients) == 0 {
			delete(gw.channels, chanName)
		}
	}
}

func (gw *Gateway) readPump(client *Client) {
	defer func() {
		gw.unregisterClient(client)
		client.conn.Close()
	}()

	client.conn.SetReadLimit(int64(gw.maxSize))
	client.conn.SetReadDeadline(time.Now().Add(gw.pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(gw.pongWait))
		return nil
	})

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %+v", err)
			}
			break
		}

		jmsg := &Message{}

		err = json.Unmarshal(msg, jmsg)
		if err != nil {
			continue
		}

		if jmsg.Action == "join" && jmsg.Channel != "" {
			gw.logger.Logf(zapcore.InfoLevel, "%+v join channel %+v", client, jmsg.Channel)
			gw.joinChannel(jmsg.Channel, client)
			continue
		}

		if jmsg.Action == "leave" && jmsg.Channel != "" {
			gw.logger.Logf(zapcore.InfoLevel, "%+v leave channel %+v", client, jmsg.Channel)
			gw.leaveChannel(jmsg.Channel, client)
			continue
		}

		log.Println("receive", string(msg))
	}
}

func (gw *Gateway) writePump(client *Client) {
	ticker := time.NewTicker(gw.pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(gw.writeWait))
			if !ok {
				// The hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(gw.writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (gw *Gateway) SendChannel(name string, msg []byte) {
	channel, ok := gw.channels[name]
	if !ok {
		return
	}

	for client := range channel {
		client.send <- msg
	}
}

func (gw *Gateway) Broadcast(msg []byte) {
	for client := range gw.clients {
		client.send <- msg
	}
}

func (gw *Gateway) On(event string, listener interfaces.IGatewayListener) error {
	switch event {
	case "channel:join", "channel:leave":
		if _, ok := gw.listeners[event]; !ok {
			gw.listeners[event] = []interfaces.IGatewayListener{}
		}
		gw.listeners[event] = append(gw.listeners[event], listener)
	default:
		return fmt.Errorf("unable to listen to unknown events")
	}

	return nil
}
