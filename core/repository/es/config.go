package es

type Config struct {
	Uris     []string `yaml:"uris"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
