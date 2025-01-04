package main

import (
	"embed"
	"fmt"
	"net/http"
)

func ServeStaticContent(fs embed.FS, fileName string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		file, err := fs.ReadFile(fileName)
		if err != nil {
			rw.Write([]byte(fmt.Sprintf("Error opening file: %s\n", err)))
		}
		rw.Write(file)
	}
}

func TokenValidationMiddlewareHandler(next http.Handler, jwtSecret []byte, cookieName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieName)
		if err != nil {
			fmt.Fprintf(w, "Error retreiving cookie: %+v\n", err)
			return
		}
		token, _, err := validateJwt(tokenCookie.Value)
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
			fmt.Fprintf(w, "Unable to fetch cookie:\n\t%+v", err)
			return
		}
		token, claims, err := validateJwt(tokenString.Value)
		if err != nil {
			fmt.Fprintf(w, "error validating web token:\n\t%+v", err)
			return
		}
		fmt.Fprintf(w, "Token is: %+v\n", token)
		fmt.Fprintf(w, "Token Validity: %+v\n", token.Valid)
		fmt.Fprintf(w, "Token expiration: %+v\n", token.Claims.(*Claims).ExpiresAt)
		fmt.Fprintf(w, "Claims: %+v\n", token.Claims.(*Claims))
		fmt.Fprintf(w, "JwtKey: %+v\n", jwtKeyMap[claims.KID])
		fmt.Fprintf(w, "Errors: %+v\n", err)
	}
}

func DisplayKeyMapHandler(w http.ResponseWriter, r *http.Request) {
	for key, value := range jwtKeyMap {
		fmt.Fprintf(w, "KEY: %s\n", key)
		fmt.Fprintf(w, "VALUE:\t%+v\n\n", value)
	}
}

func ProtectedRouteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contratulations, the token was valid and you now have access to protected resources")
}
