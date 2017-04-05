package main

import (
	"fmt"
)

type Endpoint struct {
	area      string
	serviceID string
}

func NewEndpoint(area string, serviceID string) Endpoint {
	return Endpoint{
		area:      area,
		serviceID: serviceID,
	}
}

func (e Endpoint) List() string {
	return fmt.Sprintf("https://smartcast.hs.llnwd.net/%s/%s/%s.txt", e.area, e.serviceID, e.serviceID)
}

func (e Endpoint) TS(name string) string {
	return fmt.Sprintf("https://smartcast.hs.llnwd.net/%s/%s/%s", e.area, e.serviceID, name)
}

func TrimSerial(tsFilename string) string {
	return tsFilename[3:]
}
