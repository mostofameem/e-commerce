package db

const (
	CheckValue     = "check_value"
	CheckExistence = "check_existence"
)

var OperatorsMap = map[string]string{
	"eq":  "==",
	"ne":  "!=",
	"lt":  "<",
	"lte": "<=",
	"gt":  ">",
	"gte": ">=",
}

type Operator struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type OperationGroup struct {
	Label       string     `json:"label"`
	Value       string     `json:"value"`
	Operators   []Operator `json:"operators"`
	QueryParams []string   `json:"queryParams"`
	UI          string     `json:"ui"`
}

var OperationGroups = []OperationGroup{
	{
		Label: "Check value",
		Value: CheckValue,
		Operators: []Operator{
			{
				Label: "Equal to",
				Value: "eq",
			},
			{
				Label: "Not equal to",
				Value: "ne",
			},
			{
				Label: "Greater than",
				Value: "gt",
			},
			{
				Label: "Less than",
				Value: "lt",
			},
			{
				Label: "Greater than or equal to",
				Value: "gte",
			},
			{
				Label: "Less than or equal to",
				Value: "lte",
			},
		},
		QueryParams: []string{"paths", "operators", "values"},
		UI:          "path_value_check",
	},
	{
		Label:       "Check existence",
		Value:       CheckExistence,
		QueryParams: []string{"props"},
		UI:          "existence_check",
	},
}
