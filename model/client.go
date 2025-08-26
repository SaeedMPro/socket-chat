package model

import "fmt"

type Client struct {
	Host   string
	Port   int
	UIPort int  `json:"ui_port"`
}

func (c *Client) Address() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

type Config struct {
	ClientOne Client `json:"client-one"`
	ClientTwo Client `json:"client-two"`
}
