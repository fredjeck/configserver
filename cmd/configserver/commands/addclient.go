package commands

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/spf13/cobra"
)

var AddClientCommand = &cobra.Command{
	Use:   "addclient",
	Short: "Adds a new client to a repository",
	Long: `Adds a new client to a repository.
	`,
	Run: addClient,
}

func init() {
	AddClientCommand.Flags().StringP("repository", "r", "", "target repository as configured in the configserver.yaml file")
	AddClientCommand.Flags().StringP("clientid", "i", "", "client id")
	RootCommand.AddCommand(AddClientCommand)
}

func addClient(cmd *cobra.Command, args []string) {
	clientid, err_client := cmd.Flags().GetString("clientid")
	if err_client != nil {
		return
	}
	repo, err_repo := cmd.Flags().GetString("repository")
	if err_repo != nil {
		return
	}

	secret, err := encrypt.Encrypt([]byte(repo+":"+clientid), Key)
	if err != nil {
		return
	}

	fmt.Print(string(b64.StdEncoding.EncodeToString(secret)))
}
