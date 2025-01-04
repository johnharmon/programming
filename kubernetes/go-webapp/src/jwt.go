package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKeyMap = make(map[string]*JwtKey)

func generateJWTSecret() (key []byte, keyErr error) {
	jwtSecret := make([]byte, 32)

	if _, randErr := rand.Read(jwtSecret); randErr != nil {
		return key, fmt.Errorf("failed to generate secret!\n%v", randErr)
	}
	return jwtSecret, nil
}

type Claims struct {
	Username string `json:"username"`
	KID      string `json:"kid"`
	jwt.RegisteredClaims
}

func CreateWebToken(uuidString string) *jwt.Token {
	//expirationTime := time.Now().Add(time.Hour * 1)
	expirationTime := DefaultConfig.JWTConfig.Token.Expiration.GetExpirationTime()
	claims := &Claims{
		Username: "test-user",
		KID:      uuidString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    DefaultConfig.JWTConfig.Token.Issuer,
			Audience:  jwt.ClaimStrings{DefaultConfig.JWTConfig.Token.Audience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}

func SignWebTokenWithSecret(token *jwt.Token, jwtSecret any) (signedToken string, signErr error) {
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		signErr = fmt.Errorf("error signing token: %v", err)
	}
	return signedToken, signErr
}

func SignWebToken(token *jwt.Token) (signedToken string, jwtSecret []byte, signErr error) {
	jwtSecret, err := generateJWTSecret()
	if err != nil {
		signErr = fmt.Errorf("error creating secret: %+v", err)
	}
	signedToken, err = token.SignedString(jwtSecret)
	if err != nil {
		signErr = errors.Join(signErr, fmt.Errorf("error signing token: %v", err))
	}
	return signedToken, jwtSecret, signErr
}

func TokenValidationMiddleware(next http.Handler, jwtSecret []byte, cookieName string) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieName)
		if err != nil {
			http.Error(w, "Missing req1uired cookie\n", http.StatusNotAcceptable)
			fmt.Fprintf(w, "Error retreiving cookie: %+v\n", err)
			return
		}
		token, _, err := validateJwt(tokenCookie.Value)
		if err != nil {
			http.Error(w, "Error - Cannot validate token\n", http.StatusBadRequest)
			return
		}
		if !token.Valid {
			http.Error(w, "Unauthorized - Token Invalid\n", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CreateWebTokenHandler(jwtSecret []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenUUID := uuid.NewString()
		token := CreateWebToken(tokenUUID)
		signedToken, tokenSecret, signErr := SignWebToken(token)
		if signErr != nil {
			fmt.Fprintf(w, "Error signing token: %v\n", signErr)
		}
		_, err := token.Claims.GetExpirationTime()
		if err != nil {
			fmt.Fprintf(w, "Error getting expiration time: %s", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "set_cookie",
			Value: signedToken,
			Path:  "/",
			//Expires: tokenExpiration.Time,
			//Expires: time.Now().Add(time.Hour),
			Expires: DefaultConfig.JWTConfig.Cookie.GetExpirationTime(),
		})
		fmt.Fprintln(w, signedToken)
		jwtKey := NewJwtKeyWithUUID(tokenSecret, tokenUUID)
		jwtKey.GetSecret()
		fmt.Printf("%+v\n", jwtKey)
		jwtKeyMap[tokenUUID] = jwtKey
	}
}

func validateJwt(tokenString string) (token *jwt.Token, claims *Claims, validErr error) {
	claims = &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (key any, keyErr error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			keyErr = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			jwtKeyMap[claims.KID].GetSecret()
			key = jwtKeyMap[claims.KID].KeySecret
		}
		return key, keyErr
	})
	if err != nil {
		validErr = err
	}
	return token, claims, validErr
}
