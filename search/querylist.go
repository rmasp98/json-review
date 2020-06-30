package search

import (
	"fmt"
	"kube-review/utils"
	"regexp"
	"sort"
	"strings"
)

// QueryList allow viewing, editing and saving of the list of Common Misconfigurations
type QueryList struct {
	list map[string]QueryData
}

// QueryData contains all data for a particular misconfiguration
type QueryData struct {
	Query       string    `json:"query"`
	Description string    `json:"description"`
	QueryType   QueryEnum `json:"queryType"`
}

// NewQueryList stuff
func NewQueryList() QueryList {
	return QueryList{make(map[string]QueryData)}
}

// GetNames stuff
func (q QueryList) GetNames() []string {
	var names []string
	for key := range q.list {
		names = append(names, key)
	}
	sort.Strings(names)
	return names
}

// GetDescription stuff
func (q QueryList) GetDescription(name string) string {
	return q.list[name].Description
}

// GetQuery stuff
func (q QueryList) GetQuery(name string) (string, QueryEnum) {
	return q.list[name].Query, q.list[name].QueryType
}

// GetHints stuff
func (q QueryList) GetHints(input string) []string {
	var hints []string
	r, _ := regexp.Compile(input)
	for _, query := range q.GetNames() {
		if input == query {
			return []string{}
		}
		if r.MatchString(query) {
			hints = append(hints, query+" - "+redBold+q.GetDescription(query)+reset)
		}
	}
	return hints
}

// InsertHint stuff
func (q QueryList) InsertHint(input string, index int) string {
	hints := q.GetHints(input)
	if index < len(hints) {
		return strings.Split(hints[index], " - ")[0]
	}
	return input
}

// Add stuff
func (q *QueryList) Add(name string, query string, description string, queryType QueryEnum) error {
	if queryType != REGEX && queryType != EXPRESSION {
		return fmt.Errorf("Invalid query type. Must be either Regex or Intelligent")
	}
	q.list[name] = QueryData{query, description, queryType}
	return nil
}

// Remove stuff
func (q *QueryList) Remove(name string) {
	delete(q.list, name)
}

// Save stuff
func (q QueryList) Save(file string) error {
	return utils.SaveJSON(file, q.list, true)
}

// Load stuff
func (q *QueryList) Load(file string, schemaFile string) error {
	schema, err := utils.Load(schemaFile)
	if schemaFile != "" && err != nil {
		return err
	}
	return utils.LoadJSON(file, &q.list, schema)
}
