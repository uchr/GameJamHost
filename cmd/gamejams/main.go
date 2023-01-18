package main

import (
	"context"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/services/gamejammanager"
	"GameJamPlatform/internal/services/sessionprovider"
	"GameJamPlatform/internal/services/usersmanager"
	"GameJamPlatform/internal/storages"
	"GameJamPlatform/internal/web/servers"
	"GameJamPlatform/internal/web/templatemanager"
)

func main() {
	log.Init(log.DebugLevel, "logs")

	ctx := context.Background()

	storageConfig, err := storages.NewConfig()
	if err != nil {
		log.Panic(err, "failed to create storage config")
	}

	repo, err := storages.NewStorage(ctx, storageConfig)
	if err != nil {
		log.Panic(err, "failed to create storage")
	}

	tm, err := templatemanager.NewManager("web/template")
	if err != nil {
		log.Panic(err, "failed to create templates")
	}

	serverConfig, err := servers.NewConfig()
	if err != nil {
		log.Panic(err, "failed to create server config")
	}

	u := usersmanager.NewUsers(repo)

	sp := sessionprovider.NewProvider(repo)

	service := gamejammanager.NewService(repo)
	server := servers.NewServer(service, tm, u, sp, serverConfig)
	err = server.Run()
	if err != nil {
		log.Panic(err, "failed to start server")
	}
}
