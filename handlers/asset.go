package handlers

import (
	"context"
	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/helpers"
	"github.com/ttmbank/backend/storage"
	"net/http"
)

// GetAssets
func GetAssets(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAssets")
	log.Debugf("GetAssets:: w: %v", w)

	assets, err := db.GetAllAssets(ctx, helpers.GetPagination(ctx, req, nil))
	if err != nil {
		log.Errorf(err)
	}
	handlers.ReturnResult(ctx, w, assets)
}

// GetAssetsCount
func GetAssetsCount(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAssetsCount")
	log.Debugf("GetAssetsCount:: w: %v", w)

	assets, err := db.GetAllAssets(ctx, helpers.GetPagination(ctx, req, nil))
	if err != nil {
		log.Errorf(err)
	}
	handlers.ReturnResult(ctx, w, len(assets))
}
