package commands

import (
	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	home          string
	Logger        *zap.Logger
	Configuration *config.Config
	RootCommand   = &cobra.Command{
		Use:   "configserver",
		Short: "Externalize our configuration in distributed systems",
		Long:  `configserver allows you to server your configuration files from git repositories with style`,
	}
)

func Run(args []string) error {
	RootCommand.SetArgs(args)
	return RootCommand.Execute()
}

// Execute executes the root command.
func Execute() error {
	return RootCommand.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCommand.PersistentFlags().StringVarP(&home, "configuration", "c", "", "configuration files location (default is $CONFIGSERVER_HOME or /var/run/configserver)")
}

func initConfig() {
	Logger = logging.NewLogger()
	var err error
	Configuration, err = config.ReadFromPath(home)

	if err != nil {
		Logger.Sugar().Fatal(err)
	}

	Logger.Sugar().Infof("Configuration loaded from '%s'", Configuration.LoadedFrom)
}
