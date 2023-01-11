package main

import (
	"context"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/servers"
	"GameJamPlatform/internal/services"
	"GameJamPlatform/internal/storages"
	"GameJamPlatform/internal/templates"
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

	tmpl, err := templates.NewTemplates("web/template")
	if err != nil {
		log.Panic(err, "failed to create templates")
	}

	serverConfig, err := servers.NewConfig()
	if err != nil {
		log.Panic(err, "failed to create server config")
	}

	service := services.NewService(repo)
	server := servers.NewServer(service, tmpl, serverConfig)
	err = server.Run()
	if err != nil {
		log.Panic(err, "failed to start server")
	}
}
