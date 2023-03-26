package main

import (
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/fredjeck/configserver/encrypt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	command := &cobra.Command{
		Use:   "keygen",
		Short: "Configserver Keygen generates a random key for encrypting sensitive values",
		Long:  `Configserver Keygen generates a random key for encrypting sensitive values`,
		Run: func(cmd *cobra.Command, args []string) {

			key := encrypt.NewEncryptionKey()
			encoded := b64.StdEncoding.EncodeToString(key[:])

			logger.Sugar().Infof("Generated key is (base64) %s", encoded)

			if _, err := os.Stat("/var/run/configserver"); os.IsNotExist(err) {
				err := os.MkdirAll("/var/run/configserver", 0700)
				if err != nil {
					logger.Sugar().Errorf("Cannot create path %s: %s", "/var/run/configserver", err.Error())
					logger.Fatal("Failed to init storage factory", zap.Error(err))
				}

			}

			err := os.WriteFile("/var/run/configserver/encryption.key", []byte(encoded), 0644)
			if err != nil {
				logger.Sugar().Errorf("Cannot create keyfile: %s", err.Error())
			}

		},
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
