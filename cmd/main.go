// main is the main executable for ConfigServer
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/fredjeck/configserver/internal/server"
)

var configurationPath string

func init() {
	const (
		defaultConfiguration = "/var/run/configserver/configserver.yml"
		configurationUsage   = "path to the configuration file"
	)
	flag.StringVar(&configurationPath, "configuration", defaultConfiguration, configurationUsage)
	flag.StringVar(&configurationPath, "c", defaultConfiguration, configurationUsage+" (shorthand)")

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, `Usage:
configserver -c /path/to/configuration.yml

Starts a new configserver instance using the provided configuration. 
If the configuration is omitted, attempts to locate the configuration in the folder pointed by the CONFIGSERVER_HOME environment variable 
`)

		flag.PrintDefaults()
	}
}

func main() {
	configuration.InitLogging()
	flag.Parse()
	c, err := configuration.LoadFrom(configurationPath)
	if err != nil {
		slog.Error("Configuration cannot be loaded, exiting ...", "error", err)
		os.Exit(1)
	}
	c.LogEnvironment()
	s := server.NewConfigServer(c)
	s.Start()
}
