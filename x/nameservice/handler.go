package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//. "github.com/caosbad/nameservice/x/nameservice/keeper"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/caosbad/nameservice/x/nameservice/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case types.MsgSetAuction:
			return handleMsgSetAuction(ctx, keeper, msg)
		case types.MsgBidName:
			return handleMsgBidName(ctx, keeper, msg)
		case types.MsgClaimName:
			return handleMsgClaimName(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}


// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return &sdk.Result{}, nil                // return
}


// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}


// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.DeleteWhois(ctx, msg.Name)
	return &sdk.Result{}, nil
}


// Handle a message to set auction
func handleMsgSetAuction(ctx sdk.Context, keeper Keeper, msg types.MsgSetAuction) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	if keeper.GetAuction(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Auction already started")
	}

	keeper.SetAuction(ctx, msg.Name, true)
	keeper.SetPrice(ctx, msg.Name, msg.Price)
	keeper.AddAuction(ctx, msg.Name)

	return &sdk.Result{}, nil
}

// Handle a message bid
func handleMsgBidName(ctx sdk.Context, keeper Keeper, msg types.MsgBidName) (*sdk.Result, error) {
	if !keeper.GetAuction(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Canot bid")
	}
	if msg.Bider.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Owner canot bid")
	}

	//if keeper.GetAuction(ctx, msg.Name) {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Auction already started")
	//}
	if keeper.GetPrice(ctx, msg.Name).IsAllGTE(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}

	bidBlockHeight := keeper.GetBidHeight(ctx, msg.Name)
	currentBlockHeight := ctx.BlockHeight()
	// have bider before
	if bidBlockHeight > 0 {
		if currentBlockHeight-bidBlockHeight >= 100 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Auction is over") //
		}
		// return funds to the bider before
		_, err := keeper.CoinKeeper.AddCoins(ctx, keeper.GetBidUser(ctx, msg.Name), keeper.GetPrice(ctx, msg.Name)) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
		_, err1 := keeper.CoinKeeper.SubtractCoins(ctx, msg.Bider, msg.Bid) // If so, deduct the Bid amount from the sender
		if err1 != nil {
			return nil, err
		}

	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Bider, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
	}

	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	keeper.SetBidUser(ctx, msg.Name, msg.Bider)
	keeper.SetBidHeight(ctx, msg.Name, currentBlockHeight)
	// TODO add auction logic
	//keeper.SetAuction(ctx, msg.Name, true)
	return &sdk.Result{}, nil
}

// Handle a message to set auction
func handleMsgClaimName(ctx sdk.Context, keeper Keeper, msg types.MsgClaimName) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	currentBlockHeight := ctx.BlockHeight()
	bidBlockHeight := keeper.GetBidHeight(ctx, msg.Name)
	// need 100 blocks
	if currentBlockHeight-bidBlockHeight < 100 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Auction is not over ")
	}
	// bider  or name owner can claim the name
	if !msg.Owner.Equals(keeper.GetBidUser(ctx, msg.Name)) && !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner ")
	}

	keeper.SetBidHeight(ctx, msg.Name, 0)
	keeper.SetOwner(ctx, msg.Name, keeper.GetBidUser(ctx, msg.Name))
	keeper.SetAuction(ctx, msg.Name, false)
	return &sdk.Result{}, nil
}
