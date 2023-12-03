package commands

import (
	"github.com/fredjeck/configserver/internal/repository"
	"log/slog"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/encryption"
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
	keystore              *encryption.Keystore
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

	if lastError != nil {
		slog.Error("ConfigServer was not able to start due to missing or invalid configuration file", "err", lastError)
		os.Exit(1)
	}

	configuration.LogEnvironment()

	keystore, lastError = encryption.LoadKeyStoreFromPath(configuration.CertsLocation)
	if lastError != nil {
		slog.Error("ConfigServer was not able to load its keystore", "err", lastError)
		os.Exit(1)
	}

	mgr, lastError := repository.NewManager(configuration.GitConfiguration)
	mgr.Start()
	if lastError != nil {
		slog.Error("ConfigServer was not able to start its GIT repository service", "err", lastError)
		os.Exit(1)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	slog.Info("Starting ConfigServer ...")
	server.New(configuration, keystore).Start()
}
