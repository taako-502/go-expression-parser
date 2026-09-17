package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/taako-502/go-expression-parser/parser"
)

//go:embed web/index.html
var indexHTML []byte

type response struct {
	Tokens []parser.Token `json:"tokens,omitempty"`
	AST    *parser.Node   `json:"ast,omitempty"`
	Result *float64       `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	})
	mux.HandleFunc("GET /api/evaluate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		ast, tokens, err := parser.Parse(r.URL.Query().Get("expression"))
		out := response{Tokens: tokens, AST: ast}
		if err == nil {
			var value float64
			value, err = parser.Eval(ast)
			if err == nil {
				out.Result = &value
			}
		}
		if err != nil {
			out.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
		}
		_ = json.NewEncoder(w).Encode(out)
	})
	return mux
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "HTTP server address")
	flag.Parse()
	server := &http.Server{Addr: *addr, Handler: handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("ブラウザで http://%s を開いてください（終了: Ctrl+C）", *addr)
	log.Fatal(server.ListenAndServe())
}
