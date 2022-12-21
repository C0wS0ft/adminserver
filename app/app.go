package app

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"net/http"

	h "github.com/ttmbank/backend/adminserver/handlers"
	"github.com/ttmbank/backend/adminserver/middlewares"
	"github.com/ttmbank/backend/api/scannerapi"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/storage"
)

type App struct {
	Router     *mux.Router
	DB         *storage.Storage
	ScannerAPI *scannerapi.ScannerAPI
}

// Initialize
func (a *App) Initialize(ctx context.Context) {
	log := logger.FromContext(ctx).WithField("m", "Initialize")
	log.Debugf("Initialize:: ")

	db := new(storage.Storage)
	err := db.InitPostgress(ctx, viper.GetString("db.host"), viper.GetInt("db.port"), viper.GetString("db.name"), viper.GetString("db.user"), viper.GetString("db.password"), nil)
	if err != nil {
		log.Errorf("Cannot connect to db connection: %v", err)
		return
	}

	a.ScannerAPI = scannerapi.NewScannerAPI(ctx)
	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters(ctx)
}

// Run
func (a *App) Run(ctx context.Context, host string) {
	log := logger.FromContext(ctx).WithField("m", "Run")
	log.Debugf("Run:: host: %v", host)

	a.Router.Use(middlewares.LoggingMiddleware)
	a.Router.Use(middlewares.AuthMiddlewareGenerator(ctx, a.DB.DB))

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)(a.Router)

	log.Errorf(http.ListenAndServe(host, cors))
}

// Get
func (a *App) Get(ctx context.Context, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log := logger.FromContext(ctx).WithField("m", "Get")
	log.Debugf("Get:: path: %v, f: %v", path, f)

	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(ctx context.Context, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log := logger.FromContext(ctx).WithField("m", "Post")
	log.Debugf("Post:: path: %v, f: %v", path, f)

	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(ctx context.Context, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log := logger.FromContext(ctx).WithField("m", "Put")
	log.Debugf("Put:: path: %v, f: %v", path, f)

	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(ctx context.Context, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log := logger.FromContext(ctx).WithField("m", "Delete")
	log.Debugf("Delete:: path: %v, f: %v", path, f)

	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// setRouters
func (a *App) setRouters(ctx context.Context) {
	log := logger.FromContext(ctx).WithField("m", "setRouters")
	log.Debugf("setRouters:: ")

	a.Post(ctx, "/api/admin/login", a.UserLogin)
	a.Get(ctx, "/api/admin/seeds", a.GetSeeds)
	a.Get(ctx, "/api/admin/seeds/count", a.GetSeedsCount)
	a.Get(ctx, "/api/admin/assets", a.GetAssets)
	a.Get(ctx, "/api/admin/assets/count", a.GetAssetsCount)
	a.Get(ctx, "/api/admin/addresses/count", a.GetAddressesCount)
	a.Get(ctx, "/api/admin/addresses/{asset_id}/{seed_hash}", a.GetAddresses)
	a.Get(ctx, "/api/admin/address/{address}", a.GetAddressesByAddressString)
	a.Get(ctx, "/api/admin/tx/{tx_hash}", a.GetTx)
	a.Get(ctx, "/api/admin/txs/count", a.GetTxsCount)
	a.Get(ctx, "/api/admin/txs/{address}/{asset_id}", a.GetTxs)
	a.Get(ctx, "/api/admin/rescan/{currency}/{block_from}/{block_to}", a.RescanBlocks)

}

// UserLogin
func (a *App) UserLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "UserLogin")
	log.Debugf("UserLogin:: w: %v", w)

	h.UserLogin(a.DB, w, r)
}

// GetSeeds
func (a *App) GetSeeds(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetSeeds")
	log.Debugf("GetSeeds:: w: %v", w)

	h.GetSeeds(a.DB, w, r)
}

// GetSeedsCount
func (a *App) GetSeedsCount(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetSeedsCount")
	log.Debugf("GetSeedsCount:: w: %v", w)

	h.GetSeedsCount(a.DB, w, r)
}

// GetAssets
func (a *App) GetAssets(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAssets")
	log.Debugf("GetAssets:: w: %v", w)

	h.GetAssets(a.DB, w, r)
}

// GetAssetsCount
func (a *App) GetAssetsCount(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAssetsCount")
	log.Debugf("GetAssetsCount:: w: %v", w)

	h.GetAssetsCount(a.DB, w, r)
}

// GetAddresses
func (a *App) GetAddresses(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddresses")
	log.Debugf("GetAddresses:: w: %v", w)

	h.GetAddresses(a.DB, w, r)
}

// GetAddressesByAddressString
func (a *App) GetAddressesByAddressString(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddressesByAddressString")
	log.Debugf("GetAddressesByAddressString:: w: %v", w)

	h.GetAddressesByAddressString(a.DB, w, r)
}

// GetAddressesCount
func (a *App) GetAddressesCount(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetAddressesCount")
	log.Debugf("GetAddressesCount:: w: %v", w)

	h.GetAddressesCount(a.DB, w, r)
}

// GetTxs
func (a *App) GetTxs(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTxs")
	log.Debugf("GetTxs:: w: %v", w)

	h.GetTransactions(a.DB, w, r)
}

// GetTx
func (a *App) GetTx(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTx")
	log.Debugf("GetTx:: w: %v", w)

	h.GetTransaction(a.DB, w, r)
}

// GetTxsCount
func (a *App) GetTxsCount(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "GetTxsCount")
	log.Debugf("GetTxsCount:: w: %v", w)

	h.GetTransactionsCount(a.DB, w, r)
}

// RescanBlocks
func (a *App) RescanBlocks(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "RescanBlocks")
	log.Debugf("RescanBlocks:: w: %v", w)

	h.RescanBlocks(a.DB, a.ScannerAPI, w, r)
}
