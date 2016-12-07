package adapters

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"

	"git.xuvasi.com/gocode/faye-go"
	"git.xuvasi.com/gocode/faye-go/transport"
)

/* HTTP handler that can be dropped into the standard http handlers list */
func FayeHandler(server faye.Server) http.Handler {
	// websocketHandler := websocket.Handler(transport.WebsocketServer(server))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {

			server.Logger().Debugf("Client connecting and using web sockets.")

			ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
			if _, ok := err.(websocket.HandshakeError); ok {
				http.Error(w, "Not a websocket handshake", 400)
				return
			} else if err != nil {
				log.Println(err)
				return
			}

			transport.WebsocketServer(server)(ws)
		} else {
			if r.Method == "POST" {
				server.Logger().Debugf("Client connecting and using long-polling.")

				var v interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&v); err == nil {
					transport.MakeLongPoll(v, server, w)
				} else {
					log.Fatal(err)
				}
			} else {
				server.Logger().Debugf("Trying to make a best effort connection...")

				uri := r.RequestURI

				first_parts := strings.Split(uri, "/faye?message=")

				raw_str := strings.Split(first_parts[1], "&")

				json_str, _ := url.QueryUnescape(raw_str[0])

				var v interface{}
				json.Unmarshal([]byte(json_str), v)
				transport.MakeLongPoll(v, server, w)
			}
		}
	})
}

// func handler(w http.ResponseWriter, r *http.Request) {
//     ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
//     if _, ok := err.(websocket.HandshakeError); ok {
//         http.Error(w, "Not a websocket handshake", 400)
//         return
//     } else if err != nil {
//         log.Println(err)
//         return
//     }
//     ... Use conn to send and receive messages.
// }
