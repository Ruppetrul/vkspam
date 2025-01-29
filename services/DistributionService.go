package services

import (
	"vkspam/repositories"
)

type DistributionService interface {
}

type distributionService struct {
	repo repositories.DistributionRepository
}

func NewDistributionService(repo repositories.DistributionRepository) DistributionService {
	return &distributionService{repo: repo}
}
