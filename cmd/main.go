package main

import (
	config "github.com/fredjeck/configserver/internal/configuration"
	"github.com/fredjeck/configserver/internal/server"
)

func main() {
	c := config.DefaultConfiguration
	c.Environment.Kind = "development"
	c.Server.PassPhrase = "To infinity and beyond"
	c.Repositories.CheckoutLocation = "/tmp/configserver"
	c.Repositories.Configuration = []*config.Repository{
		{
			Name: "configserver-samples-integration",
			Url:  "https://github.com/fredjeck/configserver-samples",
			//Branch:                 "integration",
			RefreshIntervalSeconds: 3600,
			Clients: []string{
				"myclientid",
				"sample_client",
			},
		},
	}
	config.InitLogging()
	c.LogEnvironment()
	s := server.NewConfigServer(c)
	s.Start()
}
