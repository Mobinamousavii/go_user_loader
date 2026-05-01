package app

import (
	"fmt"
	"goproject/internal/api"
	"goproject/internal/config"
	"goproject/internal/loader"
	"goproject/internal/service"
	"log"
	"net/http"
)

func Run(args []string) error {
	cfg, err := config.Parse(args)

	if err != nil {
		return err
	}

	svc := service.NewService()
	log.Printf("Starting processing with %d workers", cfg.Workers)

	records, errCh := loader.LoadFile(cfg.Filepath)

	if err := svc.ProcessStream(records, cfg.Workers); err != nil {
		return err
	}

	if err := <-errCh; err != nil {
		return err
	}

	log.Printf("Loaded %d valid users", svc.ValidCount())
	log.Printf("Skipped %d invalid users", svc.InvalidCount())

	srv := api.New(svc)
	addr := fmt.Sprintf(":%d", cfg.Port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv.Handler(),
	}

	return httpServer.ListenAndServe()
}
