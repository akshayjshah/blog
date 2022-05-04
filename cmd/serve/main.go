package main

import (
	"context"
	"embed"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

//go:embed static recipes favicon.ico *.html
var site embed.FS

var redirects = map[string]string{
	"/books/grit/":      "/grit/",
	"/books/sourdough/": "/sourdough/",
}

func serve(w http.ResponseWriter, r *http.Request) {
	if to, ok := redirects[r.URL.Path]; ok {
		http.Redirect(w, r, to, http.StatusFound)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/") {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		} else {
			path = strings.TrimSuffix(path, "/") + ".html"
		}
		page, err := site.Open(path)
		if err != nil {
			http.Redirect(w, r, "/404/", http.StatusFound)
			return
		}
		io.Copy(w, page)
		return
	}
	w.Header().Set("Cache-Control", "max-age=31536000") // cache for 1y
	http.FileServer(http.FS(site)).ServeHTTP(w, r)
}

func main() {
	hostport := "localhost:8080"
	if p := os.Getenv("PORT"); p != "" {
		hostport = ":" + p
	}
	handler := h2c.NewHandler(
		http.HandlerFunc(serve),
		&http2.Server{
			IdleTimeout: 30 * time.Second,
		},
	)
	srv := &http.Server{
		Addr:              hostport,
		Handler:           handler,
		ReadHeaderTimeout: 100 * time.Millisecond,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown: %v", err)
		}
		close(done)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen and serve: %v", err)
	}
	<-done
}
