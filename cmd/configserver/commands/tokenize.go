package commands

import (
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	tokenizeCmd.Flags().StringVarP(&privateKeyPath, "keys", "k", "", "Path to a directory where your id_rsa key can be found")
	_ = tokenizeCmd.MarkFlagRequired("keys")

	tokenizeCmd.Flags().StringVarP(&outFilePath, "out", "o", "", "Destination file path, if not provided the original file will be overwritten")

	tokenizeCmd.Flags().StringVarP(&inFilePath, "file", "f", "", "Destination file path, if not provided the original file will be overwritten")
	_ = tokenizeCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(tokenizeCmd)
}

var (
	tokenizeCmd = &cobra.Command{
		Use:   "tokenize",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			tokenize()
		},
	}
	privateKeyPath string
	outFilePath    string
	inFilePath     string
)

func tokenize() {
	_, err := os.Stat(inFilePath)
	dieIfError(err, "'%s' file does not exist or is not accessible", inFilePath)

	bytes, err := os.ReadFile(inFilePath)
	dieIfError(err, "'%s' cannot read file", inFilePath)

	vault, err := encryption.LoadKeyVault(privateKeyPath, false)
	dieIfError(err, "unable to initialize private key")

	tokenized, err := encryption.Tokenize(bytes, vault)
	dieIfError(err, "an error occurred while tokenizing file")

	if len(outFilePath) == 0 {
		outFilePath = inFilePath
	}

	err = os.WriteFile(outFilePath, tokenized, 0644)
	dieIfError(err, "an error occurred while writing the output file")
}
