package coinbasepro

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
)

type Profile struct {
	Active    bool   `json:"active"`
	CreatedAt Time   `json:"created_at"`
	ID        string `json:"id"`
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
}

type ProfileFilter struct {
	Active bool `json:"active"`
}

func (p ProfileFilter) Params() []string {
	var params []string
	if p.Active {
		params = append(params, "active")
	}
	return params
}

type ProfileTransferSpec struct {
	Amount   float64      `json:"amount,string"`
	Currency CurrencyName `json:"currency"`
	From     string       `json:"from"`
	To       string       `json:"to"`
}

type ProfileTransfer struct {
	Amount   float64      `json:"amount"`
	Currency CurrencyName `json:"currency"`
	From     string       `json:"from"`
	To       string       `json:"to"`
}

// ListProfiles retrieves a list of Profiles (portfolio equivalents). A given user can have a maximum of 10 profiles.
// The list is not paginated.
func (c *CoinbasePro) ListProfiles(ctx context.Context, filter ProfileFilter) ([]Profile, error) {
	var profiles []Profile
	path := fmt.Sprintf("/%s/%s/", coinbaseproProfiles, qsx.Query(filter.Params()))
	return profiles, c.API.Get(ctx, path, &profiles)
}

// GetProfile retrieves the details of a single Profile.
func (c *CoinbasePro) GetProfile(ctx context.Context, profileID string) (Profile, error) {
	var profile Profile
	path := fmt.Sprintf("/%s/%s/", coinbaseproProfiles, profileID)
	return profile, c.API.Get(ctx, path, &profile)
}

// CreateProfileTransfer transfers funds between user Profiles.
func (c *CoinbasePro) CreateProfileTransfer(ctx context.Context, transferSpec ProfileTransferSpec) (ProfileTransfer, error) {
	var transfer ProfileTransfer
	path := fmt.Sprintf("/%s", coinbaseproProfilesTransfer)
	return transfer, c.API.Post(ctx, path, transferSpec, &transfer)
}
