package commands

import (
	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/fredjeck/configserver/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	home          string
	Configuration *config.Config
	Key           *[32]byte
	RootCommand   = &cobra.Command{
		Use:   "configserver",
		Short: "Externalize your configuration in distributed systems",
		Long:  `Configserver allows you to serve your configuration files from git repositories with style`,
		Run:   startServer,
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
	var err error

	env := strings.ToLower(os.Getenv("CONFIGSERVER_ENV"))
	if strings.Contains(env, "dev") {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
	zap.L().Sugar().Infof("CONFIGSERVER_ENV='%s'", env)

	Configuration, err = config.ReadFromPath(home)

	if err != nil {
		zap.L().Sugar().Fatal("Cannot read the configuration file", err)
	}

	zap.L().Sugar().Infof("Configuration loaded from '%s'", Configuration.LoadedFrom)

	Key, err = encrypt.ReadEncryptionKey(Configuration.EncryptionKeyPath(), true)
	if err != nil {
		zap.L().Sugar().Fatal(err)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	srv := server.New(Configuration, Key)
	srv.Start()
}
