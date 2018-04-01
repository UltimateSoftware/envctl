package cmd

import (
	"fmt"
	"os"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/spf13/cobra"
)

func newLoginCmd(ctl container.Controller, s db.Store) *cobra.Command {
	loginDesc := "log in to the current environment"

	loginLongDesc := `login - Log in to the current environment

"login" will log in to the current environment using the shell specified in
the config file.`

	msgEnvOff := `Wait! The environment isn't ready yet!

To get it ready, run "envctl create".
`

	runLogin := func(cmd *cobra.Command, args []string) {
		env, err := s.Read()
		if err != nil {
			fmt.Printf("error reading data store: %v\n", err)
			os.Exit(1)
		}

		if !env.Initialized() {
			fmt.Print(msgEnvOff)
			os.Exit(1)
		}

		if err := ctl.Attach(env.Container); err != nil {
			fmt.Printf("error logging in to environment: %v\n", err)
			os.Exit(1)
		}
	}

	return &cobra.Command{
		Use:   "login",
		Short: loginDesc,
		Long:  loginLongDesc,
		Run:   runLogin,
	}
}
