package up

import "encoding/json"

// Request struct
type Request struct {
	Method     string            `json:"method"`
	URI        string            `json:"uri"`
	Parameters map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}

func NewRequestFromRawString(raw string) (Request, error) {
	req := Request{}
	err := json.Unmarshal([]byte(raw), &req)
	return req, err
}
