package cmd

import (
	"fmt"
	"os"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/UltimateSoftware/envctl/pkg/container/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile = "envctl.yaml"

var rootDesc = "Control your development environments"

var rootLongDesc = `envctl - Control your development environments

A common pattern is to have some sort of tool like Vagrant or Docker to simulate
or mimic production environments on developer workstations. There are _many_
ways to skin this cat.

envctl is a tool for easily controlling these environments. The only thing it
needs is a configuration file, "envctl.yaml", for it to know what to do.
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envctl",
	Short: rootDesc,
	Long:  rootLongDesc,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("image: %v\n", viper.GetString("image"))

	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	ctl := initCtl()
	s := initStore()
	cfg := initConfig()

	rootCmd.AddCommand(newCreateCmd(ctl, s, cfg))
	rootCmd.AddCommand(newDestroyCmd(ctl, s))
	rootCmd.AddCommand(newStatusCmd(s))
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newLoginCmd(ctl, s))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() *viper.Viper {
	cfg := viper.New()

	cfg.SetConfigFile(cfgFile)

	// If a config file is found, read it in. Swallow the error because
	// if anything, the config file will be created with `init`.
	cfg.ReadInConfig()

	return cfg
}

func initStore() db.Store {
	var err error
	jsonStore, err := db.NewJSONStore(".envctl/")
	if err != nil {
		fmt.Printf("error creating environment store: %v\n", err)
		os.Exit(1)
	}

	return jsonStore
}

func initCtl() container.Controller {
	var err error
	ctl, err := docker.NewController()
	if err != nil {
		fmt.Printf("error creating Docker controller: %v\n", err)
		os.Exit(1)
	}

	return ctl
}
