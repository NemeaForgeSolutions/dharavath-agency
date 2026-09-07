package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"dharavath-agency/configs"
	deliveryHttp "dharavath-agency/internal/delivery/http"
	"dharavath-agency/internal/delivery/http/handler"
	"dharavath-agency/internal/repository/memory"
	"dharavath-agency/internal/service"
	"dharavath-agency/internal/view"
)

func setupTestApp(t *testing.T) http.Handler {
	cfg := configs.LoadConfig()

	propRepo, err := memory.NewPropertyRepo()
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}
	leadRepo := memory.NewLeadRepo()

	propService := service.NewPropertyService(propRepo, propRepo)
	leadService := service.NewLeadService(leadRepo)
	catalogService := service.NewCatalogService(propRepo)

	viewEngine, err := view.NewEngine()
	if err != nil {
		t.Fatalf("Failed to init view engine: %v", err)
	}

	pageHandler := handler.NewPageHandler(propService, catalogService, viewEngine, cfg.Company)
	propertyHandler := handler.NewPropertyHandler(propService, viewEngine, cfg.Company)
	leadHandler := handler.NewLeadHandler(leadService)
	adminService := service.NewAdminService(propRepo, leadRepo)
	adminHandler := handler.NewAdminHandler(propService, adminService, viewEngine, cfg.Company)
	seoHandler := handler.NewSEOHandler(propService, catalogService, "https://dharavathagency.in")

	return deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		PageHandler:     pageHandler,
		PropertyHandler: propertyHandler,
		LeadHandler:     leadHandler,
		AdminHandler:    adminHandler,
		SEOHandler:      seoHandler,
	})
}

func TestEnterprisePages(t *testing.T) {
	app := setupTestApp(t)

	routes := []string{
		"/",
		"/properties",
		"/properties/the-aurora-residences-4bhk-financial-district",
		"/locations",
		"/commercial",
		"/rent",
		"/sell",
		"/new-projects",
		"/agents",
		"/insights",
		"/about",
		"/contact",
		"/static/css/styles.css",
	}

	for _, route := range routes {
		// Test standard request
		req := httptest.NewRequest("GET", route, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Route %s returned status %d, expected %d", route, rec.Code, http.StatusOK)
		}

		// Test gzip request (browser default)
		reqGzip := httptest.NewRequest("GET", route, nil)
		reqGzip.Header.Set("Accept-Encoding", "gzip, deflate, br")
		recGzip := httptest.NewRecorder()
		app.ServeHTTP(recGzip, reqGzip)

		if recGzip.Code != http.StatusOK {
			t.Errorf("Route %s with gzip returned status %d, expected %d", route, recGzip.Code, http.StatusOK)
		}

		// Verify Content-Type is NEVER application/x-gzip
		contentType := recGzip.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/x-gzip") {
			t.Errorf("Route %s returned application/x-gzip instead of actual MIME type: %s", route, contentType)
		}

		if route == "/static/css/styles.css" {
			if !strings.Contains(contentType, "text/css") {
				t.Errorf("CSS route returned wrong Content-Type: %s", contentType)
			}
		} else {
			if !strings.Contains(contentType, "text/html") {
				t.Errorf("HTML route %s returned wrong Content-Type: %s", route, contentType)
			}
		}
	}
}

func TestHTMXSearchEndpoint(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest("GET", "/properties/search?city=Hyderabad&type=Villa", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Jubilee Hills") {
		t.Errorf("Expected Jubilee Hills villa in search response")
	}
}

func TestHTMXLeadSubmission(t *testing.T) {
	app := setupTestApp(t)

	formData := url.Values{
		"name":          {"Deepika Padukone"},
		"phone":         {"+91 98200 55555"},
		"email":         {"deepika@example.com"},
		"propertyId":    {"PROP-HYD-002"},
		"propertyTitle": {"Serene Palms Luxury Triplex Villa"},
		"type":          {"site_visit"},
	}

	req := httptest.NewRequest("POST", "/api/leads", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Inquiry Confirmed, Deepika Padukone!") {
		t.Errorf("Expected confirmation card for Deepika Padukone")
	}
}

