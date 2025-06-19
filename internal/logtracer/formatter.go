package logtracer

import (
	"fmt"

	"github.com/coreos/go-systemd/sdjournal"
)

func formatter(entry *sdjournal.JournalEntry) (string, error) {
	msg, ok := entry.Fields["MESSAGE"]
	if !ok {
		return "", fmt.Errorf("no MESSAGE field present in journal entry")
	}

	return fmt.Sprintf("%s\n", msg), nil
}
