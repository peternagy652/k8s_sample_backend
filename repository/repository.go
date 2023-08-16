package repository

import (
	"github.com/peternagy652/k8s_sample_backend/model"
)

type Repository interface {
	AddPerson(person model.Person) (string, error)
	GetPersonByID(id string) (model.Person, error)
	GetPersons() ([]model.Person, error)
	ModifyPerson(person model.Person) (model.Person, error)
	DeletePersonByID(id string) error
}
