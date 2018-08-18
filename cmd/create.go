package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/UltimateSoftware/envctl/internal/config"
	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func newCreateCmd(
	ctl container.Controller,
	s db.Store,
	l config.Loader,
) *cobra.Command {
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

		cfg, err := l.Load()
		if err != nil {
			fmt.Printf("error reading config file: %v\n", err)
			os.Exit(1)
		}

		name := uuid.New().String()
		baseImage := cfg.Image
		shell := cfg.Shell
		mount := cfg.Mount

		if mount == "" {
			fmt.Println("no mount specified, defaulting to /mnt/repo...")
			mount = "/mnt/repo"
		}

		envs, err := parseVariables(cfg)
		if err != nil {
			fmt.Printf("error getting environment variables: %v\n", err)
			os.Exit(1)
		}

		pwd, err := os.Getwd()
		if err != nil {
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
			Envs:    envs,
			NoCache: !(*cfg.CacheImage),
			User:    cfg.User,
			Ports:   cfg.Ports,
		}

		fmt.Println("creating your environment...")

		newMeta, err := ctl.Create(meta)
		if err != nil {
			fmt.Printf("error creating environment: %v\n", err)
			os.Exit(1)
		}

		rawcmds := cfg.Bootstrap
		if len(rawcmds) > 0 {
			fmt.Println("running bootstrap steps...")

			script := &bytes.Buffer{}
			for _, rawcmd := range rawcmds {
				_, err := script.WriteString(fmt.Sprintf("%v\n", rawcmd))
				if err != nil {
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
				fmt.Printf("error opening tmp script for writing: %v\n", err)
				os.Exit(1)
			}

			if _, err := io.Copy(f, script); err != nil {
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
				fmt.Printf("error running %v: %v\n", cmdarr, err)
				s.Create(db.Environment{
					Status:    db.StatusError,
					Container: newMeta,
				})
				os.Exit(1)
			}
		}

		fmt.Println("saving environment...")
		err = s.Create(db.Environment{
			Status:    db.StatusReady,
			Container: newMeta,
		})
		if err != nil {
			fmt.Printf("error saving environment: %v\n", err)
			os.Exit(1)
		}
	}

	return &cobra.Command{
		Use:   "create",
		Short: createDesc,
		Long:  createLongDesc,
		Run:   runCreate,
	}
}

func parseVariables(cfg config.Opts) ([]string, error) {
	rawenvs := cfg.Variables

	// This supports dynamic evaluation of environment variables so secrets
	// don't have to be checked into the repo, but config files don't have
	// to be generated from templates either.
	envs := []string{}
	for k, v := range rawenvs {
		if v[0] == '$' {
			v = os.Getenv(v[1:])
		}

		if v == "" {
			return []string{}, fmt.Errorf("missing variable %v", k)
		}

		envs = append(envs, fmt.Sprintf("%v=%v", k, v))
	}

	return envs, nil
}
