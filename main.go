package main

import (
	"log"
	"os"

	"github.com/guessi/eks-node-diagnostic/cmd"
)

func main() {
	app := cmd.Entry()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
