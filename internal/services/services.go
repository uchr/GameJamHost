package services

import "GameJamPlatform/internal/storages"

type Service struct {
	repo storages.Repo
}

func NewService(repo storages.Repo) *Service {
	return &Service{repo: repo}
}
