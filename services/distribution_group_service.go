package services

import (
	"vkspam/models"
	"vkspam/repositories"
)

type DistributionGroupService interface {
	GetList(userId int) ([]models.DistributionGroup, error)
	Save(group models.DistributionGroup) (int, error)
	Delete(id int) error
	Get(id int) (*models.DistributionGroup, error)
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

func (s *distributionGroupService) Save(group models.DistributionGroup) (int, error) {
	return s.repo.Save(group)
}

func (s *distributionGroupService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *distributionGroupService) Get(id int) (*models.DistributionGroup, error) {
	return s.repo.Get(id)
}
