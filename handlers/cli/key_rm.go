package cli

import (
	"log/slog"

	"github.com/99designs/keyring"
	"github.com/alecthomas/kong"
)

type RemoveKeyCliHandler struct {
	ServiceNameFlag `embed:""`
	PublicKey       string `arg:"" help:"public key to remove from the keyring"`
}

func (h RemoveKeyCliHandler) Run(logger *slog.Logger, kr keyring.Keyring) error {
	logger.Info("Trying to remove key", slog.String("public_key", h.PublicKey))
	return kr.Remove(h.PublicKey)
}

func (h RemoveKeyCliHandler) AfterApply(kongCtx *kong.Context) error {
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
