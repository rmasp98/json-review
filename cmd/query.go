package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	queryList []string
	queryCmd  = &cobra.Command{
		Use:   "query",
		Short: "Find common issues in config",
		Long:  "This tool will automatically run and output the defined list of search commands",
		Run:   queryRun,
	}
)

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringArrayVarP(&queryList, "queries", "q", []string{}, "List of queries to run")
}

func queryRun(cmd *cobra.Command, args []string) {
	fmt.Println()
}
