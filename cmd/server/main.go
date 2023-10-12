package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/gonozov0/go-musthave-devops/internal/server"
	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
	inmemory "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	postgres "github.com/gonozov0/go-musthave-devops/internal/server/repository/postgres"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %s", err.Error())
	}

	var repo repository.Repository
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	if cfg.DatabaseDSN != "" {
		repo, err = postgres.NewPGRepository(cfg.DatabaseDSN)
		if err != nil {
			log.Fatalf("Could not init postgres repository: %s", err.Error())
		}
	} else if cfg.FileStoragePath != "" {
		wg.Add(1)
		repo, err = inmemory.NewInMemoryRepositoryWithFileStorage(
			ctx,
			wg,
			cfg.FileStoragePath,
			cfg.StoreInterval,
			cfg.RestoreFlag,
		)
		if err != nil {
			log.Fatalf("Could not init in memory repository: %s", err.Error())
		}
	} else {
		repo = inmemory.NewInMemoryRepository()
	}

	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	router := application.NewRouter(repo)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		log.Infof("Starting server on port %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-stopChan:
		log.Info("Received signal to stop. Shutting down...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorf("Server shutdown failed:%+v", err)
		}
		cancel()
		wg.Wait()
	case err := <-errChan:
		log.Fatalf("Server error: %s", err.Error())
	}
}
