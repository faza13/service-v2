package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v9"
	"log"
	"service/config"
)

func NewElasticClient(ctx context.Context, cfg *config.Config) *elasticsearch.Client {
	elasticCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.Elastic.Host,
		},
		//Transport: nrelasticsearch.NewRoundTripper(nil),
	}

	es, err := elasticsearch.NewClient(elasticCfg)

	if err != nil {
		log.Fatal("failed to connect to elasticsearch:")
	}

	return es
}
