package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/ttmbank/backend/api/scannerapi"
	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/models"
	"github.com/ttmbank/backend/storage"
)

// RescanBlocks
func RescanBlocks(ctx context.Context, db *storage.Storage, scannerApi *scannerapi.ScannerAPI, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "RescanBlocks")
	log.Debugf("RescanBlocks:: w: %v", w)

	vars := mux.Vars(req)
	currencyString, ok := vars["currency"]
	if !ok {
		handlers.ERROR_ADDRESS_NOT_SET(w)
		return
	}
	blockFromString, ok := vars["block_from"]
	if !ok {
		handlers.ERROR_ASSET_NOT_SET(w)
		return
	}
	blockToString, ok := vars["block_to"]
	if !ok {
		handlers.ERROR_ASSET_NOT_SET(w)
		return
	}
	if blockFromString == "all" {
		basicAsset := db.FindAssetBySymbol(ctx, currencyString, models.AssetTypeBasic)
		if basicAsset == nil {
			handlers.ERROR_ASSET_NOT_FOUND(w, currencyString)
			return
		}
		txes := db.FindAssetAllTxs(ctx, basicAsset, nil)
		blockMap := make(map[uint64]bool)
		for _, t := range txes {
			blockMap[t.Block] = true
		}
		for blockNumber, _ := range blockMap {
			result, err := scannerApi.RescanBlocks(ctx, currencyString, uint(blockNumber), uint(blockNumber))
			if err != nil || !result {
				log.Infof("Error rescanning block: %v, %v, %v", blockNumber, result, err)
			}
		}

	} else {
		blockFrom, err := strconv.Atoi(blockFromString)
		if err != nil {
			handlers.ERROR_BAD_REQUEST(w, "Cannot parse block from")
			return
		}
		blockTo, err := strconv.Atoi(blockToString)
		if err != nil {
			handlers.ERROR_BAD_REQUEST(w, "Cannot parse block to")
			return
		}

		result, err := scannerApi.RescanBlocks(ctx, currencyString, uint(blockFrom), uint(blockTo))
		if err != nil {
			handlers.ERROR_BAD_REQUEST(w, err.Error())
			return
		}

		handlers.ReturnResult(ctx, w, result)
	}
}
