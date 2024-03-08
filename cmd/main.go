package main

import "github.com/fredjeck/configserver/server"

func main() {
	c := &server.Configuration{
		PassPhrase:       "This is a passphrase used to protect yourself",
		ListenOn:         "127.0.0.1:4200",
		SecretExpiryDays: 60,
	}
	s := server.NewConfigServer(c)
	s.Start()
}
