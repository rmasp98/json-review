package search

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse stuff
func Parse(input string) ([]Command, error) {
	p := parser{strings.Trim(input, " "), 0, []string{}, Command{}, []Command{}}
	return p.parse()
}

type parser struct {
	input          string
	charIter       int
	outputs        []string
	currentCommand Command
	outputCommands []Command
}

func (p *parser) parse() ([]Command, error) {
	var openBrackets = 0
	for p.charIter < len(p.input) {
		if err := p.checkOperator(); err != nil {
			return []Command{}, err
		}
		openBrackets += p.checkBrackets("(")
		if err := p.checkFunction(); err != nil {
			return []Command{}, err
		}
		openBrackets -= p.checkBrackets(")")
		if p.currentCommand.function != CMDNULL {
			p.outputCommands = append(p.outputCommands, p.currentCommand)
			p.currentCommand = Command{}
		}
	}
	if openBrackets != 0 {
		return []Command{}, fmt.Errorf("No close bracket")
	}
	return p.outputCommands, nil
}

func (p parser) getNextSlice(numChars int) string {
	if numChars == -1 {
		if p.charIter < len(p.input) {
			return p.input[p.charIter:]
		}
	} else if numChars+p.charIter <= len(p.input) {
		return p.input[p.charIter : p.charIter+numChars]
	}
	return ""
}

func (p *parser) stripLeft(numChars int) {
	p.charIter += numChars
	for p.getNextSlice(1) == " " {
		p.charIter++
	}
}

func (p *parser) checkOperator() error {
	for _, operator := range Operators {
		if p.getNextSlice(len(operator)) == operator {
			p.stripLeft(len(operator))
			if len(p.getNextSlice(-1)) == 0 {
				return fmt.Errorf("Hanging operator")
			}
			p.currentCommand.operator = operator
			return nil
		}
	}
	// This is broken for when starts with bracket
	if p.charIter < 8 {
		return nil
	}
	return fmt.Errorf("Missing operator")
}

func (p *parser) checkBrackets(bracket string) int {
	numBrackets := 0
	for p.getNextSlice(1) == bracket {
		p.stripLeft(1)
		numBrackets++

		p.currentCommand.bracket = bracket
		p.outputCommands = append(p.outputCommands, p.currentCommand)
		p.currentCommand = Command{}
	}
	return numBrackets
}

func (p *parser) checkFunction() error {
	if strings.EqualFold(p.getNextSlice(13), CMDFINDRELATIVE.String()+"(") {
		p.stripLeft(13)
		p.currentCommand.function = CMDFINDRELATIVE
		if err := p.parseArguments(p.getNextSlice(-1), findRelArgs); err != nil {
			return err
		}
	} else if strings.EqualFold(p.getNextSlice(10), CMDFINDNODES.String()+"(") {
		p.stripLeft(10)
		p.currentCommand.function = CMDFINDNODES
		if err := p.parseArguments(p.getNextSlice(-1), findArgs); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Invalid function name")
	}
	return nil
}

func (p *parser) parseArguments(args string, template []argTemplate) error {
	arguments := strings.Split(p.getNextSlice(-1), ")")
	if len(arguments) > 1 {
		var argMap = map[string]string{}
		var kwargsActive = false
		for index, arg := range strings.Split(arguments[0], ",") {
			var name string
			var finalArg string
			arg := strings.Trim(arg, " ")
			if arg != "" {
				// Ensures the check for "=" is outside a regex
				if strings.Contains(strings.Split(arg, "\"")[0], "=") {
					kwargsActive = true
					split := strings.Split(arg, "=")
					name = strings.Trim(split[0], " ")
					finalArg = strings.Trim(split[1], " ")
					var argType string
					for _, argTemp := range template {
						if strings.EqualFold(argTemp.name, name) {
							argType = argTemp.argType
						}
					}
					if err := p.validateArgument(finalArg, argType); err != nil {
						return err
					}
				} else if !kwargsActive && index < len(template) {
					if err := p.validateArgument(arg, template[index].argType); err != nil {
						return err
					}
					name = template[index].name
					finalArg = arg
				} else {
					return fmt.Errorf("Not allowed a normal argument after a keyword argument")
				}
			} else {
				return fmt.Errorf("Argument empty")
			}
			if name == "output" {
				p.currentCommand.output = finalArg
			} else {
				if name == "regex" {
					finalArg = strings.Trim(finalArg, "\"")
				}
				argMap[name] = finalArg
			}
		}
		// Strip arguments plus final bracket
		p.stripLeft(len(arguments[0]) + 1)
		p.currentCommand.input = argMap
		return nil
	}
	return fmt.Errorf("No close bracket for function")
}

func (p *parser) validateArgument(argument, argType string) error {
	switch argType {
	case "regex":
		if (argument)[0] == '"' && (argument)[len(argument)-1] == '"' {
			argument := strings.Trim(argument, "\"")
			if _, err := regexp.Compile(argument); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Regex has not been quoted")
		}
	case "MatchType":
		if !strings.EqualFold(argument, "Any") && !strings.EqualFold(argument, "Key") && !strings.EqualFold(argument, "Value") {
			return fmt.Errorf("MatchType invalid")
		}
	case "bool":
		if !strings.EqualFold(argument, "True") && !strings.EqualFold(argument, "False") {
			return fmt.Errorf("Bool invalid")
		}
	case "int":
		if _, err := strconv.Atoi(argument); err != nil {
			return err
		}
	case "input":
		var exists = false
		for _, output := range p.outputs {
			if argument == output {
				exists = true
			}
		}
		if !exists {
			return fmt.Errorf("Input (%s) is not created before being called", argument)
		}
	case "output":
		p.outputs = append(p.outputs, argument)
	default:
		return fmt.Errorf("Invalid argument type: '%s'", argType)
	}
	return nil
}
