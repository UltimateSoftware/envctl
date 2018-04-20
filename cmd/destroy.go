package cmd

import (
	"fmt"
	"os"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/spf13/cobra"
)

func newDestroyCmd(ctl container.Controller, s db.Store) *cobra.Command {
	destroyDesc := "destroy an instance of a development environment"
	destroyLongDesc := `destroy - Destroy an instance of a development environment
`

	msgEnvOff := `The environment is off!

To create it, run "envctl create".`

	runDestroy := func(cmd *cobra.Command, args []string) {
		env, err := s.Read()
		if err != nil {
			fmt.Printf("error reading data store: %v\n", err)
			os.Exit(1)
		}

		if !env.Initialized() {
			fmt.Println(msgEnvOff)
			s.Delete()
			os.Exit(1)
		}

		fmt.Println("destroying environment... ")

		if err := ctl.Remove(env.Container); err != nil {
			fmt.Printf("error destroying environment: %v\n", err)
			os.Exit(1)
		}

		if err := s.Delete(); err != nil {
			fmt.Printf("error deleting data store: %v\n", err)
			os.Exit(1)
		}
	}

	return &cobra.Command{
		Use:   "destroy",
		Short: destroyDesc,
		Long:  destroyLongDesc,
		Run:   runDestroy,
	}
}
