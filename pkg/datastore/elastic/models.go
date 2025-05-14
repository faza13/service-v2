package elastic

type QueryRequest struct {
	Must          []map[string]interface{} `json:"must,omitempty"`
	MustNot       []map[string]interface{} `json:"must_not,omitempty"`
	Should        []map[string]interface{} `json:"should,omitempty"`
	Filter        []map[string]interface{} `json:"filter,omitempty"`
	Sort          []map[string]interface{} `json:"sort,omitempty"`
	Score         []map[string]interface{} `json:"score,omitempty"`
	Collapse      map[string]interface{}   `json:"collapse,omitempty"`
	Aggeregations map[string]interface{}   `json:"aggregations,omitempty"`
	Fields        map[string]interface{}   `json:"fields,omitempty"`
	ScriptFields  map[string]interface{}   `json:"script_fields,omitempty"`
}

type QueryFunction struct {
	FunctionScore FunctionScoreObj `json:"function_score"`
}

type FunctionScoreObj struct {
	Query     QueryObj    `json:"query"`
	Boost     string      `json:"boost"`
	Functions interface{} `json:"functions"`
	MaxBoost  int         `json:"max_boost"`
	ScoreMode string      `json:"score_mode"`
	BoostMode string      `json:"boost_mode"`
	MinScore  int         `json:"min_score"`
}

type QueryObj struct {
	Bool BoolQuery `json:"bool"`
}

type BoolQuery struct {
	Must    []map[string]interface{} `json:"must,omitempty"`
	MustNot []map[string]interface{} `json:"must_not,omitempty"`
	Should  []map[string]interface{} `json:"should,omitempty"`
	Filter  []map[string]interface{} `json:"filter,omitempty"`
}

type QueryRequestRange struct {
	TimeZone string
	GTE      string
	LTE      string
}
