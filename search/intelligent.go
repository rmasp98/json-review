package search

import (
	"fmt"
	"strings"
)

// GetHints stuff
func GetHints(input string) string {
	return ""
}

// Intelligent stuff
type Intelligent struct {
	input           string
	charIterator    int
	lenInput        int
	numOpenBrackets int
	commands        []Command
	commandIndex    int
}

// NewIntelligent stuff
func NewIntelligent(input string) (Intelligent, error) {
	trimmedInput := strings.Trim(input, " ")
	intelligent := Intelligent{trimmedInput, 0, len(trimmedInput), 0, []Command{}, 0}
	if err := intelligent.validate(); err != nil {
		return Intelligent{}, err
	}
	return intelligent, nil
}

// Execute stuff
func (i *Intelligent) Execute(nodeList sNodeList) {
	filterNodes := i.executeCommands([]int{}, nodeList)
	nodeList.ApplyFilter(filterNodes)
}

// GetCommands returns commands (needed for testing)
func (i Intelligent) GetCommands() []Command {
	return i.commands
}

func (i *Intelligent) executeCommands(currentIndices []int, nodeList sNodeList) []int {
	var nodeIndices []int
	for index := i.commandIndex; index < len(i.commands); index++ {
		command := i.commands[index]
		list := command.RunConitional(currentIndices, nodeList)
		if command.HasOpenBracket() {
			i.commandIndex = index + 1
			list = i.executeCommands(nodeIndices, nodeList)
			index = i.commandIndex
		}
		if command.Operator != "" {
			nodeIndices = command.RunOperation(nodeIndices, list)
		} else {
			nodeIndices = append(nodeIndices, list...)
		}
		if command.HasCloseBracket() {
			i.commandIndex = index
			return nodeIndices
		}
	}

	return nodeIndices
}

func (i *Intelligent) validate() error {
	var err error
	var currentCommand Command
	for i.charIterator < i.lenInput {
		if currentCommand.Operator, err = i.checkOperators(); err != nil {
			return i.createError(err.Error())
		}
		i.checkBrackets('(', 1, &currentCommand)
		if currentCommand.Control, err = i.check(controls); err != nil {
			return i.createError("A control is missing or invalid")
		}
		var condition string
		if condition, err = i.check(conditionals); err != nil {
			return i.createError("A conditional is missing or invalid")
		}
		currentCommand.Equal = (condition == "==")
		if currentCommand.Regex, err = i.checkQuotes(); err != nil {
			return i.createError(err.Error())
		}
		i.checkBrackets(')', -1, &currentCommand)
		if currentCommand.Control != "" {
			i.commands = append(i.commands, currentCommand)
			currentCommand = Command{"", false, "", "", ""}
		}
	}

	if i.numOpenBrackets != 0 {
		return fmt.Errorf("Mismatch of brackets")
	}
	return nil
}

func (i *Intelligent) stripLeft(num int) {
	i.charIterator += num
	for i.charIterator < i.lenInput && i.input[i.charIterator] == ' ' {
		i.charIterator++
	}
}

func (i *Intelligent) check(checkList []string) (string, error) {
	for _, check := range checkList {
		var checkEnd = i.charIterator + len(check)
		if checkEnd < i.lenInput && i.input[i.charIterator:checkEnd] == check {
			i.stripLeft(len(check))
			return check, nil
		}
	}
	return "", fmt.Errorf("")
}

func (i *Intelligent) checkQuotes() (string, error) {
	if quote := i.getQuote(); quote != "" {
		substring := strings.Split(i.input[i.charIterator+1:], quote)
		// If quote at end of string, len(subs) would be at least 1
		if len(substring) > 1 {
			i.stripLeft(len(substring[0]) + 2)
			return substring[0], nil
		}
		return "", fmt.Errorf("Missing end quote")
	}
	return "", fmt.Errorf("Regex is missing or has not been quoted")
}

func (i *Intelligent) checkBrackets(bracket byte, incrementor int, command *Command) {
	for i.charIterator < i.lenInput && i.input[i.charIterator] == bracket {
		command.Bracket = string(bracket)
		i.commands = append(i.commands, *command)
		*command = Command{"", false, "", "", ""}
		i.numOpenBrackets = i.numOpenBrackets + incrementor
		i.stripLeft(1)
	}
}

func (i *Intelligent) checkOperators() (string, error) {
	var oldPosition = i.charIterator
	operator, _ := i.check(operators)
	if i.charIterator != oldPosition && i.charIterator == i.lenInput {
		return "", fmt.Errorf("Input cannot end with an operator")
	} else if i.charIterator != 0 && i.charIterator == oldPosition && i.charIterator < i.lenInput {
		return "", fmt.Errorf("Missing or hanging operator")
	}
	return operator, nil
}

func (i Intelligent) getQuote() string {
	if i.charIterator < i.lenInput {
		switch i.input[i.charIterator] {
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

func (i Intelligent) createError(err string) error {
	return fmt.Errorf("%s\n%s\n%s", i.input, strings.Repeat(" ", i.charIterator)+"^", err)
}
