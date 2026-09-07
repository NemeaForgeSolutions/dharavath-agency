package memory

import (
	"fmt"
	"sync"
	"time"

	"dharavath-agency/internal/domain"
)

type LeadRepo struct {
	mu    sync.RWMutex
	leads []domain.Lead
}

func NewLeadRepo() *LeadRepo {
	return &LeadRepo{
		leads: make([]domain.Lead, 0),
	}
}

func (r *LeadRepo) Save(lead *domain.Lead) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lead.ID = fmt.Sprintf("LEAD-%d", time.Now().UnixNano())
	lead.CreatedAt = time.Now()
	r.leads = append(r.leads, *lead)
	return nil
}

func (r *LeadRepo) FindAll() []domain.Lead {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Lead, len(r.leads))
	copy(res, r.leads)
	return res
}
