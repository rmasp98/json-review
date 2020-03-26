package main

import (
	"fmt"
	"io/ioutil"
	"kube-review/jsontree"
	"kube-review/ui"
	"log"
	"os"
)

func main() {
	f, _ := os.OpenFile("kube-review-debug-log", os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)

	fmt.Println("Hello World!")

	if len(os.Args) > 1 {
		var jsonData jsontree.NodeList
		var errJSON error
		{
			rawJSON, errRead := ioutil.ReadFile(os.Args[1])
			if errRead != nil {
				fmt.Printf("%s does not exist\n", os.Args[1])
				os.Exit(1)
			}
			// jsonData, errJSON = jsontree.CreateTree(string(rawJSON))
			// if errRead != nil || errJSON != nil {
			// 	fmt.Printf("%s does not contain correct JSON\n", os.Args[1])
			// 	os.Exit(1)
			// }
			jsonData, errJSON = jsontree.NewNodeList(string(rawJSON))
			if errRead != nil || errJSON != nil {
				fmt.Printf("%s does not contain correct JSON\n", os.Args[1])
				fmt.Print(errRead, errJSON)
				os.Exit(1)
			}
		}

		cui := ui.NewCursesUI(&jsonData)
		cui.Run()
	}
}
