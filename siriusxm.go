package siriusxm

import (
	"encoding/json"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

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
