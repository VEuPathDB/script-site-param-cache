package x

import "encoding/json"

func JsonMarshal(val interface{}) []byte {
	tmp, err := json.Marshal(val)
	FailFast(err)
	return tmp
}
