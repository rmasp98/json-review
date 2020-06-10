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
	nodes    *[]Node
	complete bool
	callback func()
	wg       sync.WaitGroup
}

// NewParser creates a new Parser...
func NewParser(nodes *[]Node, callback func()) Parser {
	return Parser{nodes, false, callback, sync.WaitGroup{}}
}

// Parse stuff
func (p *Parser) Parse(jsonData []byte) error {
	(*p.nodes) = append((*p.nodes), NewNode("Root", "", 0))
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.createNode(data, 0)
	}()
	return nil
}

// WaitForComplete will block the thread until parsing is complete
func (p *Parser) WaitForComplete() {
	p.wg.Wait()
}

// CreateNodeList parses nodeList from file
// func (p *Parser) CreateNodeList(jsonData []byte) error {
// p.nodes = []Node{NewNode("Root", "", 0)}
// var data interface{}
// if err := json.Unmarshal(jsonData, &data); err != nil {
// 	return err
// }
// // Could change this into a goroutine (return err should then kill program)
// if err := p.createNode(data, 0); err != nil {
// 	return err
// }
// p.complete = true
// //////////////////////////////////
// return nil
// }

// GetNodeList returns current list of nodes plus whether it has been full loaded
// func (p Parser) GetNodeList(block bool) ([]Node, bool) {
// return p.nodes, p.complete
// }

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
