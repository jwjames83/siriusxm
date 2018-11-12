package siriusxm

import (
	"encoding/json"
)

func skipRoot(jsonBlob []byte) json.RawMessage {
	var root map[string]json.RawMessage

	if err := json.Unmarshal(jsonBlob, &root); err != nil {
		panic(err)
	}
	for _, v := range root {
		return v
	}
	return nil
}
