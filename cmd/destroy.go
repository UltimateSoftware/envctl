package cmd

import (
	"fmt"
	"os"

	"github.com/UltimateSoftware/envctl/internal/print"
	"github.com/spf13/cobra"
)

var destroyDesc = "destroy an instance of a development environment"

var destroyLongDesc = `destroy - Destroy an instance of a development environment
`

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: destroyDesc,
	Long:  destroyLongDesc,
	Run:   runDestroy,
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func runDestroy(cmd *cobra.Command, args []string) {
	env, err := jsonStore.Read()
	if err != nil {
		fmt.Printf("error reading data store: %v\n", err)
		os.Exit(1)
	}

	msgEnvOff := `The environment is off!

To create it, run "envctl create".`

	if !env.Initialized() {
		fmt.Println(msgEnvOff)
		jsonStore.Delete() // this will have been set up in initStore
		os.Exit(1)
	}

	fmt.Print("destroying environment... ")

	if err := dockerClient.RemoveContainer(env.Container); err != nil {
		print.Error()
		fmt.Printf("error destroying environment: %v\n", err)
		os.Exit(1)
	}

	if err := dockerClient.RemoveImage(env.Image); err != nil {
		print.Error()
		fmt.Printf("error destroying environment: %v\n", err)
		os.Exit(1)
	}

	if err := jsonStore.Delete(); err != nil {
		print.Error()
		fmt.Printf("error deleting data store: %v\n", err)
		os.Exit(1)
	}

	print.OK()
}
