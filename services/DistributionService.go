package services

import (
	"vkspam/models"
	"vkspam/repositories"
)

type DistributionService interface {
	Save(distribution models.Distribution) error
	Get(id int) (*models.Distribution, error)
	GetListByGroup(groupId int) (*[]models.Distribution, error)
	DeleteById(id int) error
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

func (s *distributionService) GetListByGroup(groupId int) (*[]models.Distribution, error) {
	return s.repo.GetListByGroup(groupId)
}

func (s *distributionService) DeleteById(id int) error {
	return s.repo.DeleteById(id)
}
