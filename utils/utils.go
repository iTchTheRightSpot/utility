package utils

import "time"

type RequestContextKey string

type RequestBody struct {
	Id     string `json:"request_id,omitempty"`
	Ip     string `json:"ip_address,omitempty"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

const RequestKey RequestContextKey = "REQUEST_KEY"

func Timezone(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = "UTC"
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	return location, nil
}