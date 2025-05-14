package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	jsoniter "github.com/json-iterator/go"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Elastic struct {
	esClient *elasticsearch.Client
	index    string
}

func NewElastic(es *elasticsearch.Client, index string) *Elastic {
	return &Elastic{
		esClient: es,
		index:    index,
	}
}

func (elastic *Elastic) bulkIndex(ctx context.Context, jsonData []interface{}) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  elastic.index,    // The default index name
		Client: elastic.esClient, // The Elasticsearch client
		//NumWorkers:    constants.NumWorkersBulkIndex, // The number of worker goroutines
		//FlushBytes:    constants.FlushBytesBulkIndex, // The flush threshold in bytes
		//FlushInterval: 30 * time.Second,              // The periodic flush interval
	})

	indexID := "id"

	for _, data := range jsonData {
		// Prepare the data payload: encode article to JSON
		payload, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Cannot encode data with ID %d: %s", data.(map[string]interface{})[indexID], err)
		}

		err = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(int(data.(map[string]interface{})[indexID].(float64))),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(payload),

				// OnFailure is called for each failed operation
				OnFailure: elastic.onBulkIndexFailure,
			},
		)
		if err != nil {
			log.Printf("Error adding data with ID %d to bulk indexer: %s", data.(map[string]interface{})[indexID], err)
		}
	}

	if err := bi.Close(ctx); err != nil {
		log.Printf("Error closing bulk indexer: %s", err)
	}

	return err
}

func (elastic *Elastic) onBulkIndexFailure(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	} else {
		log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
	}
}

func (elastic *Elastic) createOrUpdate(ctx context.Context, id int, jsonData interface{}) error {
	req := esapi.IndexRequest{
		Index:      elastic.index,
		DocumentID: strconv.Itoa(id),
		Body:       esutil.NewJSONReader(&jsonData),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, elastic.esClient)
	if err != nil {
		return fmt.Errorf("create or update elastic error: %s", err)
	}
	defer res.Body.Close()
	return nil
}

func (elastic *Elastic) searchDataPagination(ctx context.Context, jsonSearch map[string]interface{}, dest interface{}) error {
	es := elastic.esClient
	jsonByte, _ := json.Marshal(jsonSearch)
	buf := bytes.NewReader(jsonByte)
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(elastic.index),
		es.Search.WithBody(buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("error search elastic: %s", err)
	}
	if res.StatusCode != 200 {
		return errors.Join(
			fmt.Errorf("error status code: %d", res.StatusCode),
			fmt.Errorf("json_search: %s", jsonSearch),
		)
	}

	if err := jsoniter.NewDecoder(res.Body).Decode(dest); err != nil {
		return fmt.Errorf("error decode json: %s", err)
	}

	return nil
}

func (elastic *Elastic) getById(ctx context.Context, id int, dest interface{}) error {
	req := esapi.GetRequest{
		Index:      elastic.index,
		DocumentID: strconv.Itoa(id),
	}
	res, err := req.Do(ctx, elastic.esClient)
	if err != nil {
		return fmt.Errorf("get by id elastic error: %s", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(dest)
	return nil
}

func (elastic *Elastic) update(ctx context.Context, id int, jsonData interface{}) error {
	req := esapi.UpdateRequest{
		Index:      elastic.index,
		DocumentID: strconv.Itoa(id),
		Body:       esutil.NewJSONReader(&jsonData),
	}
	res, err := req.Do(ctx, elastic.esClient)
	if err != nil {
		return fmt.Errorf("update elastic error: %s", err)
	}
	defer res.Body.Close()
	return nil
}

func (elastic *Elastic) deleteByIDs(arrId []int) error {
	jsonData := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": arrId,
			},
		},
	}
	res, err := elastic.esClient.DeleteByQuery([]string{elastic.index}, esutil.NewJSONReader(&jsonData))
	if err != nil {
		return fmt.Errorf("delete by ids elastic error: %s", err)
	}
	defer res.Body.Close()
	return nil
}

