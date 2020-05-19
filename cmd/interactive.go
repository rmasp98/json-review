package cmd

import (
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
	nodeList := getConfig()
	ui.Run(nodeList)
}
