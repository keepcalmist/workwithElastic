package storage

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func InitElastic()(*elasticsearch.Client, error){
	return elasticsearch.NewDefaultClient()
}