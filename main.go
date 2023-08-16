package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/peternagy652/k8s_sample_backend/handlers"
	"github.com/peternagy652/k8s_sample_backend/repository"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Info("Some uuid: ", uuid.New().String())

	// repo := repository.NewInmemoryRepository()

	repo, err := repository.NewPostgreRepository("peti", "pass", "k8s", "")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(bindRepository(repo))

	e.GET("/", handlers.Hello)
	e.POST("/person", handlers.AddPersonHandler)
	e.GET("/person/:id", handlers.GetPersonByIDHandler)
	e.GET("/persons", handlers.GetPersonsHandler)
	e.POST("/generate", handlers.GeneratePersons)

	e.Logger.Fatal(e.Start(":7992"))
}

func bindRepository(repo repository.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("repository", repo)
			return next(c)
		}
	}
}
