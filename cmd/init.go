package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	initDesc := "Initialize development environment"

	initLongDesc := `init - Initialize development environment

"init" will generate a file called "envctl.yaml".

This file has sane defaults, but might need to be edited, and should be checked
into version control.
`

	tpl := `---
image: ubuntu:latest

shell: /bin/bash

bootstrap:
- echo 'Environment initialized' > /envctl

variables:
  FOO: bar
`

	runInit := func(cmd *cobra.Command, args []string) {
		fmt.Println("creating config file... ")

		if _, err := os.Stat(cfgFile); err == nil {
			fmt.Printf("cannot overwrite %v\n", cfgFile)
			os.Exit(1)
		}

		f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("error opening %v: %v\n", cfgFile, err)
			os.Exit(1)
		}

		_, err = f.WriteString(tpl)
		if err != nil {
			fmt.Printf("error writing %v: %v\n", cfgFile, err)
			os.Exit(1)
		}
	}

	return &cobra.Command{
		Use:   "init",
		Short: initDesc,
		Long:  initLongDesc,
		Run:   runInit,
	}
}
