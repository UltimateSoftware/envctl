package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var loginDesc = "log in to the current environment"

var loginLongDesc = `login - Log in to the current environment

"login" will log in to the current environment using the shell specified in
the config file.
`

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: loginDesc,
	Long:  loginLongDesc,
	Run:   runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	env, err := jsonStore.Read()
	if err != nil {
		fmt.Printf("error reading data store: %v\n", err)
		os.Exit(1)
	}

	msgEnvOff := `Wait! The environment isn't ready yet!

To get it ready, run "envctl create".
`

	if !env.Initialized() {
		fmt.Print(msgEnvOff)
		os.Exit(1)
	}

	if err := dockerClient.AttachContainer(env.Container); err != nil {
		fmt.Printf("error logging in to environment: %v\n", err)
		os.Exit(1)
	}
}
