package json

import (
	"encoding/json"
)

func Serialize(v interface{}) string {

	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func SerializeToByte(v interface{}) []byte {

	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bytes
}

func Deserialize(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}
