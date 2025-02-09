package network

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var RemoveRouteCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-route NETWORK FLAGS",
			Short:                 "Remove a route from a network",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().IPNet("destination", net.IPNet{}, "Destination network or host (required)")
		cmd.MarkFlagRequired("destination")

		cmd.Flags().IP("gateway", net.IP{}, "Gateway IP address (required)")
		cmd.MarkFlagRequired("gateway")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		gateway, _ := cmd.Flags().GetIP("gateway")
		destination, _ := cmd.Flags().GetIPNet("destination")
		idOrName := args[0]
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", idOrName)
		}

		opts := hcloud.NetworkDeleteRouteOpts{
			Route: hcloud.NetworkRoute{
				Gateway:     gateway,
				Destination: &destination,
			},
		}
		action, _, err := s.Client().Network().DeleteRoute(s, network, opts)
		if err != nil {
			return err
		}
		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}
		cmd.Printf("Route removed from network %d\n", network.ID)

		return nil
	},
}
