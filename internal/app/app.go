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

	userlist, err := loader.LoadFile(cfg.Filepath)

	if err != nil {
		return err
	}

	validuser := service.SetValidUsers(userlist)

	invaliduser := service.SetInvalidUsers(userlist)

	fmt.Println(invaliduser)

	log.Printf("Loaded %d valid users", len(validuser))
	log.Printf("Skipped %d invalid users", len(invaliduser))

	userhandler := api.UserHandler(cfg.Filepath)
	cacheuserhandler := api.CacheUsersHandler(validuser)
	invaliduserhandler := api.InvalidUserHandler(invaliduser)

	http.HandleFunc("/users-nocaching", userhandler)
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/users", cacheuserhandler)
	http.HandleFunc("/users/invalid", invaliduserhandler)

	addr := fmt.Sprintf(":%d", cfg.Port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	return httpServer.ListenAndServe()
}
