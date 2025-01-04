package main

import (
	"fmt"
	"log"
	"net/http"
)

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
	http.HandleFunc("/jwt/token/keymap", TokenValidationMiddleware(http.HandlerFunc(DisplayKeyMapHandler), secret, "set_cookie"))
	http.HandleFunc("/index", ServeEmbeddedStaticContent(files, "files/index.html"))
	http.Handle("/jwt/token/protected2", TokenValidationMiddlewareHandler(http.HandlerFunc(ProtectedRouteHandler), secret, "set_cookie"))
	fmt.Printf("Starting server on port 8080")
	serveErr := http.ListenAndServe(":8080", nil)
	if serveErr != nil {
		fmt.Printf("Error starting server: %v", serveErr)
	}
}
