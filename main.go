package main

import (
	"deduplicate/commands"
	"fmt"
	"os"
)

func main() {
	rootCmd := commands.NewCmdRoot()
	if _, err := rootCmd.ExecuteC(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
