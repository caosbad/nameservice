package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
	IsAuction   bool           `json:"isAuction"`
	BlockHeight int64          `json:"blockHeight"`
	BidUser     sdk.AccAddress `json:"bidUser"`
}

type Auction struct {
	Whois
	Name string `json:name`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

func NewAuction(whois Whois, name string) Auction {
	return Auction{
		Whois: whois,
		Name:  name,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s Value: %s Price: %s IsAuction: %s `, w.Owner, w.Value, w.Price, w.IsAuction))
}

// implement fmt.Stringer
func (a Auction) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s Owner: %s Value: %s Price: %s IsAuction: %s BlockHeight: %s`,a.Name, a.Owner, a.Value, a.Price, a.IsAuction, a.BlockHeight))
}