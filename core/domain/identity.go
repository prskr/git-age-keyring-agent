package domain

import (
	"cmp"
	"net/url"
	"slices"
)

type Identity struct {
	PublicKey  string   `json:"publicKey"`
	PrivateKey string   `json:"privateKey"`
	Remote     *url.URL `json:"remote,omitempty"`
}

func (i Identity) MatchesRemotes(remotes ...*url.URL) bool {
	if i.Remote == nil {
		return true
	}

	_, found := slices.BinarySearchFunc(remotes, i.Remote, func(u1 *url.URL, u2 *url.URL) int {
		if i := cmp.Compare(u1.Host, u2.Host); i != 0 {
			return i
		}

		if u1.Path == "" || u2.Path == "" {
			return 0
		}

		return cmp.Compare(u1.Path, u2.Path)
	})

	return found
}
