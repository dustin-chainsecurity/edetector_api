package elastic

// import (
// 	"edetector_API/config"
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"github.com/elastic/go-elasticsearch/v6"
// 	"github.com/elastic/go-elasticsearch/v6/esapi"
// )

// var es *elasticsearch.Client

// func SetElkClient() error {
// 	var err error
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{"http://" + config.Viper.GetString("ELASTIC_HOST") + ":" + config.Viper.GetString("ELASTIC_PORT")},
// 	}
// 	es, err = elasticsearch.NewClient(cfg)
// 	return err
// }

// func GetAllIndices() ([]string, error) {
// 	res, err := es.Cat.Indices(
// 		es.Cat.Indices.WithIndex("*"),
// 		es.Cat.Indices.WithFormat("json"),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	var indices []map[string]interface{}
// 	if res.IsError() {
// 		return nil, fmt.Errorf("error retrieving indices: %s", res.Status())
// 	}
// 	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
// 		return nil, fmt.Errorf("error parsing indices response: %s", err)
// 	}

// 	var indexNames []string
// 	for _, index := range indices {
// 		indexNames = append(indexNames, index["index"].(string))
// 	}
// 	return indexNames, nil
// }

// func ClearAllIndices() error {
// 	query := `{
// 		"query": {
// 			"match_all": {}
// 		}
// 	}`
	
// 	indices, err := GetAllIndices()
// 	if err != nil {
// 		return err
// 	}

// 	for _, index := range indices {
// 		if err := clearIndex(index, query); err != nil {
// 			return err
// 		} else {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func clearIndex(index string, query string) error {
// 	req := esapi.DeleteByQueryRequest{
// 		Index:        []string{index},
// 		Query:        strings.NewReader(query),
// 		Refresh:      "true", // Refresh the index after deletion for immediate effect
// 		IgnoreErrors: true,   // Continue even if some documents could not be deleted
// 	}

// 	res, err := req.Do(context.Background(), client)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return fmt.Errorf("error deleting documents from index %s: %s", index, res.Status())
// 	}
	
// 	return nil
// }