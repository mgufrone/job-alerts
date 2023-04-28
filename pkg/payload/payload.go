package payload

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"strconv"
	"strings"
)

type Payload struct {
	data []byte
}

func New(data []byte) Payload {
	return Payload{data: data}
}

func keyParser(key string) []string {
	keys := strings.Split(key, ".")
	for k, v := range keys {
		if val, err := strconv.Atoi(v); err == nil && val >= 0 {
			keys[k] = "[" + v + "]"
		}
	}
	return keys
}
func (p *Payload) Get(key string) ([]byte, error) {
	by, _, _, err := jsonparser.Get(p.data, keyParser(key)...)
	return by, err
}
func (p *Payload) As(key string, out interface{}) error {
	if key == "" {
		return json.Unmarshal(p.data, out)
	}
	by, err := p.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(by, out)
}
func (p *Payload) String(key string) (string, error) {
	return jsonparser.GetUnsafeString(p.data, keyParser(key)...)
}
func (p *Payload) Int(key string) (int64, error) {
	return jsonparser.GetInt(p.data, keyParser(key)...)
}

func (p *Payload) Bool(key string) (bool, error) {
	return jsonparser.GetBoolean(p.data, keyParser(key)...)
}
func (p *Payload) JSON() ([]byte, error) {
	return p.data, nil
}
