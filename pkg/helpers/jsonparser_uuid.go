package helpers

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/google/uuid"
)

func GetUUID(data []byte, keys ...string) (res uuid.UUID, err error) {
	if len(keys) == 0 {
		err = errors.New("at least provide a key")
		return
	}
	var (
		str []byte
		vt  jsonparser.ValueType
	)
	if str, vt, _, err = jsonparser.Get(data, keys...); err == nil {
		if len(str) == 4 && vt == jsonparser.Null {
			return
		}
		res, err = uuid.ParseBytes(str)
	}
	return
}
