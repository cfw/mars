package memcached

import "fmt"

type Config struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func (m *Config) Url() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}
