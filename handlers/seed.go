package handlers

import (
	"context"
	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/helpers"
	"github.com/ttmbank/backend/storage"
	"net/http"
)

// GetSeeds
func GetSeeds(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetSeeds")
	log.Debugf("GetSeeds:: w: %v", w)

	seeds := db.GetSeeds(ctx, helpers.GetPagination(ctx, req, nil))
	handlers.ReturnResult(ctx, w, seeds)
}

// GetSeedsCount
func GetSeedsCount(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetSeedsCount")
	log.Debugf("GetSeedsCount:: w: %v", w)

	seeds := db.GetSeeds(ctx, helpers.GetPagination(ctx, req, nil))
	handlers.ReturnResult(ctx, w, len(seeds))
}
