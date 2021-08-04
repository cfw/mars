package sqlx

import "fmt"

type Config struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Driver      string `yaml:"driver"`
	Debug       bool   `yaml:"debug"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
}

func (d *Config) Url() string {
	return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s", d.Host, d.Port, d.Username, d.Password, d.Database)
}
