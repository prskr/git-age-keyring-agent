package grpc_test

import (
	"encoding/json"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"

	agentv1 "buf.build/gen/go/git-age/agent/protocolbuffers/go/agent/v1"
	"connectrpc.com/connect"

	"github.com/prskr/git-age-keyring-agent/core/domain"
	"github.com/prskr/git-age-keyring-agent/handlers/grpc"
	"github.com/prskr/git-age-keyring-agent/internal/testx"
)

func TestKeyVaultServer_GetIdentities(t *testing.T) {
	const (
		publicKey  = `age1h3mfl6yg40artzltlrmsq9skxetnu8qgwr5jcyx8z0c3zxr77v2sqyy2xl`
		privateKey = `AGE-SECRET-KEY-1YCDE7N3QLSP4LPUGDDQ245GFWFLLL042UFLY556DKG36P66ZHUDSZ805HM`
	)

	t.Parallel()

	tests := []struct {
		name    string
		remotes []string
		keyring keyring.Keyring
		expectF testx.ValueAssertionFunc[[]string]
	}{
		{
			name:    "SSH remote",
			remotes: []string{"ssh://git@github.com/prskr/git-age-keyring-agent.git"},
			keyring: keyring.NewArrayKeyring([]keyring.Item{
				{
					Key: publicKey,
					Data: mustMarshal(domain.Identity{
						Remote:     "git@github.com:prskr/git-age-keyring-agent.git",
						PublicKey:  publicKey,
						PrivateKey: privateKey,
					}),
				},
			}),
			expectF: func(t assert.TestingT, val []string, vals ...[]string) bool {
				return assert.Equal(t, []string{privateKey}, val)
			},
		},
		{
			name:    "HTTPS remote",
			remotes: []string{"https://github.com/prskr/git-age-keyring-agent.git"},
			keyring: keyring.NewArrayKeyring([]keyring.Item{
				{
					Key: publicKey,
					Data: mustMarshal(domain.Identity{
						Remote:     "https://github.com/prskr/git-age-keyring-agent.git",
						PublicKey:  publicKey,
						PrivateKey: privateKey,
					}),
				},
			}),
			expectF: func(t assert.TestingT, val []string, vals ...[]string) bool {
				return assert.Equal(t, []string{privateKey}, val)
			},
		},
		{
			name:    "No remote in keyring",
			remotes: []string{"https://github.com/prskr/git-age-keyring-agent.git"},
			keyring: keyring.NewArrayKeyring([]keyring.Item{
				{
					Key: publicKey,
					Data: mustMarshal(domain.Identity{
						PublicKey:  publicKey,
						PrivateKey: privateKey,
					}),
				},
			}),
			expectF: func(t assert.TestingT, val []string, vals ...[]string) bool {
				return assert.Equal(t, []string{privateKey}, val)
			},
		},
		{
			name:    "Remote mismatch",
			remotes: []string{"https://github.com/prskr/git-age-keyring-agent.git"},
			keyring: keyring.NewArrayKeyring([]keyring.Item{
				{
					Key: publicKey,
					Data: mustMarshal(domain.Identity{
						Remote:     "https://github.com/prskr/git-age.git",
						PublicKey:  publicKey,
						PrivateKey: privateKey,
					}),
				},
			}),
			expectF: func(t assert.TestingT, val []string, vals ...[]string) bool {
				return assert.Empty(t, val)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := testx.Context(t)
			srv := grpc.NewAgentServer(tt.keyring)
			res, err := srv.GetIdentities(
				ctx,
				connect.NewRequest(&agentv1.GetIdentitiesRequest{Remotes: tt.remotes}),
			)
			assert.NoError(t, err)
			tt.expectF(t, res.Msg.Keys)
		})
	}
}

func TestKeyVaultServer_StoreIdentity(t *testing.T) {
	const (
		publicKey  = `age1gus7chkg7r8cdh4detngm46t2nsvnv20nrwlrg9ngz6a533awejq965prs`
		privateKey = `AGE-SECRET-KEY-1J6G7GNY4MX8HUN7GD0TF4TVE6WWDDK2QNQMFXH0TZH3VLAZPQUMSVPHUC8`
	)
	t.Parallel()

	tests := []struct {
		name    string
		keyring keyring.Keyring
		req     *agentv1.StoreIdentityRequest
		errF    assert.ErrorAssertionFunc
		expectF testx.ValueAssertionFunc[keyring.Keyring]
	}{
		{
			name:    "Success case - empty keyring",
			keyring: keyring.NewArrayKeyring(nil),
			req: &agentv1.StoreIdentityRequest{
				PublicKey:  publicKey,
				PrivateKey: privateKey,
			},
			errF: assert.NoError,
			expectF: func(t assert.TestingT, val keyring.Keyring, vals ...keyring.Keyring) bool {
				keys, err := val.Keys()
				return assert.NoError(t, err) && assert.Len(t, keys, 1)
			},
		},
		{
			name: "Success case - non-empty keyring - no conflict",
			keyring: keyring.NewArrayKeyring([]keyring.Item{
				{
					Key: publicKey,
					Data: mustMarshal(domain.Identity{
						Remote:     "https://github.com/prskr/git-age.git",
						PublicKey:  publicKey,
						PrivateKey: privateKey,
					}),
				},
			}),
			req: &agentv1.StoreIdentityRequest{
				Remote:     "https://github.com/prskr/git-age-keyring-agent.git",
				PublicKey:  publicKey,
				PrivateKey: privateKey,
			},
			errF: assert.NoError,
			expectF: func(t assert.TestingT, val keyring.Keyring, vals ...keyring.Keyring) bool {
				keys, err := val.Keys()
				return assert.NoError(t, err) && assert.Len(t, keys, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := testx.Context(t)
			srv := grpc.NewAgentServer(tt.keyring)
			_, err := srv.StoreIdentity(
				ctx,
				connect.NewRequest(tt.req),
			)
			tt.errF(t, err)
			if tt.expectF != nil {
				tt.expectF(t, srv.KeyRing)
			}
		})
	}
}

func mustMarshal(val any) []byte {
	data, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}

	return data
}
