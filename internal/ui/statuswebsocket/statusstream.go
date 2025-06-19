package statuswebsocket

import (
	"log"
	"net/http"

	"github.com/digital-dream-labs/vector-bluetooth/ble"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
} // use default options

// StatusStreamer contains the channel being used
type StatusStreamer struct {
	ch chan ble.StatusChannel
}

// New returns a populated StatusStreamer
func New(ch chan ble.StatusChannel) *StatusStreamer {
	return &StatusStreamer{
		ch: ch,
	}
}

// DownloadStatus sends out data
func (s *StatusStreamer) DownloadStatus(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	for {
		b := <-s.ch
		log.Printf("downloadstatus: OTAStatus: %+v\n", b.OTAStatus)
		_ = c.WriteJSON(b)
	}
}
