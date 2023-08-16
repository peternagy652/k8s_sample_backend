package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/peternagy652/k8s_sample_backend/model"
	"github.com/peternagy652/k8s_sample_backend/repository"
)

type AddPersonResult struct {
	ID string `json:"ID"`
}

func AddPersonHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	repo := c.Get("repository").(repository.Repository)

	p := model.Person{}
	err := json.NewDecoder(c.Request().Body).Decode(&p)
	if err != nil {
		c.Logger().Error("Failed to read request body: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uuid, err := repo.AddPerson(p)
	if err != nil {
		c.Logger().Error("Failed to add person to repository: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	r := AddPersonResult{ID: uuid}

	return c.JSON(http.StatusOK, r)
}

type GetPersonResult struct {
	Person model.Person `json:"Person"`
}

func GetPersonByIDHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param("id")
	if id == "" {
		c.Logger().Error("Failed to get id from request")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	repo := c.Get("repository").(repository.Repository)

	person, err := repo.GetPersonByID(id)
	if err != nil {
		c.Logger().Error("Failed to get person from repository: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	r := GetPersonResult{Person: person}

	return c.JSON(http.StatusOK, r)
}

type GetPersonsResult struct {
	Persons []model.Person `json:"Persons"`
}

func GetPersonsHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	repo := c.Get("repository").(repository.Repository)

	persons, err := repo.GetPersons()
	if err != nil {
		c.Logger().Error("Failed to get data from repository: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	r := GetPersonsResult{Persons: persons}

	return c.JSON(http.StatusOK, r)
}

type GeneratePersonsResult struct {
	IDs []string `json:"IDs"`
}

func GeneratePersons(c echo.Context) error {
	defer c.Request().Body.Close()

	repo := c.Get("repository").(repository.Repository)

	stringCount := c.QueryParam("count")
	count, err := strconv.ParseInt(stringCount, 10, 0)
	if err != nil {
		c.Logger().Error("Failed to get count: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	c.Logger().Debug("Going to generate %s new Persons.", stringCount)

	r := GeneratePersonsResult{
		IDs: make([]string, 0),
	}

	var i int64
	for i = 0; i < count; i++ {
		p := generatePerson(i)
		id, err := repo.AddPerson(p)
		if err != nil {
			c.Logger().Error("Failed to generate new person: %s", err)
			continue
		}

		r.IDs = append(r.IDs, id)
	}

	return c.JSON(http.StatusOK, r)
}

func generatePerson(index int64) model.Person {
	return model.Person{
		Id:   uuid.New().String(),
		Name: fmt.Sprintf("Person %d", index),
		Age:  int32(rand.Intn(70) + 20),
	}
}
