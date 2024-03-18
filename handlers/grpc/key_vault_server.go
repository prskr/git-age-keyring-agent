package grpc

import (
	"context"

	"buf.build/gen/go/git-age/agent/connectrpc/go/agent/v1/agentv1connect"
	agentv1 "buf.build/gen/go/git-age/agent/protocolbuffers/go/agent/v1"
	"connectrpc.com/connect"
	"github.com/99designs/keyring"
)

var _ agentv1connect.KeyVaultServiceHandler = (*KeyVaultServer)(nil)

func NewAgentServer(kr keyring.Keyring) *KeyVaultServer {
	return &KeyVaultServer{
		keyRing: kr,
	}
}

type KeyVaultServer struct {
	keyRing keyring.Keyring
}

func (a *KeyVaultServer) GetKeys(context.Context, *connect.Request[agentv1.GetKeysRequest]) (*connect.Response[agentv1.GetKeysResponse], error) {
	keys, err := a.keyRing.Keys()
	if err != nil {
		return nil, err
	}

	keysResponse := new(agentv1.GetKeysResponse)

	for _, key := range keys {
		item, err := a.keyRing.Get(key)
		if err != nil {
			return nil, err
		}

		keysResponse.Keys = append(keysResponse.Keys, string(item.Data))
	}

	return connect.NewResponse(keysResponse), nil
}
