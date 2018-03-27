package cmd

import (
	"fmt"
	"os"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var jsonStore *db.JSONStore

var dockerClient *docker.Client

var rootDesc = "Control your development environments"

var rootLongDesc = `envctl - Control your development environments

A common pattern is to have some sort of tool like Vagrant or Docker to simulate
or mimic production environments on developer workstations. There are _many_
ways to skin this cat.

envctl is a tool for easily controlling these environments. The only thing it
needs is a configuration file, "envctl.yml", for it to know what to do.
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
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initStore)
	cobra.OnInitialize(initDocker)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is envctl.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		cfgFile = "envctl.yaml"
	}

	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in. Swallow the error because
	// if anything, the config file will be created with `init`.
	viper.ReadInConfig()
}

func initStore() {
	var err error
	jsonStore, err = db.NewJSONStore(".envctl/")
	if err != nil {
		fmt.Printf("error creating environment store: %v\n", err)
		os.Exit(1)
	}
}

func initDocker() {
	var err error
	dockerClient, err = docker.NewClient()
	if err != nil {
		fmt.Printf("error creating docker client: %v\n", err)
		os.Exit(1)
	}
}
