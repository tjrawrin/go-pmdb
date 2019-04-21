package render

import (
	"encoding/json"
	"net/http"
)

type data struct {
	Data interface{} `json:"data"`
}

// JSON renders a simple JSON response and sets the content type and status.
func JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	result, err := json.Marshal(data{Data: v})
	if err != nil {
		return err
	}

	w.Write(result)
	return nil
}
