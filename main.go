
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
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
    // Forward relevant headers from the incoming request to the upstream,
    // excluding hop-by-hop headers and Authorization (Digest lib sets it).
    if d.Header == nil {
        d.Header = make(http.Header)
    }
    hopByHop := map[string]struct{}{
        "Connection":        {},
        "Proxy-Connection":  {},
        "Keep-Alive":        {},
        "Proxy-Authenticate":{},
        "Proxy-Authorization":{},
        "Te":                {},
        "Trailer":           {},
        "Transfer-Encoding": {},
        "Upgrade":           {},
        "Authorization":     {}, // will be set by digest client
        "Host":              {}, // set by http client
        "Content-Length":    {}, // computed by http client
    }
    for key, values := range req.Header {
        if _, skip := hopByHop[key]; skip {
            continue
        }
        // Go's http.Header is case-insensitive, but ensure standard casing
        // when setting on the outbound request.
        canonicalKey := http.CanonicalHeaderKey(key)
        for _, v := range values {
            // Some servers are strict about empty header values; skip empties
            if strings.TrimSpace(v) == "" {
                continue
            }
            d.Header.Add(canonicalKey, v)
        }
    }
	resp, err := d.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Digest request failed: %v\n", err)
		http.Error(w, "Upstream request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

    // Copy upstream response headers back to the client, excluding hop-by-hop
    // headers which the Go server manages itself.
    respHopByHop := map[string]struct{}{
        "Connection":        {},
        "Proxy-Connection":  {},
        "Keep-Alive":        {},
        "Proxy-Authenticate":{},
        "Proxy-Authorization":{},
        "Te":                {},
        "Trailer":           {},
        "Transfer-Encoding": {},
        "Upgrade":           {},
    }
    for key, values := range resp.Header {
        if _, skip := respHopByHop[key]; skip {
            continue
        }
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
