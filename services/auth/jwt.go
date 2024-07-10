package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayobami6/pickitup/config"
	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

var UserKey contextKey= "UserID"


func CreateJWT(secret []byte, userId int) (string, error) {
	exp, err := strconv.ParseInt(config.GetEnv("JWT_EXPIRATION", "25000"), 10, 64)
	if err != nil {
		return "", err
	}
	expiration := time.Second * time.Duration(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}

func Auth(handlerFunc http.HandlerFunc, userStore models.UserRepo) http.HandlerFunc {
	// return the http.Handler function
	return func(w http.ResponseWriter, r *http.Request) {
        tokenString, err := utils.GetTokenFromRequest(r)
		if err!= nil {
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }
        token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.GetEnv("JWT_SECRET", "")), nil
        })
        if err!= nil ||!token.Valid {
            forbidden(w)
            return
        }
        // get claims from jwt
		claims := token.Claims.(jwt.MapClaims)
		userID, err := strconv.Atoi(claims["userID"].(string))
		if err!= nil {
            forbidden(w)
            return
        }
		// verify user from the database
		user, err := userStore.GetUserByID(userID)
        if err!= nil {
            forbidden(w)
            return
        }
        // set user as context value
        ctx := context.WithValue(r.Context(), UserKey, user.ID)
        // call the original handler with the updated context
        handlerFunc(w, r.WithContext(ctx))

	}
}

func forbidden(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, "Unauthorized",)
}


func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
    if !ok {
        return -1
    }
    return userID
}