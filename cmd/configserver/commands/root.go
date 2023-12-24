package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "configserver",
	Short: "Securely distribute your configuration files with style",
	Long: `Configserver is a git based distributed configuration file system, 
				tailor made with love for cloud based systems
				find out more on https://github.com/fredjeck/configserver`,
}

func Execute() {
	defCmd := "serve"
	var cmdFound bool
	cmd := rootCmd.Commands()

	for _, a := range cmd {
		for _, b := range os.Args[1:] {
			if a.Name() == b {
				cmdFound = true
				break
			}
		}
	}
	if !cmdFound {
		args := append([]string{defCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	err := rootCmd.Execute()
	dieIfError(err, "Unexpected error")
}

// dieIfError exits the program in case of error and takes care of logging the failure's details
func dieIfError(err error, message string, args ...interface{}) {
	m := message
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	if err != nil {
		fmt.Print(fmt.Errorf("%s: %w", m, err))
		os.Exit(1)
	}
}
