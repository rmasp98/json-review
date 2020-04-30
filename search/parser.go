package search

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Parse uses the parser class to validate an intelligent search string
// and return a list of commands
func Parse(input string) ([]Command, error) {
	trimmedInput := strings.Trim(input, " ")
	p := parser{trimmedInput, 0, len(trimmedInput), 0, []Command{}}
	return p.parse()
}

type parser struct {
	input           string
	charIterator    int
	lenInput        int
	numOpenBrackets int
	outputCommands  []Command
}

func (p *parser) parse() ([]Command, error) {
	var err error
	var currentCommand Command
	for p.charIterator < p.lenInput {
		if currentCommand.Operator, err = p.checkOperators(); err != nil {
			return []Command{}, p.createError(err.Error())
		}
		p.checkBrackets('(', 1, &currentCommand)
		if currentCommand.Control, err = p.check(controls); err != nil {
			return []Command{}, p.createError("A control is missing or invalid")
		}
		var condition string
		if condition, err = p.check(conditionals); err != nil {
			return []Command{}, p.createError("A condition is missing or invalid")
		}
		currentCommand.Equal = (condition == "==")
		if currentCommand.Regex, err = p.checkQuotes(); err != nil {
			return []Command{}, p.createError(err.Error())
		}
		p.checkBrackets(')', -1, &currentCommand)
		if currentCommand.Control != "" {
			p.outputCommands = append(p.outputCommands, currentCommand)
			currentCommand = Command{"", false, "", "", ""}
		}
	}

	if err := p.validateCommands(); err != nil {
		return []Command{}, err
	}
	return p.outputCommands, nil
}

func (p *parser) stripLeft(num int) {
	p.charIterator += num
	for p.charIterator < p.lenInput && p.input[p.charIterator] == ' ' {
		p.charIterator++
	}
}

func (p *parser) check(checkList []string) (string, error) {
	orderedCheckList := orderCheckList(checkList)
	for _, check := range orderedCheckList {
		var checkEnd = p.charIterator + len(check)
		if checkEnd < p.lenInput && strings.EqualFold(p.input[p.charIterator:checkEnd], check) {
			p.stripLeft(len(check))
			return check, nil
		}
	}
	return "", fmt.Errorf("")
}

func (p *parser) checkQuotes() (string, error) {
	if quote := p.getQuote(); quote != "" {
		substring := strings.Split(p.input[p.charIterator+1:], quote)
		// If quote at end of string, len(subs) would be at least 1
		if len(substring) > 1 {
			p.stripLeft(len(substring[0]) + 2)
			if _, err := regexp.Compile(substring[0]); err != nil {
				return "", fmt.Errorf("Regex is not valid")
			}
			return substring[0], nil
		}
		return "", fmt.Errorf("Missing end quote")
	}
	return "", fmt.Errorf("Regex is missing or has not been quoted")
}

func (p *parser) checkBrackets(bracket byte, incrementor int, command *Command) {
	for p.charIterator < p.lenInput && p.input[p.charIterator] == bracket {
		command.Bracket = string(bracket)
		p.outputCommands = append(p.outputCommands, *command)
		*command = Command{"", false, "", "", ""}
		p.numOpenBrackets = p.numOpenBrackets + incrementor
		p.stripLeft(1)
	}
}

func (p *parser) checkOperators() (string, error) {
	var oldPosition = p.charIterator
	operator, _ := p.check(operators)
	if p.charIterator != oldPosition && p.charIterator == p.lenInput {
		return "", fmt.Errorf("Input cannot end with an operator")
	} else if p.charIterator != 0 && p.charIterator == oldPosition && p.charIterator < p.lenInput {
		return "", fmt.Errorf("Missing or hanging operator")
	}
	return operator, nil
}

func (p parser) getQuote() string {
	if p.charIterator < p.lenInput {
		switch p.input[p.charIterator] {
		case '"':
			return "\""
		case '\'':
			return "'"
		case '`':
			return "`"
		}
	}
	return ""
}

func (p parser) validateCommands() error {
	if p.numOpenBrackets != 0 {
		return fmt.Errorf("Mismatch of brackets")
	}
	controlLevel := 0 // 0=no input, 1=input but no brackets 2+=input and brackets
	for _, command := range p.outputCommands {
		if command.Control != "" && !isMainCommand(command) && controlLevel < 2 {
			return fmt.Errorf("Conditional: '%s' will not have any input", command.GetConditionalString())
		} else if isMainCommand(command) && controlLevel == 0 {
			controlLevel++
		}
		if command.Bracket == "(" && controlLevel > 0 {
			controlLevel++
		} else if command.Bracket == ")" && controlLevel > 0 {
			controlLevel--
		}
	}
	return nil
}

func orderCheckList(checkList []string) []string {
	orderedCheckList := append([]string{}, checkList...)
	sort.SliceStable(orderedCheckList, func(i, j int) bool {
		return len(orderedCheckList[i]) > len(orderedCheckList[j])
	})
	return orderedCheckList
}

func (p parser) createError(err string) error {
	return fmt.Errorf("%s\n%s", strings.Repeat(" ", p.charIterator)+"^", err)
}