func TestSEORoutes(t *testing.T) {
	app := setupTestApp(t)

	// 1. Test sitemap.xml
	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for /sitemap.xml, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") {
		t.Errorf("Expected XML Content-Type for /sitemap.xml, got %s", contentType)
	}

	sitemapBody := rec.Body.String()
	if !strings.Contains(sitemapBody, "<urlset") || !strings.Contains(sitemapBody, "https://dharavathagency.in") {
		t.Errorf("Expected sitemap XML to contain <urlset> and base domain")
	}

	// 2. Test robots.txt
	req = httptest.NewRequest("GET", "/robots.txt", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for /robots.txt, got %d", rec.Code)
	}

	robotsBody := rec.Body.String()
	if !strings.Contains(robotsBody, "Disallow: /admin/") {
		t.Errorf("Expected robots.txt to disallow /admin/")
	}
	if !strings.Contains(robotsBody, "Sitemap: https://dharavathagency.in/sitemap.xml") {
		t.Errorf("Expected robots.txt to link to sitemap.xml")
	}
}

func TestAdminDashboard(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest("GET", "/admin", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for /admin, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Property Portfolio Management") {
		t.Errorf("Expected dashboard header in /admin")
	}
	if !strings.Contains(body, "Total Listings") {
		t.Errorf("Expected KPI cards in /admin")
	}
}

func TestAdminPropertyCRUD(t *testing.T) {
	app := setupTestApp(t)

	// 1. GET /admin/properties/new form
	req := httptest.NewRequest("GET", "/admin/properties/new", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for /admin/properties/new, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Create New Real Estate Listing") {
		t.Errorf("Expected form title in new property page")
	}

	// 2. POST /admin/properties/new (Create)
	newPropForm := url.Values{
		"title":       {"Skyline Heights Penthouse"},
		"transaction": {"Sale"},
		"type":        {"Penthouse"},
		"subType":     {"Duplex Sky Mansion"},
		"developer":   {"Prestige Group"},
		"city":        {"Bengaluru"},
		"location":    {"Whitefield"},
		"state":       {"Karnataka"},
		"facing":      {"East"},
		"floor":       {"32nd Floor"},
		"parking":     {"3 Covered"},
		"price":       {"65000000"},
		"bedrooms":    {"4"},
		"bathrooms":   {"5"},
		"area":        {"4800"},
		"carpetArea":  {"3900"},
		"furnishing":  {"Fully Furnished"},
		"possession":  {"Ready to Move"},
		"reraStatus":  {"RERA Approved"},
		"reraNumber":  {"PRM/KA/RERA/1251"},
		"image":       {"https://images.unsplash.com/photo-1600585154340-be6161a56a0c"},
		"description": {"Exclusive duplex sky mansion with private pool."},
		"amenities":   {"Private Pool, Concierge, Sky Lounge"},
		"featured":    {"true"},
		"verified":    {"true"},
	}

	req = httptest.NewRequest("POST", "/admin/properties/new", strings.NewReader(newPropForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("Expected redirect 303 after property creation, got %d", rec.Code)
	}

	// 3. Verify it shows up on dashboard
	req = httptest.NewRequest("GET", "/admin", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), "Skyline Heights Penthouse") {
		t.Errorf("Expected newly created property to appear on admin dashboard")
	}

	// 4. Test toggle featured via HTMX on existing property PROP-HYD-001
	req = httptest.NewRequest("POST", "/admin/properties/toggle-featured/PROP-HYD-001", nil)
	req.Header.Set("HX-Request", "true")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for toggle featured HTMX, got %d", rec.Code)
	}
	toggleBody := rec.Body.String()
	if !strings.Contains(toggleBody, "hx-post=\"/admin/properties/toggle-featured/PROP-HYD-001\"") {
		t.Errorf("Expected rendered toggle badge in HTMX response")
	}

	// 5. Test edit form GET
	req = httptest.NewRequest("GET", "/admin/properties/edit/PROP-HYD-001", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for edit form, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Edit Listing") {
		t.Errorf("Expected Edit Listing title in edit form")
	}

	// 6. Test delete property via HTMX
	req = httptest.NewRequest("POST", "/admin/properties/delete/PROP-HYD-001", nil)
	req.Header.Set("HX-Request", "true")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for HTMX delete, got %d", rec.Code)
	}

	// Verify it was removed
	req = httptest.NewRequest("GET", "/admin", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "PROP-HYD-001") {
		t.Errorf("Property PROP-HYD-001 should have been deleted from catalog")
	}
}

