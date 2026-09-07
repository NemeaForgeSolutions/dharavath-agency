package view

import (
	"dharavath-agency/internal/domain"
)

type PageData struct {
	Title        string
	Description  string
	ActivePage   string
	CanonicalURL string
	Robots       string
	OGType       string
	OGImage      string
	SiteURL      string
	Company      domain.CompanyInfo
	Data         interface{}
	IsHTMX       bool
}
