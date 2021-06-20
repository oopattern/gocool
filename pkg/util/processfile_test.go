package util

import (
	"errors"
	"github.com/oopattern/gocool/log"
	"strings"
	"testing"
)

type Recommends struct {
	Recoms map[string]int64 `json:"recoms"`
}

func (r *Recommends) ParseLine(line string, tp string) error {
	fields := strings.Split(line, "\t")
	if len(fields) < 1 {
		log.Error("parse line[%s] format error", line)
		return errors.New("parse line format error")
	}

	id := fields[0]
	r.Recoms[id] = 1
	return nil
}

func TestProcessFile(t *testing.T) {
	r := Recommends{Recoms: make(map[string]int64)}

	if err := HandleFile("recommend_file", "", &r); err != nil {
		t.Error(err)
	}
}
