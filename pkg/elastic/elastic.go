package elastic

import (
	"edetector_API/config"
	"fmt"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

var es *elasticsearch.Client

func SetElkClient() error {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{"http://" + config.Viper.GetString("ELASTIC_HOST") + ":" + config.Viper.GetString("ELASTIC_PORT")},
	}
	es, err = elasticsearch.NewClient(cfg)
	return err
}

func getAllIndices() ([]string, error) {
	res, err := es.Cat.Indices(
		es.Cat.Indices.WithIndex("*"),
		es.Cat.Indices.WithFormat("json"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var indices []map[string]interface{}
	if res.IsError() {
		return nil, fmt.Errorf("error retrieving indices: %s", res.Status())
	}

	// if err := es.JSONDecode(res.Body, &indices); err != nil {
	// 	return nil, fmt.Errorf("error parsing indices response: %s", err)
	// }

	var indexNames []string
	for _, index := range indices {
		indexNames = append(indexNames, index["index"].(string))
	}

	return indexNames, nil	
}