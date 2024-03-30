package cli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"filippo.io/age"
	"github.com/99designs/keyring"
	"github.com/alecthomas/kong"
	giturls "github.com/whilp/git-urls"

	"github.com/prskr/git-age-keyring-agent/core/domain"
)

type ImportCliHandler struct {
	ServiceNameFlag `embed:""`
	Comment         string `short:"c" name:"comment" help:"Comment to add in file"`
	Remote          string `short:"r" name:"remote" help:"Remote for which this key should be considered"`
	PrivateKey      string `arg:""`
}

func (h ImportCliHandler) Run(kr keyring.Keyring) error {
	ageIdentity, err := age.ParseX25519Identity(h.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	var parsedUrl *url.URL

	if h.Remote != "" {
		parsedUrl, err = giturls.Parse(h.Remote)
		if err != nil {
			return fmt.Errorf("failed to parse remote: %w", err)
		}
	}

	id := domain.Identity{
		PublicKey:  ageIdentity.Recipient().String(),
		PrivateKey: ageIdentity.String(),
		Remote:     parsedUrl,
	}

	itemData, err := json.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to marshal identity: %w", err)
	}

	item := keyring.Item{
		Key:         ageIdentity.Recipient().String(),
		Data:        itemData,
		Description: h.Comment,
	}

	if err := kr.Set(item); err != nil {
		return err
	}

	return nil
}

func (h ImportCliHandler) AfterApply(kongCtx *kong.Context) error {
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
