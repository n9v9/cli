package sshkey

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create FLAGS",
			Short: "Create a SSH key",
			Args:  cobra.NoArgs,
		}
		cmd.Flags().String("name", "", "Key name (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().String("public-key", "", "Public key")
		cmd.Flags().String("public-key-from-file", "", "Path to file containing public key")
		cmd.MarkFlagsMutuallyExclusive("public-key", "public-key-from-file")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		publicKey, _ := cmd.Flags().GetString("public-key")
		publicKeyFile, _ := cmd.Flags().GetString("public-key-from-file")
		labels, _ := cmd.Flags().GetStringToString("label")

		if publicKeyFile != "" {
			var (
				data []byte
				err  error
			)
			if publicKeyFile == "-" {
				data, err = io.ReadAll(os.Stdin)
			} else {
				data, err = os.ReadFile(publicKeyFile)
			}
			if err != nil {
				return nil, nil, err
			}
			publicKey = string(data)
		}

		opts := hcloud.SSHKeyCreateOpts{
			Name:      name,
			PublicKey: publicKey,
			Labels:    labels,
		}
		sshKey, _, err := s.Client().SSHKey().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		cmd.Printf("SSH key %d created\n", sshKey.ID)

		return sshKey, util.Wrap("ssh_key", hcloud.SchemaFromSSHKey(sshKey)), nil
	},
}
