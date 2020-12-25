package storage

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func InitElastic()(*elasticsearch.Client, error){
	cfg := elasticsearch.Config{Addresses: []string{
		"http://0.0.0.0:9200",
	} }
	return elasticsearch.NewClient(cfg)
}