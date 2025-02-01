package main

import (
	"fmt"
	"net/http"
)

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\nThis is a test page!\n")
}

func main() {
	http.HandleFunc("/test", testResponse)
	http.HandleFunc("/jwt/token/get", CreateWebTokenHandler())
	http.HandleFunc("/jwt/token/validate", ValidateWebTokenHandlerDebugger())
	http.HandleFunc("/jwt/token/keymap", TokenValidationMiddleware(http.HandlerFunc(DisplayKeyMapHandler), "set_cookie"))
	http.HandleFunc("/index", ServeEmbeddedStaticContent(files, "files/index.html"))
	http.HandleFunc("/jwt/token/protected", TokenValidationMiddleware(http.HandlerFunc(ProtectedRouteHandler), "set_cookie"))
	http.Handle("/jwt/token/protected2", TokenValidationMiddlewareHandler(http.HandlerFunc(ProtectedRouteHandler), "set_cookie"))
	fmt.Printf("Starting server on port 8080\n")
	serveErr := http.ListenAndServe(":8080", nil)
	if serveErr != nil {
		fmt.Printf("Error starting server: %v", serveErr)
	}
}
