package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/peternagy652/k8s_sample_backend/model"
	"github.com/peternagy652/k8s_sample_backend/util"
)

type InmemoryRepository struct {
	personCache []model.Person
	personsLock *sync.RWMutex
}

func NewInmemoryRepository() Repository {
	repo := InmemoryRepository{
		personCache: make([]model.Person, 0),
		personsLock: &sync.RWMutex{},
	}

	return &repo
}

func (i *InmemoryRepository) AddPerson(person model.Person) (string, error) {
	i.personsLock.Lock()
	defer i.personsLock.Unlock()

	person.Id = uuid.New().String()

	i.personCache = append(i.personCache, person)

	return person.Id, nil
}

func (i *InmemoryRepository) GetPersonByID(id string) (model.Person, error) {
	i.personsLock.RLock()
	defer i.personsLock.RUnlock()

	for _, person := range i.personCache {
		if person.Id == id {
			return person, nil
		}
	}

	return model.Person{}, &util.NotFoundError{}
}

func (i *InmemoryRepository) GetPersons() ([]model.Person, error) {
	i.personsLock.RLock()
	defer i.personsLock.RUnlock()

	return i.personCache, nil
}

func (i *InmemoryRepository) ModifyPerson(person model.Person) (model.Person, error) {
	i.personsLock.Lock()
	defer i.personsLock.Unlock()

	for j, p := range i.personCache {
		if p.Id == person.Id {
			i.personCache[j] = person
			return person, nil
		}
	}

	return model.Person{}, &util.NotFoundError{}
}

func (i *InmemoryRepository) DeletePersonByID(id string) error {
	i.personsLock.Lock()
	defer i.personsLock.Unlock()

	for j, p := range i.personCache {
		if p.Id == id {
			// Don't care about the order
			i.personCache[j] = i.personCache[len(i.personCache)-1]
			i.personCache = i.personCache[:len(i.personCache)-1]
			return nil
		}
	}

	return &util.NotFoundError{}
}

func (i *InmemoryRepository) Clear() error {
	i.personsLock.Lock()
	defer i.personsLock.Unlock()

	i.personCache = make([]model.Person, 0)

	return nil
}
