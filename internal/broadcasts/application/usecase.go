package application

import (
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type usecase struct {
	repository       domain.BroadcastsRepository
	broadcastUsecase domain.BroadcastConnectionUsecase
	viewerUsecase    domain.ViewerConnectionUsecase
	config           *config.Config
	logger           domain.Logger
}

type BroadcastUsecaseParams struct {
	Repository       domain.BroadcastsRepository
	BroadcastUsecase domain.BroadcastConnectionUsecase
	ViewerUsecase    domain.ViewerConnectionUsecase
	Config           *config.Config
	Logger           domain.Logger
}

// Update the NewBroadcastUsecase function to accept the new struct
func NewBroadcastUsecase(params BroadcastUsecaseParams) domain.BroadcastsUsecases {
	return &usecase{
		repository:       params.Repository,
		broadcastUsecase: params.BroadcastUsecase,
		viewerUsecase:    params.ViewerUsecase,
		config:           params.Config,
		logger:           params.Logger,
	}
}
