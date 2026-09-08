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
	Type          string    `json:"type"`   // "site_visit", "valuation_request", "sell_request", "general_contact"
	Status        string    `json:"status"` // "pending", "approved", "rejected"
	PreferredDate string    `json:"preferredDate"`
	AgentID       *string   `json:"agentId,omitempty"`
	UserID        *string   `json:"userId,omitempty"`
	CreatedBy     *string   `json:"createdBy,omitempty"`
	UpdatedBy     *string   `json:"updatedBy,omitempty"`
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
