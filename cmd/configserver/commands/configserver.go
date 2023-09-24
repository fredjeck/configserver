package commands

import (
	"log/slog"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server"
	"github.com/spf13/cobra"
)

var (
	ConfigServerCommand = &cobra.Command{
		Use:   "configserver",
		Short: "Securely distribute your configuration files with style",
		Long:  `Configserver is a git based distributed configuration file system, tailor made with a lot of love for cloud based systems`,
		Run:   startServer,
	}
	configurationFilePath string
	configuration         *config.Configuration
	lastError             error
)

func init() {
	cobra.OnInitialize(initialize)
	ConfigServerCommand.PersistentFlags().StringVarP(&configurationFilePath, "configuration", "c", "", "Configuration file location (default is $CONFIGSERVER_HOME/configserver.yml or /var/run/configserver/configserver.yml)")
}

func Run(args []string) error {
	ConfigServerCommand.SetArgs(args)
	return ConfigServerCommand.Execute()
}

// initialize loads the configuration and prepares configserver for running
func initialize() {
	config.InitLogging()
	configuration, lastError = config.LoadFrom(configurationFilePath)
	slog.Info("Environment",
		config.EnvConfigServerEnvironment, configuration.Environment.Kind,
	)

	if lastError != nil {
		slog.Error("ConfigServer was not able to start due to missing or invalid configuration file", "err", lastError)
		os.Exit(1)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	slog.Info("Starting ConfigServer ...")
	server.New(configuration).Start()
}
