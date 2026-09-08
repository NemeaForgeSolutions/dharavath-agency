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
	UpdateStatus(id string, status string) error
	FindByID(id string) (*Lead, error)
}

type CatalogRepository interface {
	FindAgents() []Agent
	FindAgentByID(id string) (*Agent, error)
	CreateAgent(agent *Agent) error
	UpdateAgent(agent *Agent) error
	DeleteAgent(id string) error
	FindProjects() []Project
	FindLocations() []LocationInsight
	FindInsights() []InsightArticle
	FindTestimonials() []Testimonial
	FindWhyChooseUs() []WhyChooseUsItem
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
	UpdateLastLogin(id string) error
	Count() int
	FindAll() []User
}