func (elastic *Elastic) generateQuery(json_objects *QueryRequest, start int, size int) map[string]interface{} {
	result_query := QueryFunction{}
	result_query.FunctionScore.MaxBoost = 100
	result_query.FunctionScore.MinScore = 0
	result_query.FunctionScore.ScoreMode = "sum"
	result_query.FunctionScore.BoostMode = "sum"
	result_query.FunctionScore.Boost = "15"
	result_query.FunctionScore.Query.Bool.Must = json_objects.Must
	result_query.FunctionScore.Query.Bool.MustNot = json_objects.MustNot
	result_query.FunctionScore.Query.Bool.Should = json_objects.Should
	result_query.FunctionScore.Query.Bool.Filter = json_objects.Filter
	if len(json_objects.Score) != 0 {
		result_query.FunctionScore.Functions = json_objects.Score
	}

	output_map_query := make(map[string]interface{})
	output_map_query["from"] = start
	output_map_query["size"] = size
	output_map_query["query"] = result_query
	if len(json_objects.Sort) != 0 {
		output_map_query["sort"] = json_objects.Sort
	}
	if len(json_objects.Aggeregations) != 0 {
		output_map_query["aggregations"] = json_objects.Aggeregations
	}
	if len(json_objects.Collapse) != 0 {
		output_map_query["collapse"] = json_objects.Collapse
	}
	if len(json_objects.Fields) != 0 {
		output_map_query["_source"] = json_objects.Fields
	}
	if len(json_objects.ScriptFields) != 0 {
		output_map_query["script_fields"] = json_objects.ScriptFields
	}

	return output_map_query
}

func (elastic *Elastic) GetPagination(page int, limit int) (startElastic int, limitElastic int) {
	return GetPagination(page, limit)
}

func (elastic *Elastic) GenerateMatchAndQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"match": map[string]interface{}{
			field: map[string]interface{}{
				"query":    value,
				"operator": "and",
			},
		},
	}
}

func (elastic *Elastic) GenerateMatchQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"match": map[string]interface{}{
			field: value,
		},
	}
}

func (elastic *Elastic) GenerateMatchPhraseQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"match_phrase": map[string]interface{}{
			field: value,
		},
	}
}

func (elastic *Elastic) GenerateTermQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	}
}

func (elastic *Elastic) GenerateFilterRangeQuery(field string, value interface{}, operand string) map[string]interface{} {
	return map[string]interface{}{
		"range": map[string]interface{}{
			field: map[string]interface{}{
				operand: value,
			},
		},
	}
}

func (elastic *Elastic) GenerateFilterRangeWithTimeZoneQuery(field string, data QueryRequestRange) map[string]interface{} {
	rangeQuery := map[string]interface{}{
		"time_zone": data.TimeZone,
	}

	if data.GTE != "" {
		rangeQuery["gte"] = data.GTE
	}

	if data.LTE != "" {
		rangeQuery["lte"] = data.LTE
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			field: rangeQuery,
		},
	}
}

func (elastic *Elastic) GenerateMatchListQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"terms": map[string]interface{}{
			field: value,
		},
	}
}

func (elastic *Elastic) GenerateMatchPrefixQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"match_phrase_prefix": map[string]interface{}{
			field: value,
		},
	}
}

func (elastic *Elastic) GenerateExistQuery(field string) map[string]interface{} {
	return map[string]interface{}{
		"exist": map[string]interface{}{
			"field": field,
		},
	}
}

func (elastic *Elastic) GenerateNestedQuery(path string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path":  path,
			"query": value,
		},
	}
}

