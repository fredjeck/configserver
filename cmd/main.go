package main

import (
	"github.com/fredjeck/configserver/cmd/configserver/commands"
	"os"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
