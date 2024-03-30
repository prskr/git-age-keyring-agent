package cli

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/99designs/keyring"
	"github.com/alecthomas/kong"

	"github.com/prskr/git-age-keyring-agent/core/domain"
	"github.com/prskr/git-age-keyring-agent/core/ports"
)

type ListKeysCliHandler struct {
	ServiceNameFlag `embed:""`
}

func (h ListKeysCliHandler) Run(stdout ports.STDOUT, kr keyring.Keyring) error {
	keys, err := kr.Keys()
	if err != nil {
		return fmt.Errorf("failed to list keys in keyring: %w", err)
	}

	tw := tabwriter.NewWriter(stdout, 0, 0, 1, ' ', 0)
	defer func() {
		_ = tw.Flush()
	}()

	if _, err := fmt.Fprintln(tw, "Public Key\tComment\tRemote"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, key := range keys {
		item, err := kr.Get(key)
		if err != nil {
			return fmt.Errorf("failed to get key from keyring: %w", err)
		}

		var id domain.Identity
		if err := json.Unmarshal(item.Data, &id); err != nil {
			return fmt.Errorf("failed to unmarshal identity: %w", err)
		}

		remote := "<none>"
		if id.Remote != nil {
			remote = id.Remote.String()
		}

		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", id.PublicKey, item.Description, remote); err != nil {
			return fmt.Errorf("failed to write identity: %w", err)
		}
	}

	return nil
}

func (h ListKeysCliHandler) AfterApply(kongCtx *kong.Context) error {
	keyRingCfg := keyring.Config{
		ServiceName: h.ServiceName,
	}

	kr, err := keyring.Open(keyRingCfg)
	if err != nil {
		return err
	}

	kongCtx.BindTo(kr, (*keyring.Keyring)(nil))

	return nil
}
