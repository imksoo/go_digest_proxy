
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	digest "github.com/xinsnake/go-http-digest-auth-client"
)

func handler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	user := os.Getenv("DIGEST_USER")
	pass := os.Getenv("DIGEST_PASS")
	backend := os.Getenv("BACKEND_URL")

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to read request body: %v\n", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	targetURL := backend + req.URL.Path
	if req.URL.RawQuery != "" {
		targetURL += "?" + req.URL.RawQuery
	}

	d := digest.NewRequest(user, pass, req.Method, targetURL, string(bodyBytes))
	resp, err := d.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Digest request failed: %v\n", err)
		http.Error(w, "Upstream request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	duration := time.Since(start)
	log.Printf("%s %s -> %s %d (%s)", req.Method, req.URL.Path, targetURL, resp.StatusCode, duration)
}

func main() {
	_ = godotenv.Load()
	listen := os.Getenv("PORT")
	if listen == "" {
		listen = "8080"
	}

	http.HandleFunc("/", handler)
	log.Printf("Proxy server listening on :%s", listen)
	log.Fatal(http.ListenAndServe(":"+listen, nil))
}

