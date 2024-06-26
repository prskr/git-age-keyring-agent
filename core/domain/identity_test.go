package domain_test

import (
	"net/url"
	"testing"

	"github.com/prskr/git-age-keyring-agent/core/domain"
)

func TestIdentity_MatchesRemotes(t *testing.T) {
	t.Parallel()
	remote1 := &url.URL{Host: "localhost", Path: "/path1"}
	remote2 := &url.URL{Host: "localhost", Path: "/path2"}
	remote3 := &url.URL{Host: "localhost1"}
	tests := []struct {
		name         string
		identity     domain.Identity
		inputRemotes []*url.URL
		want         bool
	}{
		{
			name:         "nil remote in Identity",
			identity:     domain.Identity{Remote: ""},
			inputRemotes: []*url.URL{remote2},
			want:         true,
		},
		{
			name:         "no matching remote in input",
			identity:     domain.Identity{Remote: "https://localhost/path1"},
			inputRemotes: []*url.URL{remote3},
			want:         false,
		},
		{
			name:         "matching remote in input",
			identity:     domain.Identity{Remote: "https://localhost/path1"},
			inputRemotes: []*url.URL{remote1},
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.identity.MatchesRemotes(tt.inputRemotes...); got != tt.want {
				t.Errorf("Identity.MatchesRemotes() = %v, want %v", got, tt.want)
			}
		})
	}
}
