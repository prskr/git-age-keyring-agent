package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"slices"

	"buf.build/gen/go/git-age/agent/connectrpc/go/agent/v1/agentv1connect"
	agentv1 "buf.build/gen/go/git-age/agent/protocolbuffers/go/agent/v1"
	"connectrpc.com/connect"
	"github.com/99designs/keyring"
	giturls "github.com/whilp/git-urls"

	"github.com/prskr/git-age-keyring-agent/core/domain"
)

var (
	_                 agentv1connect.IdentitiesStoreServiceHandler = (*KeyVaultServer)(nil)
	ErrConflictingKey                                              = errors.New("conflicting key")
)

func NewAgentServer(kr keyring.Keyring) *KeyVaultServer {
	return &KeyVaultServer{
		KeyRing: kr,
	}
}

type KeyVaultServer struct {
	KeyRing keyring.Keyring
}

func (a *KeyVaultServer) GetIdentities(
	_ context.Context,
	req *connect.Request[agentv1.GetIdentitiesRequest],
) (*connect.Response[agentv1.GetIdentitiesResponse], error) {
	keys, err := a.KeyRing.Keys()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slices.Sort(req.Msg.Remotes)

	urls, err := a.parseRemotesToURLs(req.Msg.Remotes)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	keysResponse := new(agentv1.GetIdentitiesResponse)

	for _, key := range keys {
		item, err := a.KeyRing.Get(key)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		var id domain.Identity
		if err := json.Unmarshal(item.Data, &id); err != nil {
			slog.Error("failed to unmarshal identity", slog.String("err", err.Error()))
			continue
		}

		if id.MatchesRemotes(urls...) {
			keysResponse.Keys = append(keysResponse.Keys, id.PrivateKey)
		}
	}

	return connect.NewResponse(keysResponse), nil
}

func (a *KeyVaultServer) StoreIdentity(
	_ context.Context,
	req *connect.Request[agentv1.StoreIdentityRequest],
) (*connect.Response[agentv1.StoreIdentityResponse], error) {
	if !req.Msg.Overwrite {
		if _, err := a.KeyRing.Get(req.Msg.PublicKey); err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, ErrConflictingKey)
		}
	}

	if req.Msg.Remote != "" {
		_, err := giturls.Parse(req.Msg.Remote)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	id := domain.Identity{
		PublicKey:  req.Msg.PublicKey,
		PrivateKey: req.Msg.PrivateKey,
		Remote:     req.Msg.Remote,
	}

	itemData, err := json.Marshal(id)
	if err != nil {
		slog.Error("Failed to marshal identity", slog.String("err", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	item := keyring.Item{
		Key:         req.Msg.PublicKey,
		Data:        itemData,
		Description: req.Msg.Comment,
	}

	if err := a.KeyRing.Set(item); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(new(agentv1.StoreIdentityResponse)), nil
}

// parseRemotesToURLs parses raw remotes to URL objects
func (a *KeyVaultServer) parseRemotesToURLs(remotes []string) (urls []*url.URL, err error) {
	urls = make([]*url.URL, 0, len(remotes))
	for _, raw := range remotes {
		var parsed *url.URL
		if parsed, err = giturls.Parse(raw); err != nil {
			break
		}
		urls = append(urls, parsed)
	}
	return urls, err
}
