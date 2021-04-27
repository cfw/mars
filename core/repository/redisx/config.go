package redisx

import "fmt"

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
}

func (r *Config) Url() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}