package service_test

import (
	"testing"

	"dharavath-agency/internal/domain"
	"dharavath-agency/internal/repository/memory"
	"dharavath-agency/internal/service"
)

func TestPropertyService(t *testing.T) {
	repo, err := memory.NewPropertyRepo()
	if err != nil {
		t.Fatalf("Failed to initialize memory repo: %v", err)
	}

	svc := service.NewPropertyService(repo, repo)

	if svc.TotalCount() == 0 {
		t.Errorf("Expected properties > 0, got %d", svc.TotalCount())
	}

	// Filter by city
	hydProps := svc.ListProperties(domain.PropertyFilter{City: "Hyderabad"})
	if len(hydProps) == 0 {
		t.Errorf("Expected Hyderabad properties, got none")
	}
	for _, p := range hydProps {
		if p.City != "Hyderabad" {
			t.Errorf("Expected city Hyderabad, got %s", p.City)
		}
	}

	// Detail lookup
	prop, agent, similar, err := svc.GetDetail("the-aurora-residences-4bhk-financial-district")
	if err != nil {
		t.Fatalf("Failed to get detail: %v", err)
	}
	if prop.Bedrooms != 4 {
		t.Errorf("Expected 4 BHK, got %d", prop.Bedrooms)
	}
	if agent == nil {
		t.Errorf("Expected assigned agent to be loaded")
	}
	if len(similar) == 0 {
		t.Errorf("Expected similar properties to be found")
	}
}

func TestLeadService(t *testing.T) {
	repo := memory.NewLeadRepo()
	svc := service.NewLeadService(repo)

	// Validation error
	_, err := svc.SubmitLead(domain.Lead{Name: "", Phone: ""})
	if err == nil {
		t.Errorf("Expected validation error for empty lead")
	}

	// Valid submission
	res, err := svc.SubmitLead(domain.Lead{
		Name:          "Anand Mahindra",
		Phone:         "+91 98200 12345",
		Email:         "anand@example.com",
		PropertyID:    "PROP-MUM-005",
		PropertyTitle: "One Bandra Coastal Penthouse",
		Type:          "site_visit",
	})
	if err != nil {
		t.Fatalf("Failed to submit valid lead: %v", err)
	}
	if res.Lead.ID == "" {
		t.Errorf("Expected generated lead ID")
	}
	if res.WhatsAppURL == "" {
		t.Errorf("Expected WhatsApp URL to be generated")
	}

	savedLeads := repo.FindAll()
	if len(savedLeads) != 1 {
		t.Errorf("Expected 1 saved lead in repository, got %d", len(savedLeads))
	}
}
