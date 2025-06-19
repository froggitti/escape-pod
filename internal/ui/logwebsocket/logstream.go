package logwebsocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

type logentry struct {
	Function     string `json:"Function"`
	Source       string `json:"Source"`
	IncomingText string `json:"incoming_text"`
	Level        string `json:"level"`
	Msg          string `json:"msg"`
	ResultIntent string `json:"result_intent"`
	Time         string `json:"time"`
}

// Log streams the log entries over a websocket
func Log(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	reader, err := getReader()
	if err != nil {
		log.Println(err)
		return
	}

	wr := newWriter()
	t := make(chan time.Time)
	defer func() {
		t <- time.Now()
		_ = reader.Close()
	}()

	go func() {
		_ = reader.Follow(t, wr)
	}()

	var line string

	for {
		b := <-wr.Chan()
		if string(b) != "\n" {
			line += string(b)
		} else {
			e := logentry{}
			_ = json.Unmarshal([]byte(stripCtlAndExtFromUTF8(line)), &e)
			line = ""

			if e.IncomingText != "" {
				if err := c.WriteJSON(e); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// GetReader returns a new journal reader
func getReader() (*sdjournal.JournalReader, error) {
	reader, err := sdjournal.NewJournalReader(
		sdjournal.JournalReaderConfig{
			Since: time.Duration(-1),
			Matches: []sdjournal.Match{
				{
					Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
					Value: "escape_pod.service",
				},
			},
			Formatter: formatter,
		},
	)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func ParsedIntent(parsed <-chan interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer func() {
			_ = c.Close()
		}()
		for p := range parsed {
			if err := c.WriteJSON(p); err != nil {
				log.Printf("write json: %+v; %v\n", p, err)
				break
			}
		}

	}
}
