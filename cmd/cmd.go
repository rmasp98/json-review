package cmd

import (
	"fmt"
	"io/ioutil"
	"kube-review/nodelist"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeFile       string
	kubeconfigFile string
	kubeContext    string
	rootCmd        = &cobra.Command{
		Use:   "kube-review",
		Short: "A review tool for kubernetes cluster config",
		Long: "This tool helps to view and search through the kubernetes" +
			"cluster configuration and output results of that search for help" +
			"with reporting. There are also a list of built in search commands" +
			"that will identify common vulnerabilities",
	}
	offlineCmd = &cobra.Command{
		Use:   "offline",
		Short: "Print kubectl command for offline review",
		Long: "This command will print the full kubectl command required to" +
			" create an offline file that can be reviewed at a later date." +
			" You can run 'kube-review offline | xclip -sel clip' to copy the" +
			"command to your clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`kubectl get $(kubectl api-resources --verbs=list -o name | grep -v -e "secrets" -e "componentstatuses" -e "priorityclass" -e "events" | paste -sd, -) --ignore-not-found --all-namespaces -o json > offline.json`)
		},
	}
)

// Execute trigger cobra to run
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(offlineCmd)
	rootCmd.PersistentFlags().StringVarP(&kubeFile, "file", "f", "", "Cluster config file")
	rootCmd.PersistentFlags().StringVar(&kubeconfigFile, "kubeconfig", "", "Path to the kubeconfig file")
	rootCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "The name of the kubeconfig context to use")
}

func getConfig() *nodelist.NodeList {
	var rawJSON []byte
	if kubeFile != "" {
		rawJSON = loadFromFile()
	} else {
		// Run kubectl get ... (after asking permission and outputing what cluster it will run on)
		fmt.Println("I have not yet implemented this so please use flag 'file'")
		os.Exit(1)
	}
	return getNodeList(rawJSON)
}

func loadFromFile() []byte {
	rawJSON, err := ioutil.ReadFile(kubeFile)
	if err != nil {
		fmt.Printf("'%s' does not exist\n", kubeFile)
		os.Exit(1)
	}
	return rawJSON
}

func getNodeList(rawJSON []byte) *nodelist.NodeList {
	if len(rawJSON) > 500000 {
		fmt.Println("This is a large file. Loading may take a few seconds...")
	}
	jsonData, err := nodelist.NewNodeList(rawJSON, true)
	if err != nil {
		fmt.Println("Could not parse JSON data. Maybe an error in the file?")
		os.Exit(1)
	}
	return &jsonData
}
