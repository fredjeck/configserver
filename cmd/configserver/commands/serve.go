package commands

import (
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/repository"
	"github.com/fredjeck/configserver/internal/server"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	configServerCommand.Flags().StringVarP(&configurationFilePath, "configuration", "c", "", "Configuration file location (default is $CONFIGSERVER_HOME/configserver.yml or /var/run/configserver/configserver.yml)")
	rootCmd.AddCommand(configServerCommand)
}

var (
	configServerCommand = &cobra.Command{
		Use:   "serve",
		Short: "Securely distribute your configuration files with style",
		Long:  `Configserver is a git based distributed configuration file system, tailor made with love for cloud based systems`,
		Run: func(cmd *cobra.Command, args []string) {
			initServer()
			slog.Info("Starting ConfigServer ...")
			server.New(configuration, repositoryManager, vault).Start()
		},
	}
	configurationFilePath string
	configuration         *config.Configuration
	repositoryManager     *repository.Manager
	vault                 *encryption.KeyVault
	lastError             error
)

// initialize loads the configuration and prepares configserver for running
func initServer() {
	config.InitLogging()
	configuration, lastError = config.LoadFrom(configurationFilePath)

	if lastError != nil {
		slog.Error("configServer was not able to start due to missing or invalid configuration file", "error", lastError)
		os.Exit(1)
	}

	configuration.LogEnvironment()

	vault, lastError = encryption.LoadKeyVault(configuration.CertsLocation, true)
	if lastError != nil {
		slog.Error("configServer was not able to load its keyvault", "error", lastError)
		os.Exit(1)
	}

	repositoryManager, lastError = repository.NewManager(configuration.GitConfiguration)
	if lastError != nil {
		slog.Error("configServer was not able to start its GIT repository service", "error", lastError)
		os.Exit(1)
	}
	repositoryManager.Start()
}
