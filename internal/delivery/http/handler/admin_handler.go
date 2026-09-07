package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"dharavath-agency/internal/domain"
	"dharavath-agency/internal/service"
	"dharavath-agency/internal/view"
)

type AdminHandler struct {
	propertyService *service.PropertyService
	adminService    *service.AdminService
	viewEngine      *view.Engine
	company         domain.CompanyInfo
}

func NewAdminHandler(
	pService *service.PropertyService,
	aService *service.AdminService,
	vEngine *view.Engine,
	company domain.CompanyInfo,
) *AdminHandler {
	return &AdminHandler{
		propertyService: pService,
		adminService:    aService,
		viewEngine:      vEngine,
		company:         company,
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats := h.adminService.GetStats()
	properties := h.propertyService.ListProperties(domain.PropertyFilter{})
	leads := h.adminService.GetRecentLeads()

	pageData := view.PageData{
		Title:       "Admin Portal - Real Estate Management",
		Description: "Dharavath Agency Internal Inventory and Advisory Management System",
		ActivePage:  "admin-dashboard",
		Company:     h.company,
		Data: map[string]interface{}{
			"Stats":      stats,
			"Properties": properties,
			"Leads":      leads,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "dashboard.html", pageData); err != nil {
		log.Printf("Error rendering admin dashboard: %v", err)
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) NewProperty(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:       "Add New Property Listing",
		Description: "Add a new luxury real estate listing to the catalog",
		ActivePage:  "admin-new-property",
		Company:     h.company,
		Data: map[string]interface{}{
			"IsEdit": false,
			"Property": &domain.Property{
				City:        "Hyderabad",
				State:       "Telangana",
				Transaction: "Sale",
				Type:        "Apartment",
				Facing:      "East",
				Possession:  "Ready to Move",
				Furnishing:  "Semi-Furnished",
				ReraStatus:  "RERA Approved",
				Verified:    true,
				Bedrooms:    3,
				Bathrooms:   3,
			},
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "property_form.html", pageData); err != nil {
		log.Printf("Error rendering new property form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	prop := h.parsePropertyFromForm(r)

	if err := h.propertyService.Create(&prop); err != nil {
		pageData := view.PageData{
			Title:      "Add New Property Listing",
			ActivePage: "admin-new-property",
			Company:    h.company,
			Data: map[string]interface{}{
				"IsEdit":   false,
				"Property": &prop,
				"Error":    err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "property_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) EditProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	prop, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pageData := view.PageData{
		Title:       "Edit Property: " + prop.Title,
		Description: "Modify existing listing details",
		ActivePage:  "admin-dashboard",
		Company:     h.company,
		Data: map[string]interface{}{
			"IsEdit":   true,
			"Property": prop,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "property_form.html", pageData); err != nil {
		log.Printf("Error rendering edit property form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	updated := h.parsePropertyFromForm(r)
	updated.ID = existing.ID
	if updated.Slug == "" {
		updated.Slug = existing.Slug
	}
	if len(updated.Gallery) == 0 {
		updated.Gallery = existing.Gallery
	}
	if len(updated.Nearby) == 0 {
		updated.Nearby = existing.Nearby
	}
	if updated.AgentID == "" {
		updated.AgentID = existing.AgentID
	}

	if err := h.propertyService.Update(&updated); err != nil {
		pageData := view.PageData{
			Title:      "Edit Property: " + updated.Title,
			ActivePage: "admin-dashboard",
			Company:    h.company,
			Data: map[string]interface{}{
				"IsEdit":   true,
				"Property": &updated,
				"Error":    err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "property_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.propertyService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if view.IsHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) ToggleFeatured(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := h.propertyService.ToggleFeatured(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prop, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if view.IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := h.viewEngine.RenderPartial(w, "admin_featured_badge", prop); err != nil {
			log.Printf("Error rendering featured badge partial: %v", err)
			http.Error(w, "Render error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) parsePropertyFromForm(r *http.Request) domain.Property {
	price, _ := strconv.ParseInt(strings.TrimSpace(r.FormValue("price")), 10, 64)
	bedrooms, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("bedrooms")))
	bathrooms, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("bathrooms")))
	area, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("area")))
	carpetArea, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("carpetArea")))

	var amenities []string
	rawAmenities := r.FormValue("amenities")
	if rawAmenities != "" {
		for _, a := range strings.Split(rawAmenities, ",") {
			cleaned := strings.TrimSpace(a)
			if cleaned != "" {
				amenities = append(amenities, cleaned)
			}
		}
	}

	return domain.Property{
		ID:           strings.TrimSpace(r.FormValue("id")),
		Title:        strings.TrimSpace(r.FormValue("title")),
		Transaction:  strings.TrimSpace(r.FormValue("transaction")),
		Type:         strings.TrimSpace(r.FormValue("type")),
		SubType:      strings.TrimSpace(r.FormValue("subType")),
		Developer:    strings.TrimSpace(r.FormValue("developer")),
		City:         strings.TrimSpace(r.FormValue("city")),
		Location:     strings.TrimSpace(r.FormValue("location")),
		State:        strings.TrimSpace(r.FormValue("state")),
		Facing:       strings.TrimSpace(r.FormValue("facing")),
		Floor:        strings.TrimSpace(r.FormValue("floor")),
		Parking:      strings.TrimSpace(r.FormValue("parking")),
		Price:        price,
		PriceDisplay: strings.TrimSpace(r.FormValue("priceDisplay")),
		Deposit:      strings.TrimSpace(r.FormValue("deposit")),
		Bedrooms:     bedrooms,
		Bathrooms:    bathrooms,
		Area:         area,
		CarpetArea:   carpetArea,
		AreaUnit:     "sq.ft.",
		Furnishing:   strings.TrimSpace(r.FormValue("furnishing")),
		Possession:   strings.TrimSpace(r.FormValue("possession")),
		ReraStatus:   strings.TrimSpace(r.FormValue("reraStatus")),
		ReraNumber:   strings.TrimSpace(r.FormValue("reraNumber")),
		Image:        strings.TrimSpace(r.FormValue("image")),
		Description:  strings.TrimSpace(r.FormValue("description")),
		Amenities:    amenities,
		Featured:     r.FormValue("featured") == "true",
		Verified:     r.FormValue("verified") == "true",
		IsNewLaunch:  r.FormValue("isNewLaunch") == "true",
	}
}