func (elastic *Elastic) AddMatchQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey]
	if ok {
		matchQuery := elastic.GenerateMatchQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddMatchPhraseQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey]
	if ok {
		matchQuery := elastic.GenerateMatchPhraseQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddMatchListQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey]
	if ok {
		matchQuery := elastic.GenerateMatchListQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddWildCardQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey]
	if ok {
		matchQuery := elastic.GenerateWildCardLikeQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddTermQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey]
	if ok {
		matchQuery := elastic.GenerateTermQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddMatchPrefixQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey].(string)
	if ok {
		matchQuery := elastic.GenerateMatchPrefixQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddMatchAndQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey].(string)
	if ok {
		matchQuery := elastic.GenerateMatchAndQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddMatchAndQueryPregReplaceIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field string) {
	value, ok := customFilter[filterKey].(string)
	if ok {
		reg, _ := regexp.Compile("[^a-zA-Z0-9]")
		value = reg.ReplaceAllString(value, " ")
		matchQuery := elastic.GenerateMatchAndQuery(field, value)
		jsonData.Must = append(jsonData.Must, matchQuery)
	}
}

func (elastic *Elastic) AddRangeQueryIfPresent(customFilter map[string]interface{}, jsonData *QueryRequest, filterKey, field, operator string) {
	if value, ok := customFilter[filterKey]; ok {
		rangeQuery := elastic.GenerateFilterRangeQuery(field, value, operator)
		jsonData.Filter = append(jsonData.Filter, rangeQuery)
	}
}

func (elastic *Elastic) AddSortQuery(customFilter map[string]interface{}, jsonData *QueryRequest, defaultSort map[string]interface{}) {
	if values, ok := customFilter["sort"]; ok {
		sortFields := strings.Split(values.(string), ",")
		for _, field := range sortFields {
			fieldName, ordering := CheckOrderingFromString(field)
			jsonData.Sort = append(jsonData.Sort, map[string]interface{}{
				fieldName: ordering,
			})
		}
	} else if len(defaultSort) > 0 {
		for field, order := range defaultSort {
			jsonData.Sort = append(jsonData.Sort, map[string]interface{}{
				field: order,
			})
		}
	}
}

func (elastic *Elastic) GenerateOrQuery(fields map[string]interface{}) map[string]interface{} {
	shouldQuery := []interface{}{}
	for field, value := range fields {
		query := strings.Split(field, "__")
		if len(query) > 1 {
			switch query[1] {
			case "custom":
				shouldQuery = append(shouldQuery, value)
			case "term":
				shouldQuery = append(shouldQuery, elastic.GenerateTermQuery(query[0], value))
			case "list":
				shouldQuery = append(shouldQuery, elastic.GenerateMatchListQuery(query[0], value))
			case "exists":
				shouldQuery = append(shouldQuery, elastic.GenerateExistQuery(query[0]))
			case "wildcard":
				shouldQuery = append(shouldQuery, elastic.GenerateWildCardLikeQuery(query[0], value))
			case "match_phrase_prefix":
				shouldQuery = append(shouldQuery, elastic.GenerateMatchPrefixQuery(query[0], value))
			case "gt", "gte", "lt", "lte":
				shouldQuery = append(shouldQuery, elastic.GenerateFilterRangeQuery(query[0], value, query[1]))
			default:
				shouldQuery = append(shouldQuery, elastic.GenerateMatchQuery(query[0], value))
			}
		} else {
			shouldQuery = append(shouldQuery, elastic.GenerateMatchQuery(query[0], value))
		}
	}
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"should": shouldQuery,
		},
	}
}

func (elastic *Elastic) GenerateWildCardLikeQuery(field string, value interface{}) map[string]interface{} {
	query := fmt.Sprintf("%v", value)
	return map[string]interface{}{
		"wildcard": map[string]interface{}{
			field: map[string]interface{}{
				"value": "*" + query + "*",
			},
		},
	}
}

func GetPagination(page int, limit int) (startElastic int, limitElastic int) {
	start := 0
	if page != 0 {
		start = (int(page) - 1) * limit
	}
	return start, limit
}

func CheckOrderingFromString(value string) (string, string) {
	first_index_value := strings.Split(value, "")[0]
	if first_index_value == "-" {
		field_name := strings.Join(strings.Split(value, "")[1:], "")
		return field_name, "desc"
	}
	return value, "asc"
}
