package services

import (
	"vkspam/models"
	"vkspam/repositories"
)

type DistributionGroupService interface {
	GetList(userId int) ([]models.DistributionGroup, error)
}

type distributionGroupService struct {
	repo repositories.DistributionGroupRepository
}

func NewDistributionGroupService(repo repositories.DistributionGroupRepository) DistributionGroupService {
	return &distributionGroupService{repo: repo}
}

func (s *distributionGroupService) GetList(userId int) ([]models.DistributionGroup, error) {
	return s.repo.GetList(userId)
}
