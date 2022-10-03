package qsx

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
)

type Websocket interface {
	Dial() (*websocket.Conn, error)
}

type Dialer struct {
	URL string
}

// Dial returns a websocket connection
func (w *Dialer) Dial() (*websocket.Conn, error) {
	var wsDialer websocket.Dialer
	wsConn, resp, err := wsDialer.Dial(w.URL, nil)
	if err != nil {
		return nil, err
	}
	_, _ = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return wsConn, nil
}
