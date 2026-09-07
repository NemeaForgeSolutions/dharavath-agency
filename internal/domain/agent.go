package domain

type Agent struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	City           string   `json:"city"`
	AreasServed    []string `json:"areasServed"`
	Specialization string   `json:"specialization"`
	Experience     string   `json:"experience"`
	Languages      []string `json:"languages"`
	ReraID         string   `json:"reraId"`
	Photo          string   `json:"photo"`
	Phone          string   `json:"phone"`
	WhatsApp       string   `json:"whatsapp"`
	Email          string   `json:"email"`
	ActiveListings int      `json:"activeListings"`
	Bio            string   `json:"bio"`
}
