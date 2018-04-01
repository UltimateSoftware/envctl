package cmd

// var initDesc = "Initialize development environment"

// var initLongDesc = `init - Initialize development environment

// "init" will generate a file called "envctl.yml" (or whatever was passed into
// --config) in the current directory.

// This file has sane defaults, but might need to be edited, and should be checked
// into version control.
// `

// // initCmd represents the init command
// var initCmd = &cobra.Command{
// 	Use:   "init",
// 	Short: initDesc,
// 	Long:  initLongDesc,
// 	Run:   runInit,
// }

// func init() {
// 	rootCmd.AddCommand(initCmd)
// }

// func runInit(cmd *cobra.Command, args []string) {
// 	tpl := `---
// image: ubuntu:latest

// shell: /bin/bash

// bootstrap:
// - echo 'Environment initialized' > /envctl

// variables:
// - FOO=bar
// `

// 	fmt.Print("creating config file... ")

// 	if _, err := os.Stat(cfgFile); err == nil {
// 		print.Error()
// 		fmt.Printf("cannot overwrite %v\n", cfgFile)
// 		os.Exit(1)
// 	}

// 	f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		print.Error()
// 		fmt.Printf("error opening %v: %v\n", cfgFile, err)
// 		os.Exit(1)
// 	}

// 	_, err = f.WriteString(tpl)
// 	if err != nil {
// 		print.Error()
// 		fmt.Printf("error writing %v: %v\n", cfgFile, err)
// 		os.Exit(1)
// 	}

// 	print.OK()
// }
