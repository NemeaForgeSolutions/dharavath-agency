package domain

import (
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type PropertyRepository interface {
	FindAll() []Property
	FindByFilter(filter PropertyFilter) []Property
	FindBySlug(slugOrID string) (*Property, error)
	FindFeatured(propType string, limit int) []Property
	FindSimilar(current *Property, limit int) []Property
	Count() int
	Create(prop *Property) error
	Update(prop *Property) error
	Delete(id string) error
	ToggleFeatured(id string) (bool, error)
}

type LeadRepository interface {
	Save(lead *Lead) error
	FindAll() []Lead
}

type CatalogRepository interface {
	FindAgents() []Agent
	FindAgentByID(id string) (*Agent, error)
	FindProjects() []Project
	FindLocations() []LocationInsight
	FindInsights() []InsightArticle
	FindTestimonials() []Testimonial
	FindWhyChooseUs() []WhyChooseUsItem
}
