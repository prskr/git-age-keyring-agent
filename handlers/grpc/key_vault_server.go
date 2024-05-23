package grpc

import (
	"context"
	"encoding/json"
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

var _ agentv1connect.IdentitiesStoreServiceHandler = (*KeyVaultServer)(nil)

func NewAgentServer(kr keyring.Keyring) *KeyVaultServer {
	return &KeyVaultServer{
		keyRing: kr,
	}
}

type KeyVaultServer struct {
	keyRing keyring.Keyring
}

func (a *KeyVaultServer) GetIdentities(
	_ context.Context,
	req *connect.Request[agentv1.GetIdentitiesRequest],
) (*connect.Response[agentv1.GetIdentitiesResponse], error) {
	keys, err := a.keyRing.Keys()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slices.Sort(req.Msg.Remotes)

	urls := make([]*url.URL, 0, len(req.Msg.Remotes))
	for _, raw := range req.Msg.Remotes {
		if parsed, err := giturls.Parse(raw); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		} else {
			urls = append(urls, parsed)
		}
	}

	keysResponse := new(agentv1.GetIdentitiesResponse)

	for _, key := range keys {
		item, err := a.keyRing.Get(key)
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
	var (
		parsedUrl *url.URL
		err       error
	)

	if req.Msg.Remote != "" {
		parsedUrl, err = giturls.Parse(req.Msg.Remote)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	id := domain.Identity{
		PublicKey:  req.Msg.PublicKey,
		PrivateKey: req.Msg.PrivateKey,
		Remote:     parsedUrl,
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

	if err := a.keyRing.Set(item); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(new(agentv1.StoreIdentityResponse)), nil
}