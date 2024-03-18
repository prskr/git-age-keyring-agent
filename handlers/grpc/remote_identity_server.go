package grpc

import (
	"context"

	"buf.build/gen/go/git-age/agent/connectrpc/go/agent/v1/agentv1connect"
	agentv1 "buf.build/gen/go/git-age/agent/protocolbuffers/go/agent/v1"
	"connectrpc.com/connect"
	"github.com/99designs/keyring"
)

var _ agentv1connect.RemoteIdentityServiceHandler = (*RemoteIdentityServer)(nil)

func NewRemoteIdentityServer(kr keyring.Keyring) *RemoteIdentityServer {
	return &RemoteIdentityServer{
		keyRing: kr,
	}
}

type RemoteIdentityServer struct {
	keyRing keyring.Keyring
}

func (r RemoteIdentityServer) Unwrap(ctx context.Context, c *connect.Request[agentv1.UnwrapRequest]) (*connect.Response[agentv1.UnwrapResponse], error) {
	// TODO implement me
	panic("implement me")
}
