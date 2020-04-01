package x

import "encoding/json"

func ReadJson(b []byte, val interface{}) {
	FailFast(json.Unmarshal(b, val))
}
