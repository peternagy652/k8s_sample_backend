package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
	"github.com/peternagy652/k8s_sample_backend/model"
	"github.com/peternagy652/k8s_sample_backend/util"
)

type PostgresRepository struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	DBName   string `json:"DBName"`
	Address  string `json:"DBAddress"`
	DB       *pg.DB `json:"DB"`
}

func NewPostgreRepository(user, password, dbname, address string) (Repository, error) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     address,
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}

	repo := PostgresRepository{
		User:     user,
		Password: password,
		DBName:   dbname,
		Address:  address,
		DB:       db,
	}

	return &repo, nil
}

func createSchema(db *pg.DB) error {
	err := db.Model((*model.Person)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) Close() {
	p.DB.Close()
}

func (p *PostgresRepository) AddPerson(person model.Person) (string, error) {
	person.Id = uuid.New().String()
	_, err := p.DB.Model(&person).Insert()
	if err != nil {
		return "", err
	}
	return person.Id, nil
}

func (p *PostgresRepository) GetPersonByID(id string) (model.Person, error) {
	person := &model.Person{Id: id}

	err := p.DB.Model(person).WherePK().Select()
	if err != nil {
		return model.Person{}, err
	}

	return *person, nil
}

func (p *PostgresRepository) GetPersons() ([]model.Person, error) {
	var persons []model.Person

	err := p.DB.Model(&persons).Select()
	if err != nil {
		return make([]model.Person, 0), err
	}

	return persons, nil
}

func (p *PostgresRepository) ModifyPerson(person model.Person) (model.Person, error) {
	return model.Person{}, &util.NotImplementedError{}
}

func (p *PostgresRepository) DeletePersonByID(id string) error {
	return &util.NotImplementedError{}
}
