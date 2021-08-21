package amqp

import (
	"fmt"
	"net/url"
)

type Config struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	VirtualHost string `yaml:"virtualHost"`
}

func (c *Config) Url() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		url.QueryEscape(c.VirtualHost))
}
