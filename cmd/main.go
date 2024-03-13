package main

import (
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server"
)

func main() {
	c := config.DefaultConfiguration
	c.Environment.Kind = "development"
	c.Repositories.CheckoutLocation = "/tmp/configserver"
	c.Repositories.Configuration = []*config.Repository{&config.Repository{
		Name:                   "go-configserver",
		Url:                    "https://github.com/fredjeck/go-configserver",
		RefreshIntervalSeconds: 3600,
		Clients: []string{
			"myclientid",
			"sample_client",
		},
	}}
	config.InitLogging()
	c.LogEnvironment()
	s := server.NewConfigServer(c)
	s.Start()
}
