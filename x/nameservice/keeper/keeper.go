package keeper

import (
	"github.com/caosbad/nameservice/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper types.BankKeeper

	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	auctionStoreKey sdk.StoreKey
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois types.Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// Deletes the entire Whois metadata struct for a name
func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
	if k.GetWhois(ctx, name).IsAuction {
		auctionStore := ctx.KVStore(k.auctionStoreKey)
		auctionStore.Delete([]byte(name))
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))

}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, auctionStoreKey sdk.StoreKey, coinKeeper types.BankKeeper) Keeper {
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		auctionStoreKey: auctionStoreKey,
		CoinKeeper:      coinKeeper,
	}
}

//func NewAuctionKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, coinKeeper types.BankKeeper) Keeper {
//	return Keeper{
//		cdc:        cdc,
//		storeKey:   storeKey,
//		CoinKeeper: coinKeeper,
//	}
//}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetAuction(ctx sdk.Context, name string, value bool) {
	whois := k.GetWhois(ctx, name)
	whois.IsAuction = value
	k.SetWhois(ctx, name, whois)
}

// GetAuction - get the auction state of a name
func (k Keeper) GetAuction(ctx sdk.Context, name string) bool {
	return k.GetWhois(ctx, name).IsAuction
}

// SetBlockHeight
func (k Keeper) SetBidHeight(ctx sdk.Context, name string, value int64) {
	whois := k.GetWhois(ctx, name)
	whois.BlockHeight = value
	k.SetWhois(ctx, name, whois)
}

// GetAuction - get the auction state of a name
func (k Keeper) GetBidHeight(ctx sdk.Context, name string) int64 {
	return k.GetWhois(ctx, name).BlockHeight
}

// SetBlockHeight
func (k Keeper) SetBidUser(ctx sdk.Context, name string, value sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.BidUser = value
	k.SetWhois(ctx, name, whois)
}

// GetAuction - get the auction state of a name
func (k Keeper) GetBidUser(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).BidUser
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) AddAuction(ctx sdk.Context, name string) {
	whois := k.GetWhois(ctx, name)
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.auctionStoreKey)
	auction := types.NewAuction(whois, name)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(auction))
}


// Deletes the entire Whois metadata struct for a name
func (k Keeper) DeleteAuction(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.auctionStoreKey)
	store.Delete([]byte(name))
}