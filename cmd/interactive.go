package cmd

import (
	"fmt"
	"kube-review/search"
	"kube-review/ui"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive session",
	Long:  "This will start an ncurses GUI to view and search the kubernetes config manually",
	Run:   interactiveRun,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func interactiveRun(cmd *cobra.Command, args []string) {
	queryList := search.NewQueryList()
	if err := queryList.Load("querylist.json", "search/queryschema.json"); err != nil {
		fmt.Println("Failed to load 'querylist.json' - " + err.Error())
		return
	}
	nodeList := getConfig()

	ui.Run(nodeList, &queryList)
}
