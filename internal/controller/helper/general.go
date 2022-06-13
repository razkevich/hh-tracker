package helper

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"math"
	"strings"
	"time"
)

// EmptyObj object is used when data doesn't want to be null on json
type EmptyObj struct{} // @name Data.EmptyObject

// RelationData Relationship data
type RelationData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
} // @name Response.RelationData

// Links browse between different pages
type Links struct {
	Current string `json:"current,omitempty"`
	First   string `json:"first,omitempty"`
	Last    string `json:"last,omitempty"`
	Next    string `json:"next,omitempty"`
	Prev    string `json:"prev,omitempty"`
	Self    string `json:"self,omitempty"`
} // @name Links

// Meta meta information, including timestamps and paging info
type Meta struct {
	Timestamps *Timestamps        `json:"timestamps,omitempty"`
	Page       *PaginationPage    `json:"page,omitempty"`
	Results    *PaginationResults `json:"results,omitempty"`
} // @name Response.Meta

// PaginationPage pagination info about the page
type PaginationPage struct {
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
} // @name params.PaginationPage

// PaginationResults pagination info about the results
type PaginationResults struct {
	Total int64 `json:"total"`
} // @name params.PaginationResults

// Timestamps timestamps
type Timestamps struct {
	CreatedAt FormattedTimestamp `gorm:"column:created_at" form:"created_at" json:"created_at" time_format:"2006-01-02" time_utc:"1" binding:"required"`
	UpdatedAt FormattedTimestamp `gorm:"column:updated_at" form:"updated_at" json:"updated_at" time_format:"2006-01-02" time_utc:"1" binding:"required"`
	DeletedAt FormattedTimestamp `gorm:"column:deleted_at" form:"deleted_at" json:"-"`
} // @name Timestamps

// FormattedTimestamp Timestamp that's properly formatted as per the style guide
type FormattedTimestamp struct {
	time.Time
} // @name FormattedTimestamp

// MarshalJSON Marshall JSON for FormattedTimestamp
func (mt FormattedTimestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(FormattedTimestamp.Format(mt, "2006-01-02T15:04:05.881Z"))
}

// GetPaginationPageAndResults Returns the page and results objects for meta info on pagination
func GetPaginationPageAndResults(limit int64, offset int64, totalCount int64) (PaginationPage, PaginationResults) {
	currentPage := (offset / limit) + 1
	totalPages := int64(math.Ceil(float64(totalCount) / float64(limit)))
	page := PaginationPage{
		Limit:   limit,
		Offset:  offset,
		Current: currentPage,
		Total:   totalPages,
	}
	results := PaginationResults{Total: totalCount}
	return page, results
}

type rqlJSON struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

// RqlOpr map of available Rql operations
var RqlOpr = map[string]string{
	"eq":   "=",
	"ne":   "!=",
	"lt":   "<",
	"gt":   ">",
	"le":   "<=",
	"ge":   ">=",
	"like": "ILIKE",
}

// ConvertRQLToJSON convert RQL JSON
func ConvertRQLToJSON(rql string) (string, error) {
	var top rqlJSON
	if err := json.Unmarshal([]byte(rql), &top); err != nil {
		return "", err
	}
	var args []rqlJSON
	if err := json.Unmarshal(top.Args, &args); err != nil {
		return "", err
	}
	filter := make(map[string]map[string]interface{})
	for _, opr := range args {
		_, ok := RqlOpr[opr.Name]
		if !ok {
			return "", fmt.Errorf("unknown operator %s", opr.Name)
		}
		var args []interface{}
		if err := json.Unmarshal(opr.Args, &args); err != nil {
			return "", err
		}
		if len(args) < 2 {
			return "", fmt.Errorf("bad arguments for %s", opr.Name)
		}
		name, ok := args[0].(string)
		if !ok {
			return "", fmt.Errorf("bad arguments for %s", opr.Name)
		}
		name = strings.ToLower(name)
		if funk.IsEmpty(filter[name]) {
			filter[name] = make(map[string]interface{})
		}
		filter[name][opr.Name] = args[1]
	}
	filterForm := make(map[string]interface{})
	for name, args := range filter {
		if len(args) == 1 {
			for arg, val := range args {
				if arg == "$eq" {
					filterForm[name] = val
				} else {
					filterForm[name] = args
				}
			}
		} else {
			filterForm[name] = args
		}
	}
	out := struct {
		Filter map[string]interface{} `json:"filter"`
	}{filterForm}
	bytes, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
