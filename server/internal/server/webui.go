package server

import (
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"fhs/app/internal/config"

	"github.com/go-fuego/fuego"
)

func RegisterWebUI(cfg *config.Config, s *fuego.Server) {
	if cfg.AppEnv != "production" {
		return
	}

	if _, err := os.Stat(filepath.Join(cfg.WebRoot, "index.html")); err != nil {
		slog.Warn("web UI disabled: no index.html under web root", "webRoot", cfg.WebRoot)
		return
	}

	s.Mux.Handle("/", spaHandler(cfg.WebRoot))
}

func spaHandler(root string) http.Handler {
	fileServer := http.FileServer(http.Dir(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlPath := path.Clean("/" + r.URL.Path)

		if urlPath == "/api" || strings.HasPrefix(urlPath, "/api/") {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// SvelteKit emits content-hashed assets under _app/immutable
		if strings.HasPrefix(urlPath, "/_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(urlPath))); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(root, "index.html"))
	})
}
