package service

import (
	"dharavath-agency/internal/domain"
)

type AdminService struct {
	propertyRepo domain.PropertyRepository
	leadRepo     domain.LeadRepository
}

func NewAdminService(pRepo domain.PropertyRepository, lRepo domain.LeadRepository) *AdminService {
	return &AdminService{
		propertyRepo: pRepo,
		leadRepo:     lRepo,
	}
}

func (s *AdminService) GetStats() domain.AdminStats {
	props := s.propertyRepo.FindAll()
	leads := s.leadRepo.FindAll()

	var totalValue int64
	featuredCount := 0

	for _, p := range props {
		totalValue += p.Price
		if p.Featured {
			featuredCount++
		}
	}

	return domain.AdminStats{
		TotalProperties:            len(props),
		TotalLeads:                 len(leads),
		FeaturedCount:              featuredCount,
		TotalPortfolioValue:        totalValue,
		TotalPortfolioValueDisplay: FormatIndianCurrency(totalValue),
	}
}

func (s *AdminService) GetRecentLeads() []domain.Lead {
	leads := s.leadRepo.FindAll()
	reversed := make([]domain.Lead, len(leads))
	for i, l := range leads {
		reversed[len(leads)-1-i] = l
	}
	if len(reversed) > 10 {
		return reversed[:10]
	}
	return reversed
}
