package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type Memcached struct {
	client *memcache.Client
}

func NewMemcached(c *Config) *Memcached {
	client := memcache.New(c.Url())
	err := client.Ping()
	if err != nil {
		panic(err)
	}
	return &Memcached{
		client: client,
	}
}
func (m *Memcached) Set(key string, value []byte, expire time.Duration) {
	_ = m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(expire.Seconds())},
	)
}
func (m *Memcached) Get(key string) (string, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}
