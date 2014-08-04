package transport

import (
	"github.com/gorilla/websocket"
	// "code.google.com/p/go.net/websocket"
	"github.com/roncohen/faye/protocol"
	"log"
)

const WebSocketConnectionPriority = 10

type Server interface {
	HandleRequest(interface{}, protocol.Connection)
}

type WebSocketConnection struct {
	ws         *websocket.Conn
	failedSend bool
}

func (wc *WebSocketConnection) Send(msgs []protocol.Message) error {
	err := wc.ws.WriteJSON(msgs)
	log.Printf("Writing to websocket: %v", msgs)
	if err != nil {
		wc.failedSend = true
	}
	return err
}

func (wc *WebSocketConnection) IsConnected() bool {
	return wc.failedSend
}

func (wc *WebSocketConnection) Close() {
	wc.ws.Close()
}

func (wc WebSocketConnection) Priority() int {
	return WebSocketConnectionPriority
}

func (lp WebSocketConnection) IsSingleShot() bool {
	return false
}

func WebsocketServer(m Server) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		var data interface{}
		wsConn := WebSocketConnection{ws, true}
		for {
			err := ws.ReadJSON(&data)
			if err != nil {
				log.Print(err)
				return
			}

			if arr := data.([]interface{}); len(arr) == 0 {
				ws.WriteJSON([]string{})
			} else {
				m.HandleRequest(data, &wsConn)
			}
		}
	}
}
