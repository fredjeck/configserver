package commands

import (
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ConfigServerCommand = &cobra.Command{
		Use:   "configserver",
		Short: "Securely distribute your configuration files with style",
		Long:  `Configserver is a git based distributed configuration file system, tailor made with a lot of love for cloud based systems`,
		Run:   startServer,
	}
	configurationFilePath string
	Configuration         *config.Configuration
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

func initialize() {
	config.InitLogging()
	Configuration, lastError = config.LoadFrom(configurationFilePath)
	if lastError != nil {
		zap.L().Sugar().Fatal("ConfigServer was not able to start due to missing or invalid configuration file: ", lastError)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	zap.L().Sugar().Infof("Starting ConfigServer ...")
	server.New(Configuration)
}
