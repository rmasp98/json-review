package main

import (
	"kube-review/cmd"
	"log"
	"os"
)

func main() {
	f, _ := os.OpenFile("kube-review-debug-log", os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)

	cmd.Execute()
}
