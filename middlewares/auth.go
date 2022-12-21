package middlewares

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"

	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/models"
)

// AuthMiddlewareGenerator
func AuthMiddlewareGenerator(ctx context.Context, db *gorm.DB) (mw func(http.Handler) http.Handler) {
	log := logger.FromContext(ctx).WithField("m", "AuthMiddlewareGenerator")
	log.Debugf("AuthMiddlewareGenerator:: ")

	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Infof(r.RequestURI)

			notAuth := []string{"/api/admin/login", "/api/admin/register"}
			requestPath := r.URL.Path

			// check if path doesn't require authorization
			for _, value := range notAuth {

				if value == requestPath {
					next.ServeHTTP(w, r)
					return
				}
			}

			tokenHeader := r.Header.Get("Authorization") // get token from HTTP header

			if tokenHeader == "" {
				handlers.ERROR_AUTH_MISSING(w)
				return
			}

			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				handlers.ERROR_AUTH_INVALID(w, tokenHeader)
				return
			}

			// obtain JWT token
			tokenPart := splitted[1]
			tk := models.AuthToken{}

			// parse JWT token
			token, err := jwt.ParseWithClaims(tokenPart, &tk, func(token *jwt.Token) (interface{}, error) {
				return []byte("cryptosecret"), nil
			})

			// cannot parse JWT token
			if err != nil {
				handlers.ERROR_AUTH_CANNOT_PARSE_TOKEN(w)
				return
			}

			// token is not valid
			if !token.Valid {
				handlers.ERROR_AUTH_TOKEN_INVALID(w)
				return
			}

			dbUser := new(models.User)
			db.First(&dbUser, "username = ?", tk.Username)
			if dbUser.ID == 0 {
				handlers.ERROR_AUTH_USER_NOT_FOUND(w, tk.Username)
				return
			}

			if dbUser.Type != models.UserTypeAdmin {
				handlers.ERROR_AUTH_NO_PERMISSION(w, fmt.Sprint(dbUser.Type))
				return
			}

			ctx := context.WithValue(r.Context(), "user", dbUser)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
	return
}
