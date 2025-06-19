package logtracer

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	ep_logtracer "github.com/DDLbots/internal-api/go/ep_logtracerpb"
)

// Deprecated
// Trace basically tails the log output of the escape pod
func (s *LogTracer) Trace(req *ep_logtracer.TraceReq, srv ep_logtracer.LogTracer_TraceServer) error {
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
		return err
	}

	w := newWriter()
	t := make(chan time.Time)
	defer func() {
		t <- time.Now()
		_ = reader.Close()
	}()

	go func() {
		_ = reader.Follow(t, w)
	}()

	var line string

	for {
		select {
		case <-srv.Context().Done():
			return nil
		case b := <-w.Chan():
			if string(b) != "\n" {
				line += string(b)
			} else {
				e := logentry{}
				_ = json.Unmarshal([]byte(stripCtlAndExtFromUTF8(line)), &e)
				line = ""

				if e.Function == "stt.(*Server).Parse" {
					if err := srv.Send(
						&ep_logtracer.TraceResp{
							Entry: []*ep_logtracer.Entry{
								{
									Timestamp:    e.Time,
									IncomingText: e.IncomingText,
									Intent:       e.ResultIntent,
								},
							},
						},
					); err != nil {
						return nil
					}
				}
			}
		}
	}
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}
