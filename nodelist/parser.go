package nodelist

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Parser is responsible for processing input json into nodeLists
type Parser struct {
	nodes      *[]Node
	parseError error
	callback   func(error)
}

// NewParser creates a new Parser...
func NewParser(nodes *[]Node, callback func(error)) Parser {
	return Parser{nodes, fmt.Errorf("Incomplete"), callback}
}

// Parse stuff
func (p *Parser) Parse(jsonData []byte, blocking bool) error {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		(*p.nodes) = append((*p.nodes), NewNode("Root", "", 0))
		p.parseError = p.createNode(data, 0)
		if p.callback != nil {
			p.callback(p.parseError)
		}
	}()
	if blocking {
		wg.Wait()
	}
	return nil
}

// IsComplete will return true if parse was successfully completed
func (p *Parser) IsComplete() bool {
	return p.parseError == nil
}

func (p *Parser) createNode(data interface{}, level int) error {
	parentIndex := len(*p.nodes) - 1
	switch elem := data.(type) {
	case string:
		(*p.nodes)[parentIndex].UpdateValue(strconv.Quote(elem))
	case float64:
		(*p.nodes)[parentIndex].UpdateValue(strconv.FormatFloat(elem, 'g', -1, 64))
	case bool:
		(*p.nodes)[parentIndex].UpdateValue(strconv.FormatBool(elem))
	case nil:
		(*p.nodes)[parentIndex].UpdateValue("null")
	case map[string]interface{}:
		return p.newMapNode(elem, level, parentIndex)
	case []interface{}:
		return p.newArrayNode(elem, level, parentIndex)
	default:
		return fmt.Errorf("Incorretly formatted Json")
	}
	return nil
}

func (p *Parser) newMapNode(data map[string]interface{}, level, parentIndex int) error {
	(*p.nodes)[parentIndex].UpdateValue("{")
	for _, key := range getOrderedMapKeys(data) {
		(*p.nodes) = append((*p.nodes), NewNode(key, "", level+1))
		if err := p.createNode(data[key], level+1); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) newArrayNode(data []interface{}, level, parentIndex int) error {
	(*p.nodes)[parentIndex].UpdateValue("[")
	for index, childInterface := range data {
		(*p.nodes) = append((*p.nodes), NewNode("[]"+strconv.Itoa(index), "", level+1))
		if err := p.createNode(childInterface, level+1); err != nil {
			return err
		}
	}
	return nil
}

func getOrderedMapKeys(data map[string]interface{}) []string {
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
