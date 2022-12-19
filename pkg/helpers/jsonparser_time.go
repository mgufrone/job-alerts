package helpers

import (
	"errors"
	"github.com/buger/jsonparser"
	"time"
)

func GetTime(data []byte, keys ...string) (res time.Time, err error) {
	if len(keys) == 0 {
		err = errors.New("at least provide a key")
		return
	}
	var s int64
	if s, err = jsonparser.GetInt(data, keys...); err == nil {
		res = time.UnixMilli(s)
	}
	return
}
