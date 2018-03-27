package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/internal/print"
	"github.com/UltimateSoftware/envctl/pkg/docker"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDesc = "create a new instance of a development environment"

var createLongDesc = `create - Create an instance of a development environment

"create" will dynamically build a development environment based on the settings
in the config file. Only one environment can exist at any time per config file.
`

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: createDesc,
	Long:  createLongDesc,
	Run:   runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) {
	env, err := jsonStore.Read()
	if err != nil {
		fmt.Printf("error reading environment state: %v\n", err)
		os.Exit(1)
	}

	msgEnvReady := `There is already an environment ready for use!

To use it, run "envctl login", or destroy it with "envctl destroy".`

	if env.Initialized() {
		fmt.Println(msgEnvReady)
		os.Exit(1)
	}

	name := uuid.New().String()
	baseImage := viper.GetString("image")
	shell := viper.GetString("shell")
	mount := viper.GetString("mount")

	if mount == "" {
		fmt.Println("no mount specified, defaulting to /mnt/repo... [ WARN ]")
		mount = "/mnt/repo"
	}

	fmt.Print("creating your environment... ")

	img, err := dockerClient.BuildImage(docker.ImageConfig{
		BaseName:  name,
		Shell:     shell,
		BaseImage: baseImage,
		Mount:     mount,
	})
	if err != nil {
		print.Error()
		fmt.Printf("error building environment: %v\n", err)
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		print.Error()
		fmt.Printf("error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	rawenvs := viper.GetStringSlice("variables")

	ccfg := docker.ContainerConfig{
		Name:      name,
		ImageName: img,
		Mounts: []docker.Mount{
			docker.Mount{
				Source:      pwd,
				Destination: mount,
			},
		},
		Env: make([]string, len(rawenvs)),
	}

	// This supports dynamic evaluation of environment variables so secrets
	// don't have to be checked into the repo, but config files don't have
	// to be generated from templates either.
	for i, rawenv := range rawenvs {
		s := strings.Split(rawenv, "=")
		k, v := s[0], s[1]

		if v[0] == '$' {
			v = os.Getenv(v[1:])
		}

		ccfg.Env[i] = fmt.Sprintf("%v=%v", k, v)
	}

	cnt, err := dockerClient.CreateContainer(ccfg)
	if err != nil {
		print.Error()
		fmt.Printf("error creating environment: %v\n", err)
		os.Exit(1)
	}

	print.OK()

	rawcmds := viper.GetStringSlice("bootstrap")

	if len(rawcmds) > 0 {
		fmt.Print("running bootstrap steps... ")

		script := &bytes.Buffer{}
		for _, rawcmd := range rawcmds {
			_, err := script.WriteString(fmt.Sprintf("%v\n", rawcmd))
			if err != nil {
				print.Error()
				fmt.Printf("error generating bootstrap script: %v\n", err)
				jsonStore.Create(db.Environment{
					Status:    db.StatusError,
					Container: cnt,
					Image:     img,
				})
				os.Exit(1)
			}
		}

		fname := ".envctl/" + uuid.New().String()
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			print.Error()
			fmt.Printf("error opening tmp script for writing: %v\n", err)
			os.Exit(1)
		}

		if _, err := io.Copy(f, script); err != nil {
			print.Error()
			fmt.Printf("error writing bootstrap script: %v\n", err)
			jsonStore.Create(db.Environment{
				Status:    db.StatusError,
				Container: cnt,
				Image:     img,
			})
			os.Exit(1)
		}

		cmdarr := []string{shell, fname}

		err = dockerClient.RunOnContainer(cmdarr, cnt)
		if err != nil {
			print.Error()
			fmt.Printf("error running %v: %v\n", cmdarr, err)
			jsonStore.Create(db.Environment{
				Status:    db.StatusError,
				Container: cnt,
				Image:     img,
			})
			os.Exit(1)
		}
		print.OK()
	}

	fmt.Print("saving environment... ")
	err = jsonStore.Create(db.Environment{
		Status:    db.StatusReady,
		Image:     img,
		Container: cnt,
	})
	if err != nil {
		print.Error()
		fmt.Printf("error saving environment: %v\n", err)
		os.Exit(1)
	}
	print.OK()
}
