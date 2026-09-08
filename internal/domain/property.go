package domain

import "strings"

type NearbyPlace struct {
	Name     string `json:"name"`
	Distance string `json:"distance"`
}

type Property struct {
	ID           string        `json:"id"`
	Slug         string        `json:"slug"`
	Title        string        `json:"title"`
	Type         string        `json:"type"`        // Apartment, Villa, Penthouse, Commercial, Plot
	SubType      string        `json:"subType"`     // High-Rise Sky Villa, Independent Gated Villa, etc.
	Transaction  string        `json:"transaction"` // Sale, Rent
	Price        int64         `json:"price"`
	PriceDisplay string        `json:"priceDisplay"`
	PricePerSqFt string        `json:"pricePerSqFt"`
	Deposit      string        `json:"deposit,omitempty"`
	Bedrooms     int           `json:"bedrooms"`
	Bathrooms    int           `json:"bathrooms"`
	Balconies    int           `json:"balconies"`
	Area         int           `json:"area"`
	CarpetArea   int           `json:"carpetArea"`
	AreaUnit     string        `json:"areaUnit"`
	Location     string        `json:"location"`
	City         string        `json:"city"`
	State        string        `json:"state"`
	Facing       string        `json:"facing"`
	Floor        string        `json:"floor"`
	Parking      string        `json:"parking"`
	Furnishing   string        `json:"furnishing"`
	Possession   string        `json:"possession"`
	PropertyAge  string        `json:"propertyAge"`
	ReraStatus   string        `json:"reraStatus"`
	ReraNumber   string        `json:"reraNumber"`
	Featured     bool          `json:"featured"`
	Verified     bool          `json:"verified"`
	IsNewLaunch  bool          `json:"isNewLaunch"`
	Image        string        `json:"image"`
	Gallery      []string      `json:"gallery"`
	Description  string        `json:"description"`
	Amenities    []string      `json:"amenities"`
	Nearby       []NearbyPlace `json:"nearby"`
	Developer    string        `json:"developer"`
	AgentID      string        `json:"agentId"`
	ProjectID    *string       `json:"projectId,omitempty"`
	LocationID   *string       `json:"locationId,omitempty"`
	CreatedBy    *string       `json:"createdBy,omitempty"`
	UpdatedBy    *string       `json:"updatedBy,omitempty"`
}

func (p *Property) Matches(f PropertyFilter) bool {
	if f.Transaction != "" && !strings.EqualFold(p.Transaction, f.Transaction) {
		return false
	}
	if f.City != "" && !strings.EqualFold(p.City, f.City) {
		return false
	}
	if f.Type != "" && !strings.EqualFold(p.Type, f.Type) {
		return false
	}
	if f.Bedrooms > 0 && p.Bedrooms != f.Bedrooms {
		return false
	}
	if f.MinPrice > 0 && p.Price < f.MinPrice {
		return false
	}
	if f.MaxPrice > 0 && p.Price > f.MaxPrice {
		return false
	}
	if f.Possession != "" && !strings.Contains(strings.ToLower(p.Possession), strings.ToLower(f.Possession)) {
		return false
	}
	if f.Furnishing != "" && !strings.EqualFold(p.Furnishing, f.Furnishing) {
		return false
	}
	if f.ReraOnly && !strings.Contains(strings.ToLower(p.ReraStatus), "rera") {
		return false
	}
	if f.Query != "" {
		q := strings.ToLower(f.Query)
		match := strings.Contains(strings.ToLower(p.Title), q) ||
			strings.Contains(strings.ToLower(p.Location), q) ||
			strings.Contains(strings.ToLower(p.City), q) ||
			strings.Contains(strings.ToLower(p.Developer), q) ||
			strings.Contains(strings.ToLower(p.Description), q)
		if !match {
			return false
		}
	}
	return true
}

type PropertyFilter struct {
	Transaction string
	City        string
	Type        string
	Bedrooms    int
	MinPrice    int64
	MaxPrice    int64
	Possession  string
	Furnishing  string
	ReraOnly    bool
	Query       string
	Sort        string // price_asc, price_desc, area_desc, featured
}

type Project struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Developer      string   `json:"developer"`
	Location       string   `json:"location"`
	LocationID     *string  `json:"locationId,omitempty"`
	LeadAgentID    *string  `json:"leadAgentId,omitempty"`
	StartingPrice  string   `json:"startingPrice"`
	Configuration  string   `json:"configuration"`
	PossessionDate string   `json:"possessionDate"`
	Status         string   `json:"status"` // Pre-Launch, Under Construction, Ready
	Badge          string   `json:"badge"`
	TotalUnits     string   `json:"totalUnits"`
	LandArea       string   `json:"landArea"`
	ReraNumber     string   `json:"reraNumber"`
	Image          string   `json:"image"`
	Overview       string   `json:"overview"`
	Highlights     []string `json:"highlights"`
	CreatedBy      *string  `json:"createdBy,omitempty"`
	UpdatedBy      *string  `json:"updatedBy,omitempty"`
}

type LocationInsight struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Tagline       string   `json:"tagline"`
	AvgPrice      string   `json:"avgPrice"`
	RentalYield   string   `json:"rentalYield"`
	Image         string   `json:"image"`
	KeyLocalities []string `json:"keyLocalities"`
	Description   string   `json:"description"`
	CreatedBy     *string  `json:"createdBy,omitempty"`
	UpdatedBy     *string  `json:"updatedBy,omitempty"`
}

type InsightArticle struct {
	ID            string  `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	Category      string  `json:"category"`
	ReadTime      string  `json:"readTime"`
	Date          string  `json:"date"`
	Summary       string  `json:"summary"`
	Author        string  `json:"author"`
	AuthorAgentID *string `json:"authorAgentId,omitempty"`
	KeyTakeaway   string  `json:"keyTakeaway"`
	CreatedBy     *string `json:"createdBy,omitempty"`
	UpdatedBy     *string `json:"updatedBy,omitempty"`
}

type Testimonial struct {
	ID         string  `json:"id"`
	Quote      string  `json:"quote"`
	Client     string  `json:"client"`
	Location   string  `json:"location"`
	Type       string  `json:"type"`
	Property   string  `json:"property"`
	PropertyID *string `json:"propertyId,omitempty"`
	AgentID    *string `json:"agentId,omitempty"`
	CreatedBy  *string `json:"createdBy,omitempty"`
	UpdatedBy  *string `json:"updatedBy,omitempty"`
}

type WhyChooseUsItem struct {
	ID          string  `json:"id"`
	Number      string  `json:"number"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Desc        string  `json:"desc"`
	CreatedBy   *string `json:"createdBy,omitempty"`
	UpdatedBy   *string `json:"updatedBy,omitempty"`
}

type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode string `json:"pincode"`
	Country string `json:"country"`
}

type Branch struct {
	City  string `json:"city"`
	Area  string `json:"area"`
	Phone string `json:"phone"`
}

type CompanyInfo struct {
	Name             string   `json:"name"`
	BrandName        string   `json:"brandName"`
	Tagline          string   `json:"tagline"`
	Headquarters     Address  `json:"headquarters"`
	Branches         []Branch `json:"branches"`
	Phone            string   `json:"phone"`
	PhoneRaw         string   `json:"phoneRaw"`
	WhatsApp         string   `json:"whatsapp"`
	Email            string   `json:"email"`
	Founded          string   `json:"founded"`
	EstablishedYear  int      `json:"establishedYear"`
	ReraBrokerNumber string   `json:"reraBrokerNumber"`
	Hours            string   `json:"hours"`
}
