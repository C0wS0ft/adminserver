package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/ttmbank/backend/common/handlers"
	"github.com/ttmbank/backend/common/logger"
	"github.com/ttmbank/backend/models"
	"github.com/ttmbank/backend/storage"
	"net/http"
)

// UserLogin
func UserLogin(ctx context.Context, db *storage.Storage, w http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(ctx).WithField("m", "UserLogin")
	log.Debugf("UserLogin:: w: %v", w)

	var user models.User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		handlers.ERROR_BAD_REQUEST(w, err.Error())
		return
	}
	dbUser := db.GetUser(ctx, user.Username)
	if dbUser == nil {
		handlers.ERROR_AUTH_USER_NOT_FOUND(w, user.Username)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	if hex.EncodeToString(hash.Sum(nil)) != dbUser.Password {
		handlers.ERROR_AUTH_BAD_PASSWORD(w, "")
		return
	}

	authToken := &models.AuthToken{Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)

	tokenString, error := token.SignedString([]byte("cryptosecret"))
	if error != nil {
		handlers.ERROR_AUTH_FORBIDDEN(w, error.Error())
		return
	}
	handlers.ReturnResult(ctx, w, tokenString)
}
