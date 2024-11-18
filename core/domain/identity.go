package domain

import (
	"cmp"
	"net/url"
	"slices"
	"strings"

	giturls "github.com/whilp/git-urls"
)

type Identity struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	Remote     string `json:"remote,omitempty"`
}

func (i Identity) MatchesRemotes(remotes ...*url.URL) bool {
	if i.Remote == "" {
		return true
	}

	parsed, err := giturls.Parse(i.Remote)
	if err != nil {
		return true
	}

	_, found := slices.BinarySearchFunc(remotes, parsed, func(u1 *url.URL, u2 *url.URL) int {
		if u1.Host == "" && u1.Path == "" {
			return 0
		}
		if i := cmp.Compare(u1.Host, u2.Host); i != 0 {
			return i
		}

		return cmp.Compare(strings.TrimLeft(u1.Path, "/"), strings.TrimLeft(u2.Path, "/"))
	})

	return found
}
