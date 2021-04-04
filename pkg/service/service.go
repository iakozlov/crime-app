package service

import "github.com/iakozlov/crime-app/pkg/repository"

type Authorization interface {
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
