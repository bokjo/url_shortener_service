package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	apiHandler "github.com/bokjo/url_shortener_service/api"
	r "github.com/bokjo/url_shortener_service/repository/redis"
	"github.com/bokjo/url_shortener_service/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	fmt.Println("Hello, World!")

	// Init data source, example Redis repo instance
	// redisURL := "redis://localhost:6379"
	redisURL := "redis://redisurlshortener"
	repo, err := r.NewRedisRepository(redisURL)

	if err != nil {
		log.Fatal(err)
	}

	// Init shortener service instance ( in needs the data source passed to it as dependency)
	service := shortener.NewShortenerService(repo)

	// Init our API handlers
	api := apiHandler.NewAPIHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", api.Get)
	router.Post("/", api.Post)

	// Start out API/http server
	errs := make(chan error, 2)

	go func() {
		fmt.Println("Shortener API listening on port 1234")
		errs <- http.ListenAndServe(":1234", router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Shortener service terminated %s", <-errs)
}
