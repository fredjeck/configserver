package main

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"path"

	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/fredjeck/configserver/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	command := &cobra.Command{
		Use:   "keygen",
		Short: "Configserver Keygen generates a random key for encrypting sensitive values",
		Long:  `Configserver Keygen generates a random key for encrypting sensitive values`,
		Run:   keygen,
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func keygen(cmd *cobra.Command, args []string) {

	logger := logging.NewLogger()
	config, err := config.Read()

	if err != nil {
		logger.Sugar().Fatal(err)
	}

	key := encrypt.NewEncryptionKey()
	encoded := b64.StdEncoding.EncodeToString(key[:])

	logger.Sugar().Infof("Generated key is (base64) %s", encoded)

	err = os.WriteFile(path.Join(config.Home, "encryption.key"), []byte(encoded), 0644)
	if err != nil {
		logger.Sugar().Errorf("Cannot create keyfile: %s", err.Error())
	}
}
