package proto

import "errors"

type LocationInfo struct {
	Country string
	Region  string
	City    string
	Isp     string
}

type Locator interface {
	Find(ipstr string) (*LocationInfo, error)
}

var ErrUnsupportedIP = errors.New("unsupported ip format")

const Null = "N/A"
