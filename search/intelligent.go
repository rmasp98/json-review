package search

// GetHints stuff
func GetHints(input string, cursorPos int) string {
	// starting at cursorPos go back until we find ( or operator
	// If substring contains quote or condition return nothing
	//Otherwise strip spaces and return any matching controls
	return ""
}

// Intelligent stuff
type Intelligent struct {
	commands     []Command
	commandIndex int
}

// NewIntelligent stuff
func NewIntelligent(input string) (Intelligent, error) {
	commands, err := Parse(input)
	if err != nil {
		return Intelligent{}, err
	}
	return Intelligent{commands, 0}, nil
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