func TestThemeAndSEOAssets(t *testing.T) {
	app := setupTestApp(t)

	// 1. Verify /favicon.ico redirects to /static/img/favicon.svg
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("Expected 301 for /favicon.ico, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/static/img/favicon.svg" {
		t.Errorf("Expected /favicon.ico location '/static/img/favicon.svg', got %q", loc)
	}

	// 2. Verify /site.webmanifest redirects to /static/site.webmanifest
	req = httptest.NewRequest("GET", "/site.webmanifest", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("Expected 301 for /site.webmanifest, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/static/site.webmanifest" {
		t.Errorf("Expected /site.webmanifest location '/static/site.webmanifest', got %q", loc)
	}

	// 3. Verify static SVG favicon is served
	req = httptest.NewRequest("GET", "/static/img/favicon.svg", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 for /static/img/favicon.svg, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Errorf("Expected SVG content in favicon")
	}

	// 4. Verify static webmanifest is valid JSON
	req = httptest.NewRequest("GET", "/static/site.webmanifest", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 for /static/site.webmanifest, got %d", rec.Code)
	}
	var manifestData map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &manifestData); err != nil {
		t.Errorf("site.webmanifest is not valid JSON: %v", err)
	}

	// 5. Verify Homepage contains theme initialization (anti-FOUC) and Schema.org
	req = httptest.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	homeHtml := rec.Body.String()
	if !strings.Contains(homeHtml, "dharavath_theme") {
		t.Errorf("Expected anti-FOUC theme script in homepage")
	}
	if !strings.Contains(homeHtml, "/static/css/styles.css") {
		t.Errorf("Expected /static/css/styles.css compiled Tailwind stylesheet in homepage")
	}
	if !strings.Contains(homeHtml, "darkMode") {
		t.Errorf("Expected Alpine darkMode state in homepage")
	}
	if !strings.Contains(homeHtml, "RealEstateAgent") {
		t.Errorf("Expected RealEstateAgent Schema.org structured data in homepage")
	}

	// 6. Verify Property Detail page contains valid JSON-LD RealEstateListing
	req = httptest.NewRequest("GET", "/properties/the-aurora-residences-4bhk-financial-district", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	propHtml := rec.Body.String()
	if !strings.Contains(propHtml, "RealEstateListing") {
		t.Errorf("Expected RealEstateListing schema in property detail page")
	}

	// Extract and validate JSON-LD
	startTag := `<script type="application/ld+json">`
	endTag := `</script>`
	startIndex := strings.Index(propHtml, startTag)
	if startIndex == -1 {
		t.Fatalf("JSON-LD script tag not found on property detail page")
	}
	startIndex += len(startTag)
	endIndex := strings.Index(propHtml[startIndex:], endTag)
	if endIndex == -1 {
		t.Fatalf("JSON-LD closing script tag not found")
	}
	jsonLdStr := strings.TrimSpace(propHtml[startIndex : startIndex+endIndex])

	var schemaObj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLdStr), &schemaObj); err != nil {
		t.Errorf("Property detail JSON-LD is not valid JSON: %v\nJSON content was:\n%s", err, jsonLdStr)
	}

	// 7. Verify Contact page contains ContactPage Schema.org
	req = httptest.NewRequest("GET", "/contact", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	contactHtml := rec.Body.String()
	if !strings.Contains(contactHtml, "ContactPage") {
		t.Errorf("Expected ContactPage schema in contact page")
	}

	// 8. Verify /health returns 200 OK with json
	req = httptest.NewRequest("GET", "/health", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK from /health, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("Expected status:ok from /health, got %s", rec.Body.String())
	}
}

