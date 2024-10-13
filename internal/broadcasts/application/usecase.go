package application

import (
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type usecase struct {
	repository domain.BroadcastsRepository
	config     *config.Config
	logger     domain.Logger
}

func NewBroadcastUsecase(repository domain.BroadcastsRepository, config *config.Config, logger domain.Logger) domain.BroadcastsUsecases {
	return &usecase{
		config:     config,
		repository: repository,
		logger:     logger,
	}
}
