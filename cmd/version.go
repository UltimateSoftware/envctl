package cmd

import (
	"fmt"

	"github.com/UltimateSoftware/envctl/internal/print"

	"github.com/spf13/cobra"
)

var envctlVersion string

func newVersionCmd() *cobra.Command {
	versionDesc := "get the current version of envctl"
	versionLongDesc := "version - Get the current version of envctl"

	return &cobra.Command{
		Use:   "version",
		Short: versionDesc,
		Long:  versionLongDesc,
		Run: func(cmd *cobra.Command, args []string) {
			if envctlVersion == "" {
				fmt.Print("no version set for this build... ")
				print.Warn()
				envctlVersion = "local"
			}

			fmt.Println(envctlVersion)
		},
	}
}
