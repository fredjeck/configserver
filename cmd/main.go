package main

import (
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server"
)

func main() {
	c := config.DefaultConfiguration
	config.InitLogging()
	c.LogEnvironment()
	s := server.NewConfigServer(c)
	s.Start()
}
