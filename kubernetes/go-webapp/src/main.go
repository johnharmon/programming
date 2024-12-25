package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func generateJWTSecret() (key []byte, keyErr error) {
	jwtSecret := make([]byte, 32)

	if _, randErr := rand.Read(jwtSecret); randErr != nil {
		return key, fmt.Errorf("Failed to generate secret!\n%v", randErr)
	}
	return jwtSecret, nil
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func validateJwt(tokenString string, jwtSecret []byte) (valid bool, token *jwt.Token, validErr error) {
	// claims := &jwt.RegisteredClaims{}
	// claims := &Claims{}
	// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (key any, keyErr error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (key any, keyErr error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			keyErr = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			key = jwtSecret
		}
		return key, keyErr
	})
	if err != nil || !token.Valid {
		valid = false
		validErr = err
	}
	return valid, token, validErr
}

func ValidateWebTokenHandler(jwtSecret []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("set_cookie")
		if err != nil {
			fmt.Fprintf(w, "Unable to fetch cookie\n")
		}
		_, token, err := validateJwt(tokenString.Value, jwtSecret)
		fmt.Fprintf(w, "Token is: %v\n", token)
		// fmt.Fprintf(w, "%v\n", token)
		fmt.Fprintf(w, "Token Validity: %v\n", token.Valid)
		fmt.Fprintf(w, "Errors: %v\n", err)
	}
}

func CreateWebTokenHandler(jwtSecret []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := CreateWebToken()
		signedToken, signErr := SignWebToken(token, jwtSecret)
		if signErr != nil {
			fmt.Fprintf(w, "Error signing token: %v\n", signErr)
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "set_cookie",
			Value: signedToken,
			Path:  "/",
		})
		// w.Header().Add("set_cookie", signedToken)
		fmt.Fprintln(w, signedToken)
	}
}

func CreateWebToken() *jwt.Token {
	expirationTime := time.Now().Add(60 * time.Minute).Unix()
	//	claims := &Claims{
	//		RegisteredClaims: jwt.RegisteredClaims{
	//			ExpiresAt: jwt.NewNumericDate(expirationTime),
	//		},
	//	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "placeholder",
		"exp":      expirationTime,
	})
	// token := jwt.New(jwt.SigningMethodHS256)
	return token
}

func SignWebToken(token *jwt.Token, jwtSecret any) (signedToken string, signErr error) {
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		signErr = fmt.Errorf("Error signing token: %v", err)
	}
	return signedToken, signErr
}

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\nThis is a test page!\n")
}

func main() {
	secret, err := generateJWTSecret()
	if err != nil {
		log.Fatalf("Error generating secret:\n\t%v", err)
	}
	http.HandleFunc("/test", testResponse)
	http.HandleFunc("/jwt/token/get", http.HandlerFunc(CreateWebTokenHandler(secret)))
	http.HandleFunc("/jwt/token/validate", http.HandlerFunc(ValidateWebTokenHandler(secret)))
	fmt.Printf("Starting server on port 8080")
	serveErr := http.ListenAndServe(":8080", nil)
	if serveErr != nil {
		fmt.Printf("Error starting server: %v", serveErr)
	}
}
