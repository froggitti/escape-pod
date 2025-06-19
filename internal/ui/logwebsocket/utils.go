package logwebsocket

import (
	"fmt"
	"strings"

	"github.com/coreos/go-systemd/sdjournal"
)

func formatter(entry *sdjournal.JournalEntry) (string, error) {
	msg, ok := entry.Fields["MESSAGE"]
	if !ok {
		return "", fmt.Errorf("no MESSAGE field present in journal entry")
	}

	return fmt.Sprintf("%s\n", msg), nil
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}
