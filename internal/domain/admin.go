package domain

type AdminStats struct {
	TotalProperties            int    `json:"totalProperties"`
	TotalLeads                 int    `json:"totalLeads"`
	PendingSellRequests        int    `json:"pendingSellRequests"`
	FeaturedCount              int    `json:"featuredCount"`
	TotalPortfolioValue        int64  `json:"totalPortfolioValue"`
	TotalPortfolioValueDisplay string `json:"totalPortfolioValueDisplay"`
}
