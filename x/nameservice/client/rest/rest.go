package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	//"github.com/caosbad/nameservice/x/nameservice/types"

	//"github.com/cosmos/cosmos-sdk/codec"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	restName = "name"
)




// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), namesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoIsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), deleteNameHandler(cliCtx)).Methods("DELETE")
	r.HandleFunc(fmt.Sprintf("/%s/names/auction", storeName), setAuctionNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/bid", storeName), bidNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names/claim", storeName), claimNameHandler(cliCtx)).Methods("POST")
}
