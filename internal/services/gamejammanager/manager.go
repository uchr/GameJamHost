package gamejammanager

import "GameJamPlatform/internal/storages"

type gameJamManager struct {
	repo storages.Repo
}

func NewService(repo storages.Repo) *gameJamManager {
	return &gameJamManager{repo: repo}
}

var _ GameJamManager = (*gameJamManager)(nil)
