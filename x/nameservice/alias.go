package nameservice

import (
	"github.com/caosbad/nameservice/x/nameservice/keeper"
	"github.com/caosbad/nameservice/x/nameservice/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	StoreKey        = types.StoreKey
	AuctionStoreKey = types.AuctionStoreKey
	QuerierRoute    = types.QuerierRoute
)

var (
	NewKeeper        = keeper.NewKeeper
	//NewAuctionKeeper = keeper.NewAuctionKeeper
	NewQuerier       = keeper.NewQuerier
	NewMsgBuyName    = types.NewMsgBuyName
	NewMsgSetName    = types.NewMsgSetName
	NewMsgDeleteName = types.NewMsgDeleteName
	NewWhois         = types.NewWhois
	NewAuction       = types.NewAuction
	ModuleCdc        = types.ModuleCdc
	RegisterCodec    = types.RegisterCodec
)

type (
	Keeper          = keeper.Keeper
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	MsgDeleteName   = types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois
)
