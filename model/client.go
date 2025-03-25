package model

import (
	"strconv"
)

type Client struct {
	Host string
	Port int
}

func (c Client) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

type Config struct {
	ClientOne Client `json:"client-one"`
	ClientTwo Client `json:"client-two"`
}
