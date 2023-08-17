package main

import (
	"errors"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/peternagy652/k8s_sample_backend/handlers"
	"github.com/peternagy652/k8s_sample_backend/repository"
)

type environment struct {
	DBUser      string
	DBPassword  string
	DBName      string
	DBAddress   string
	HostAddress string
	Repository  string
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	env := getEnvironment()

	var repo repository.Repository

	if env.Repository == "inmemory" {
		repo = repository.NewInmemoryRepository()
	} else {
		if env.DBUser == "" || env.DBPassword == "" || env.DBName == "" {
			e.Logger.Fatal(errors.New("In case of a non inmemory repository the DB user, password and name has to be specified"))
		}

		var err error
		repo, err = repository.NewPostgreRepository(env.DBUser, env.DBPassword, env.DBName, env.DBAddress)
		if err != nil {
			e.Logger.Fatal(err)
		}
	}

	e.Use(bindRepository(repo))

	e.GET("/", handlers.Hello)
	e.GET("/api/v1/hello", handlers.Hello)
	e.POST("/api/v1/person", handlers.AddPersonHandler)
	e.GET("/api/v1/person/:id", handlers.GetPersonByIDHandler)
	e.GET("/api/v1/persons", handlers.GetPersonsHandler)
	e.POST("/api/v1/generate", handlers.GeneratePersons)
	e.DELETE("/api/v1/persons", handlers.ClearPersons)

	e.Logger.Fatal(e.Start(env.HostAddress))
}

func getEnvironment() environment {
	repo := os.Getenv("REPOSITORY")
	if repo == "" {
		repo = "inmemory"
	}

	hostaddress := os.Getenv("HOST_ADDRESS")
	if hostaddress == "" {
		hostaddress = "0.0.0.0:7992"
	}

	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbaddress := os.Getenv("DB_ADDRESS")

	return environment{
		Repository:  repo,
		HostAddress: hostaddress,
		DBUser:      dbuser,
		DBPassword:  dbpassword,
		DBName:      dbname,
		DBAddress:   dbaddress,
	}
}

func bindRepository(repo repository.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("repository", repo)
			return next(c)
		}
	}
}
