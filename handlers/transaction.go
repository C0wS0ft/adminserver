package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/helpers"
	"github.com/ttmbank/backend/models"
	"github.com/ttmbank/backend/storage"
)

// GetTransactionsCount
func GetTransactionsCount(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTransactionsCount")
	log.Debugf("GetTransactionsCount:: w: %v", w)

	txsCount := db.GetTxsCount(ctx)
	handlers.ReturnResult(ctx, w, txsCount)
}

// GetTransactions
func GetTransactions(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTransactions")
	log.Debugf("GetTransactions:: w: %v", w)

	vars := mux.Vars(req)
	addressString, ok := vars["address"]
	if !ok {
		handlers.ERROR_ADDRESS_NOT_SET(w)
		return
	}
	assetIdString, ok := vars["asset_id"]
	if !ok {
		handlers.ERROR_ASSET_NOT_SET(w)
		return
	}
	assetId, err := strconv.Atoi(assetIdString)
	if err != nil {
		handlers.ERROR_ASSET_INVALID(w, assetIdString)
		return
	}

	if addressString == "all" && assetId == 0 {
		txs := db.FindAllTxs(ctx, helpers.GetPagination(ctx, req, nil))
		handlers.ReturnResult(ctx, w, txs)
		return
	}

	if addressString != "all" && assetId == 0 {
		address := db.FindAddress(ctx, addressString, "", models.AssetTypeBasic)
		if address == nil {
			return
		}
		txs := db.FindTxsByBasicAddress(ctx, address, helpers.GetPagination(ctx, req, nil))
		handlers.ReturnResult(ctx, w, txs)
		return
	}

	asset := db.GetAsset(ctx, uint(assetId))
	if addressString == "all" && assetId != 0 && asset != nil {
		txs := db.FindAssetAllTxs(ctx, asset, helpers.GetPagination(ctx, req, nil))
		handlers.ReturnResult(ctx, w, txs)
		return
	}

	address := db.FindAddressWithAsset(ctx, addressString, asset)

	if address == nil || address.ID == 0 {
		handlers.ERROR_ADDRESS_NOT_FOUND(w, addressString)
		return
	}

	txs := db.FindTxsByAddressAndAsset(ctx, address, asset)
	handlers.ReturnResult(ctx, w, txs)
}

// GetTransaction
func GetTransaction(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTransaction")
	log.Debugf("GetTransaction:: w: %v", w)

	vars := mux.Vars(req)
	txHash, ok := vars["tx_hash"]
	if !ok {
		handlers.ERROR_TRX_NOT_SET(w)
		return
	}
	tx := db.GetTx(ctx, txHash)
	if tx == nil {
		handlers.ERROR_TRX_NOT_FOUND(w)
		return
	}
	handlers.ReturnResult(ctx, w, tx)
}
