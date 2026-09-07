package domain

import (
	"errors"
	"strings"
	"time"
)

type Lead struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	PropertyID    string    `json:"propertyId"`
	PropertyTitle string    `json:"propertyTitle"`
	Message       string    `json:"message"`
	Type          string    `json:"type"` // "site_visit", "valuation_request", "general_contact"
	PreferredDate string    `json:"preferredDate"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (l *Lead) Validate() error {
	if strings.TrimSpace(l.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(l.Phone) == "" {
		return errors.New("phone number is required")
	}
	return nil
}
