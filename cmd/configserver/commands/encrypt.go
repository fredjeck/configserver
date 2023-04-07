package commands

import (
	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/spf13/cobra"
)

var EncryptCommand = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypts a provided value",
	Long:  `Encrypts the provided value into a token which can then be embedded into a configuration file served via this configserver instance`,
	Run:   encryptValue,
}

func init() {
	EncryptCommand.Flags().StringP("value", "v", "", "value to encrypt")
	RootCommand.AddCommand(EncryptCommand)
}

func encryptValue(cmd *cobra.Command, _ []string) {
	value, err := cmd.Flags().GetString("value")
	if len(value) == 0 || err != nil {
		Logger.Sugar().Fatal("Missing mandatory argument : value")
	}

	enc, err := encrypt.EncryptToken([]byte(value), Key)
	if err != nil {
		Logger.Sugar().Fatalf("an error occured while encrypting the provided value: %s", err.Error())
	}
	Logger.Sugar().Infof("Encrypted token: %s", enc)
}
