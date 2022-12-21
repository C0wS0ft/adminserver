package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/helpers"
	"github.com/ttmbank/backend/storage"
	"net/http"
	"strconv"
)

// ServeAPI serves the API for this record store
func GetAddresses(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddresses")
	log.Debugf("GetAddresses:: w: %v", w)
	var err error

	vars := mux.Vars(req)
	seedHashString, seedOk := vars["seed_hash"]

	assetIdString, assetIdOk := vars["asset_id"]

	assetId := 0

	if assetIdOk {
		assetId, err = strconv.Atoi(assetIdString)
		if err != nil {
			handlers.ERROR_ASSET_NOT_SET(w)
			return
		}
	}

	if assetId == 0 && (!seedOk || seedHashString == "all") {
		addresses := db.GetAllAddresses(ctx)
		handlers.ReturnResult(ctx, w, addresses)
		return
	}

	if assetId != 0 && (!seedOk || seedHashString == "all") {
		asset := db.GetAsset(ctx, uint(assetId))
		addresses := db.GetAllAssetsAddresses(ctx, asset)
		handlers.ReturnResult(ctx, w, addresses)
		return
	}

	if !seedOk {
		handlers.ERROR_AUTH_SEED_NOT_SET(w, seedHashString)
		return
	}

	dbSeed := db.FindSeedHash(ctx, seedHashString)
	if !seedOk {
		handlers.ERROR_AUTH_SEED_NOT_FOUND(w, seedHashString)
		return
	}

	addresses := db.GetAddressesBySeedId(ctx, dbSeed.ID, 0, helpers.GetPagination(ctx, req, nil))
	handlers.ReturnResult(ctx, w, addresses)
}

// GetAddressesCount
func GetAddressesCount(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddressesCount")
	log.Debugf("GetAddressesCount:: w: %v", w)

	addressesCount := db.GetAddressesCount(ctx)
	handlers.ReturnResult(ctx, w, addressesCount)
}

// GetAddressesByAddressString
func GetAddressesByAddressString(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddressesByAddressString")
	log.Debugf("GetAddressesByAddressString:: w: %v", w)

	vars := mux.Vars(req)
	addressString, addressOk := vars["address"]
	if !addressOk {
		handlers.ERROR_ADDRESS_NOT_SET(w)
	}
	addresses := db.FindAddresses(ctx, addressString, "all", 0)
	handlers.ReturnResult(ctx, w, addresses)
}
