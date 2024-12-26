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

func validateJwt(tokenString string, jwtSecret []byte) (token *jwt.Token, validErr error) {
	// claims := &jwt.RegisteredClaims{}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (key any, keyErr error) {
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (key any, keyErr error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			keyErr = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			key = jwtSecret
		}
		return key, keyErr
	})
	if err != nil {
		validErr = err
	}
	return token, validErr
}

func TokenValidationMiddlewareHandler(next http.Handler, jwtSecret []byte, cookieName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieName)
		if err != nil {
			fmt.Fprintf(w, "Error retreiving cookie: %+v\n", err)
			return
		}
		token, err := validateJwt(tokenCookie.Value, jwtSecret)
		if err != nil {
			fmt.Fprintf(w, "Error validating cookie: %+v\n", err)
			return
		}
		if !token.Valid {
			fmt.Fprintf(w, "Error code: %d - Unauthorized - Token Invalid\n", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func TokenValidationMiddleware(next http.Handler, jwtSecret []byte, cookieName string) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieName)
		if err != nil {
			fmt.Fprintf(w, "Error retreiving cookie: %+v\n", err)
			return
		}
		token, err := validateJwt(tokenCookie.Value, jwtSecret)
		if err != nil {
			fmt.Fprintf(w, "Error validating cookie: %+v\n", err)
			return
		}
		if !token.Valid {
			fmt.Fprintf(w, "Error code: %d - Unauthorized - Token Invalid\n", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateWebTokenHandlerDebugger(jwtSecret []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("set_cookie")
		if err != nil {
			fmt.Fprintf(w, "Unable to fetch cookie\n")
		}
		token, err := validateJwt(tokenString.Value, jwtSecret)
		fmt.Fprintf(w, "Token is: %+v\n", *token)
		fmt.Fprintf(w, "Token Validity: %+v\n", token.Valid)
		fmt.Fprintf(w, "Token expiration: %+v\n", token.Claims.(*Claims).ExpiresAt)
		fmt.Fprintf(w, "Claims: %+v\n", token.Claims.(*Claims))
		fmt.Fprintf(w, "Errors: %+v\n", err)
	}
}

func CreateWebTokenHandler(jwtSecret []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := CreateWebToken()
		signedToken, signErr := SignWebToken(token, jwtSecret)
		if signErr != nil {
			fmt.Fprintf(w, "Error signing token: %v\n", signErr)
		}
		tokenExpiration, err := token.Claims.GetExpirationTime()
		if err != nil {
			fmt.Fprintf(w, "Error getting expiration time: %s", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "set_cookie",
			Value:   signedToken,
			Path:    "/",
			Expires: tokenExpiration.Time,
		})
		fmt.Fprintln(w, signedToken)
	}
}

func ProtectedRouteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contratulations, the token was valid and you now have access to protected resources")
}

func CreateWebToken() *jwt.Token {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}

func SignWebToken(token *jwt.Token, jwtSecret any) (signedToken string, signErr error) {
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		signErr = fmt.Errorf("error signing token: %v", err)
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
	http.HandleFunc("/jwt/token/get", CreateWebTokenHandler(secret))
	http.HandleFunc("/jwt/token/validate", ValidateWebTokenHandlerDebugger(secret))
	http.HandleFunc("/jwt/token/protected", TokenValidationMiddleware(http.HandlerFunc(ProtectedRouteHandler), secret, "set_cookie"))
	http.Handle("/jwt/token/protected", TokenValidationMiddlewareHandler(http.HandlerFunc(ProtectedRouteHandler), secret, "set_cookie"))
	fmt.Printf("Starting server on port 8080")
	serveErr := http.ListenAndServe(":8080", nil)
	if serveErr != nil {
		fmt.Printf("Error starting server: %v", serveErr)
	}
}
