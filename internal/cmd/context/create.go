package context

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func newCreateCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS] NAME",
		Short:                 "Create a new context",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runCreate),
	}
	return cmd
}

func runCreate(s state.State, cmd *cobra.Command, args []string) error {
	cfg := s.Config()
	if !state.StdoutIsTerminal() {
		return errors.New("context create is an interactive command")
	}

	name := strings.TrimSpace(args[0])
	if name == "" {
		return errors.New("invalid name")
	}
	if config.ContextByName(cfg, name) != nil {
		return errors.New("name already used")
	}

	context := &config.Context{Name: name}

	var token string

	envToken := os.Getenv("HCLOUD_TOKEN")
	if envToken != "" {
		if len(envToken) != 64 {
			cmd.Println("Warning: HCLOUD_TOKEN is set, but token is invalid (must be exactly 64 characters long)")
		} else {
			cmd.Print("The HCLOUD_TOKEN environment variable is set. Do you want to use the token from HCLOUD_TOKEN for the new context? (Y/n): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if s := strings.ToLower(scanner.Text()); s == "" || s == "y" || s == "yes" {
				token = envToken
			}
		}
	}

	if token == "" {
		for {
			cmd.Printf("Token: ")
			// Conversion needed for compilation on Windows
			//                               vvv
			btoken, err := term.ReadPassword(int(syscall.Stdin))
			cmd.Print("\n")
			if err != nil {
				return err
			}
			token = string(bytes.TrimSpace(btoken))
			if token == "" {
				continue
			}
			if len(token) != 64 {
				cmd.Print("Entered token is invalid (must be exactly 64 characters long)\n")
				continue
			}
			break
		}
	}

	context.Token = token

	cfg.SetContexts(append(cfg.Contexts(), context))
	cfg.SetActiveContext(context)

	if err := cfg.Write(); err != nil {
		return err
	}

	cmd.Printf("Context %s created and activated\n", name)

	return nil
}
