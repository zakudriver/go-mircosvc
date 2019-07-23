package utils

import "encoding/json"

// strut -> json
func Struct2Json(stut interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(stut)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})

	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
