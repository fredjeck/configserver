package commands

import (
	b64 "encoding/base64"
	"path"

	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/spf13/cobra"
)

var KeygenCommand = &cobra.Command{
	Use:   "keygen",
	Short: "Generates a new encryption key",
	Long: `Generates a new encryption key and stores it under the CONFIGSERVER_HOME directory, replacing any existing key.
Warning: Generating a new key if a key already exists will render existing configurations useless.
	`,
	Run: keygen,
}

func init() {
	KeygenCommand.Flags().BoolP("dry-run", "d", false, "If true does not proceed to the encryption key replacement but rather prints it to the console")
	RootCommand.AddCommand(KeygenCommand)
}

func keygen(cmd *cobra.Command, _ []string) {
	Logger.Sugar().Info(`If an encryption key exists and is in use, running keygen will overwrite any existing key - rendering currently served configuration useless.
You will need to rehash any encrypted sensitive values in your configuration files`)

	key := encrypt.NewEncryptionKey()
	encoded := b64.StdEncoding.EncodeToString(key[:])

	Logger.Sugar().Infof("Generated key is (base64) %s", encoded)
	target := path.Join(Configuration.Home, "encryption.key")

	dr, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		Logger.Sugar().Panicf("Unable to parse the command line: %s", err)
	}

	if dr {
		Logger.Sugar().Infof("If dry-run would have been set to false, the encryption key would have been written to  '%s'", target)
		return
	}

	err = encrypt.StoreEncryptionKey(key, Configuration.EncryptionKeyPath())
	if err != nil {
		Logger.Sugar().Errorf("Cannot create keyfile: %s", err.Error())
	}
	Logger.Sugar().Infof("Keyfile generated to '%s'", target)
}
