package http

import (
	"io/fs"
	"log"
	"net/http"

	"dharavath-agency/internal/delivery/http/handler"
	"dharavath-agency/internal/delivery/http/middleware"
	"dharavath-agency/web"
)

type RouterConfig struct {
	PageHandler     *handler.PageHandler
	PropertyHandler *handler.PropertyHandler
	LeadHandler     *handler.LeadHandler
	AdminHandler    *handler.AdminHandler
	SEOHandler      *handler.SEOHandler
}

func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	staticSubFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		log.Fatalf("Failed to create static sub-filesystem: %v", err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))

	if cfg.SEOHandler != nil {
		mux.HandleFunc("GET /sitemap.xml", cfg.SEOHandler.SitemapXML)
		mux.HandleFunc("GET /robots.txt", cfg.SEOHandler.RobotsTxt)
	}
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/img/favicon.svg", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/site.webmanifest", http.StatusMovedPermanently)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /{$}", cfg.PageHandler.Home)
	mux.HandleFunc("GET /properties", cfg.PropertyHandler.List)
	mux.HandleFunc("GET /properties/search", cfg.PropertyHandler.Search)
	mux.HandleFunc("GET /properties/featured", cfg.PropertyHandler.Featured)
	mux.HandleFunc("GET /properties/{slug}", cfg.PropertyHandler.Detail)
	mux.HandleFunc("GET /locations", cfg.PageHandler.Locations)
	mux.HandleFunc("GET /commercial", cfg.PageHandler.Commercial)
	mux.HandleFunc("GET /rent", cfg.PageHandler.Rent)
	mux.HandleFunc("GET /sell", cfg.PageHandler.Sell)
	mux.HandleFunc("GET /new-projects", cfg.PageHandler.NewProjects)
	mux.HandleFunc("GET /agents", cfg.PageHandler.Agents)
	mux.HandleFunc("GET /insights", cfg.PageHandler.Insights)
	mux.HandleFunc("GET /about", cfg.PageHandler.About)
	mux.HandleFunc("GET /contact", cfg.PageHandler.Contact)

	if cfg.AdminHandler != nil {
		mux.HandleFunc("GET /admin", cfg.AdminHandler.Dashboard)
		mux.HandleFunc("GET /admin/properties", cfg.AdminHandler.Dashboard)
		mux.HandleFunc("GET /admin/properties/new", cfg.AdminHandler.NewProperty)
		mux.HandleFunc("POST /admin/properties/new", cfg.AdminHandler.CreateProperty)
		mux.HandleFunc("GET /admin/properties/edit/{id}", cfg.AdminHandler.EditProperty)
		mux.HandleFunc("POST /admin/properties/edit/{id}", cfg.AdminHandler.UpdateProperty)
		mux.HandleFunc("POST /admin/properties/delete/{id}", cfg.AdminHandler.DeleteProperty)
		mux.HandleFunc("POST /admin/properties/toggle-featured/{id}", cfg.AdminHandler.ToggleFeatured)
	}

	mux.HandleFunc("POST /api/leads", cfg.LeadHandler.Submit)

	var handler http.Handler = mux
	handler = middleware.Compress(handler)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.Logger(handler)
	handler = middleware.Recovery(handler)

	return handler
}
