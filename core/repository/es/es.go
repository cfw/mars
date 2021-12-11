package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
)

func NewElasticsearch(c *Config) *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: c.Uris,
		Username:  c.Username,
		Password:  c.Password,
	})
	_, err = es.Ping()
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Elasticsearch")
	return es
}
