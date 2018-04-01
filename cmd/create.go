package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/internal/print"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateCmd(ctl container.Controller, s db.Store) *cobra.Command {
	createDesc := "create a new instance of a development environment"
	createLongDesc := `create - Create an instance of a development environment

"create" will dynamically build a development environment based on the settings
in the config file. Only one environment can exist at any time per config file.
`

	msgEnvReady := `There is already an environment ready for use!

To use it, run "envctl login", or destroy it with "envctl destroy".`

	runCreate := func(cmd *cobra.Command, args []string) {
		env, err := s.Read()
		if err != nil {
			fmt.Printf("error reading environment state: %v\n", err)
			os.Exit(1)
		}

		if env.Initialized() {
			fmt.Println(msgEnvReady)
			os.Exit(1)
		}

		name := uuid.New().String()
		baseImage := viper.GetString("image")
		shell := viper.GetString("shell")
		mount := viper.GetString("mount")
		rawenvs := viper.GetStringSlice("variables")

		if mount == "" {
			fmt.Println("no mount specified, defaulting to /mnt/repo... [ WARN ]")
			mount = "/mnt/repo"
		}

		// This supports dynamic evaluation of environment variables so secrets
		// don't have to be checked into the repo, but config files don't have
		// to be generated from templates either.
		envs := make([]string, len(rawenvs))
		for i, rawenv := range rawenvs {
			s := strings.Split(rawenv, "=")
			k, v := s[0], s[1]

			if v[0] == '$' {
				v = os.Getenv(v[1:])
			}

			envs[i] = fmt.Sprintf("%v=%v", k, v)
		}

		pwd, err := os.Getwd()
		if err != nil {
			print.Error()
			fmt.Printf("error getting current working directory: %v\n", err)
			os.Exit(1)
		}

		meta := container.Metadata{
			BaseName:  name,
			BaseImage: baseImage,
			Shell:     shell,
			Mount: container.Mount{
				Source:      pwd,
				Destination: mount,
			},
			Envs: envs,
		}

		fmt.Print("creating your environment... ")

		newMeta, err := ctl.Create(meta)
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
					s.Create(db.Environment{
						Status:    db.StatusError,
						Container: newMeta,
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
				s.Create(db.Environment{
					Status:    db.StatusError,
					Container: newMeta,
				})
				os.Exit(1)
			}

			cmdarr := []string{shell, fname}

			err = ctl.Run(newMeta, cmdarr)
			if err != nil {
				print.Error()
				fmt.Printf("error running %v: %v\n", cmdarr, err)
				s.Create(db.Environment{
					Status:    db.StatusError,
					Container: newMeta,
				})
				os.Exit(1)
			}
			print.OK()
		}

		fmt.Print("saving environment... ")
		err = s.Create(db.Environment{
			Status:    db.StatusReady,
			Container: newMeta,
		})
		if err != nil {
			print.Error()
			fmt.Printf("error saving environment: %v\n", err)
			os.Exit(1)
		}
		print.OK()
	}

	return &cobra.Command{
		Use:   "create",
		Short: createDesc,
		Long:  createLongDesc,
		Run:   runCreate,
	}
}
