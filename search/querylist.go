package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// QueryList allow viewing, editing and saving of the list of Common Misconfigurations
type QueryList struct {
	list map[string]QueryData
}

// QueryData contains all data for a particular misconfiguration
type QueryData struct {
	Regex       string
	Description string
}

// NewQueryList stuff
func NewQueryList() QueryList {
	return QueryList{make(map[string]QueryData)}
}

// GetNames stuff
func (c QueryList) GetNames() []string {
	var names []string
	for key := range c.list {
		names = append(names, key)
	}
	sort.Strings(names)
	return names
}

// GetDescription stuff
func (c QueryList) GetDescription(name string) string {
	return c.list[name].Description
}

// GetQuery stuff
func (c QueryList) GetQuery(name string) (string, QueryEnum) {
	return c.list[name].Regex, REGEX
}

// Add stuff
func (c *QueryList) Add(name string, regex string, description string) {
	c.list[name] = QueryData{regex, description}
}

// Remove stuff
func (c *QueryList) Remove(name string) {
	delete(c.list, name)
}

// Save stuff
func (c QueryList) Save(file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	defer f.Close()
	if err != nil {
		return err
	}
	outJSON, marshelErr := json.MarshalIndent(c.list, "", "   ")
	if marshelErr != nil {
		return marshelErr
	}
	fmt.Println(string(outJSON))
	_, writeErr := f.Write(outJSON)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// Load stuff
func (c *QueryList) Load(file string) error {
	rawJSON, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		return readErr
	}
	unmarshalErr := json.Unmarshal(rawJSON, &c.list)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	return nil
}
