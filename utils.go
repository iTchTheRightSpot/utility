package main

import "time"

type RequestContextKey string

type RequestBody struct {
	Id     string `json:"id"`
	Ip     string `json:"ip"`
	Method string `json:"method"`
	Path   string `json:"path"`
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
