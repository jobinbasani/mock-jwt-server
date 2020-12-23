package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/ztrue/shutdown"
	"log"
	"mock-jwt-server/handlers"
	"net/http"
	"os"
	"time"
)

const keyId = "mock-key"

func main() {
	port := flag.Int("port", 8994, "Port to run the server")
	flag.Parse()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return
	}

	jwkPrivateKey, err := jwk.New(privateKey)
	if err != nil {
		log.Fatal("Error loading private JWK")
	}
	jwkPublicKey, err := jwk.New(privateKey.PublicKey)
	if err != nil {
		log.Fatal("Error loading public JWK")
	}
	_ = jwkPrivateKey.Set(jwk.KeyIDKey, keyId)
	_ = jwkPublicKey.Set(jwk.KeyIDKey, keyId)

	startServer(port, jwkPrivateKey, jwkPublicKey)
}

func startServer(port *int, privateKey, publicKey jwk.Key) {
	r := mux.NewRouter()

	r.HandleFunc("/.well-known/jwks.json", handlers.GetJwks(publicKey)).Methods(http.MethodGet)
	r.HandleFunc("/token/{sub}", handlers.GetTokenForSub(privateKey)).Methods(http.MethodGet)
	r.HandleFunc("/token", handlers.GenerateTokenWithPayload(privateKey)).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
		Handler: r,
	}

	go func() {
		log.Println(fmt.Sprintf("Starting server on port %d", *port))
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	shutdown.Add(func() {
		log.Println("Shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		_ = srv.Shutdown(ctx)
	})
	shutdown.Listen(os.Interrupt)
}
