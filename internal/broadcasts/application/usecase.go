package application

import "github.com/yavurb/rill/internal/broadcasts/domain"

type usecase struct {
	repository domain.BroadcastsRepository
}

func NewBroadcastUsecase(repository domain.BroadcastsRepository) domain.BroadcastsUsecases {
	return &usecase{repository}
}
