package application

import (
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type usecase struct {
	config     *config.Config
	repository domain.BroadcastsRepository
}

func NewBroadcastUsecase(repository domain.BroadcastsRepository, config *config.Config) domain.BroadcastsUsecases {
	return &usecase{
		config:     config,
		repository: repository,
	}
}
