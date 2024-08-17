package goutil

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type encodeJson struct{}

var JSON *encodeJson = &encodeJson{}

// StringifyJSON converts a map or array to a JSON string
func (encJson *encodeJson) Stringify(data interface{}, ind ...int) ([]byte, error) {
	var res []byte
	var err error
	if len(ind) != 0 {
		sp := "  "
		if len(ind) > 2 {
			sp = strings.Repeat(" ", ind[1])
		}
		res, err = json.MarshalIndent(data, strings.Repeat(" ", ind[0]), sp)
	} else {
		res, err = json.Marshal(data)
	}

	if err != nil {
		return []byte{}, err
	}
	res = bytes.ReplaceAll(res, []byte("\\u003c"), []byte("<"))
	res = bytes.ReplaceAll(res, []byte("\\u003e"), []byte(">"))

	return res, nil
}

// ParseJson converts a json string into a map of strings
func (encJson *encodeJson) Parse(b []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return res, nil
}

// DecodeJSON is useful for decoding a JSON output from the body of an http request
//
// example: goutil.DecodeJSON(r.Body)
func (encJson *encodeJson) Decode(data io.Reader) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := json.NewDecoder(data).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeepCopyJson will stringify and parse json to create a deep copy and escape pointers
func (encJson *encodeJson) DeepCopy(data map[string]interface{}) (map[string]interface{}, error) {
	b, err := encJson.Stringify(data)
	if err != nil {
		return nil, err
	}
	return encJson.Parse(b)
}
