package main

import (
	"context"
	"log"
	"os"

	"github.com/guessi/eks-node-diagnostic/cmd"
)

func main() {
	app := cmd.Entry()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
