package search

import (
	"regexp"
	"strings"
	"unicode"
)

// Expression stuff
type Expression struct {
	commands  []Command
	cmdIndex  int
	variables map[string][]int
}

// GetExpressionHints stuff
func GetExpressionHints(input string) []string {
	trimmedInput := strings.TrimRight(input, " ")
	if len(trimmedInput) > 0 && trimmedInput[len(trimmedInput)-1] == ')' {
		return Operators
	} else if function, arguments := getFunctionType(input); function != CMDNULL {
		focussedArgIndex := getFocussedArgument(arguments, function.template())
		argHints := getArgumentHints(focussedArgIndex, function.template())
		return append(argHints, getFunctionHint(function, focussedArgIndex))
	}
	return getBaseFuntionHints(input)
}

// InsertSelectedExpressionHint stuff
func InsertSelectedExpressionHint(input string, index int) string {
	var output string
	trimmedInput := strings.TrimRight(input, " ")
	if len(trimmedInput) > 0 && trimmedInput[len(trimmedInput)-1] == ')' {
		output = strings.TrimRight(input, " ") + " " + Operators[index] + " "
	} else if function, arguments := getFunctionType(input); function != CMDNULL {
		focussedArgIndex := getFocussedArgument(arguments, function.template())
		argHints := getArgumentHints(focussedArgIndex, function.template())
		if index < len(argHints) {
			input = strings.TrimRightFunc(input, func(c rune) bool {
				return c != ',' && c != '='
			})
			output = input + argHints[index]
		} else {
			output = input
		}
	} else {
		hint := strings.Split(getBaseFuntionHints(input)[index], "(")[0] + "("
		input = strings.TrimRightFunc(input, func(c rune) bool {
			return (c > 97 || c < 122) && (c > 41 || c < 90)
		})
		output = input + hint
	}
	return output
}

// NewExpression stuff
func NewExpression(input string) (Expression, error) {
	commands, err := Parse(input)
	if err != nil {
		return Expression{}, err
	}
	return Expression{commands, 0, map[string][]int{}}, nil
}

// Execute stuff
func (e Expression) Execute(nodeList sNodeList, mode FunctionEnum) error {
	nodes := e.executeCommands(nodeList)
	return nodeList.Filter(nodes)
}

func (e *Expression) executeCommands(nodeList sNodeList) []int {
	var currentIndices []int
	for index := e.cmdIndex; index < len(e.commands); index++ {
		command := e.commands[index]
		input := e.variables[command.GetInputName()]
		outName, output := command.RunFunction(input, nodeList)
		if outName != "" {
			e.variables[outName] = output
		}
		if command.HasOpenBracket() {
			e.cmdIndex = index + 1
			output = e.executeCommands(nodeList)
			index = e.cmdIndex
		}
		currentIndices = command.RunOperation(currentIndices, output)
		if command.HasCloseBracket() {
			e.cmdIndex = index
			return currentIndices
		}
	}
	return currentIndices
}

func getFunctionType(input string) (CmdFunc, []string) {
	bracketIndex := strings.LastIndex(input, "(")
	if bracketIndex > 0 && unicode.IsLetter(rune(input[bracketIndex-1])) {
		arguments := strings.Split(input[bracketIndex+1:], ",")
		if match, _ := regexp.Match("(?i)"+CMDFINDNODES.String()+"$", []byte(input[:bracketIndex])); match {
			return CMDFINDNODES, arguments
		} else if match, _ := regexp.Match("(?i)"+CMDFINDRELATIVE.String()+"$", []byte(input[:bracketIndex])); match {
			return CMDFINDRELATIVE, arguments
		}
	}
	return CMDNULL, []string{}
}

func getArgumentHints(argIndex int, argTemp []argTemplate) []string {
	if argIndex > 0 && argIndex < len(argTemp) {
		arg := argTemp[argIndex]
		switch arg.argType {
		case "MatchType":
			return []string{"ANY", "KEY", "VALUE"}
		case "bool":
			return []string{"true", "false"}
		}
	}
	return []string{}
}

func getBaseFuntionHints(input string) []string {
	strippedInput := strings.TrimFunc(input, func(c rune) bool {
		return (c <= 97 || c >= 122) && (c <= 41 || c >= 90)
	})
	var output []string
	if match, _ := regexp.Match("(?i)"+strippedInput, []byte(CMDFINDNODES.String())); match {
		output = append(output, getFunctionHint(CMDFINDNODES, -1))
	}
	if match, _ := regexp.Match("(?i)"+strippedInput, []byte(CMDFINDRELATIVE.String())); match {
		output = append(output, getFunctionHint(CMDFINDRELATIVE, -1))
	}
	return output
}

func getFunctionHint(function CmdFunc, argNum int) string {
	var out = function.String() + "("
	for index, arg := range function.template() {
		if index == argNum {
			out += redBold + arg.name + " (" + arg.description + ")" + reset + ", "
		} else {
			out += arg.name + ", "
		}
	}
	return strings.TrimRight(out, ", ") + ")"
}

func getFocussedArgument(arguments []string, argTemp []argTemplate) int {
	kwargsActive := false
	for _, arg := range arguments {
		if strings.Contains(arg, "=") {
			kwargsActive = true
		}
	}
	if kwargsActive {
		kwarg := strings.Trim(strings.Split(arguments[len(arguments)-1], "=")[0], " ")
		return getArgIndexByName(kwarg, argTemp)
	}
	return len(arguments) - 1
}
