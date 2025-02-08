package services

import (
	"vkspam/models"
	"vkspam/repositories"
)

type DistributionService interface {
	Save(distribution models.Distribution) error
	Get(id int) (*models.Distribution, error)
}

type distributionService struct {
	repo repositories.DistributionRepository
}

func NewDistributionService(repo repositories.DistributionRepository) DistributionService {
	return &distributionService{repo: repo}
}

func (s *distributionService) Save(distribution models.Distribution) error {
	return s.repo.Save(distribution)
}

func (s *distributionService) Get(id int) (*models.Distribution, error) {
	return s.repo.Get(id)
}
