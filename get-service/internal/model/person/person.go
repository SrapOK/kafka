package person

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PersonRepo struct {
	Db *redis.Client
}

type Person struct {
	Id      string `json:"id" redis:"id"`
	Name    string `json:"name" redis:"name"`
	Surname string `json:"surname" redis:"surname"`
}

func (p *PersonRepo) GetPerson(c *context.Context, key string) (*Person, error) {

	val, err := p.Db.Get(*c, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %s", err.Error())
	}

	var ps Person

	err = json.Unmarshal([]byte(val), &ps)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

func (p *PersonRepo) GetRandomPerson(c *context.Context) (*Person, error) {
	key, err := p.Db.RandomKey(*c).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get random key: %s", err.Error())
	}

	ps, err := p.GetPerson(c, key)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (p *PersonRepo) GetAllPersons(c *context.Context) ([]Person, error) {
	var list []Person

	iter := p.Db.Scan(*c, 0, "*", 0).Iterator()
	for iter.Next(*c) {
		ps, err := p.GetPerson(c, iter.Val())
		if err != nil {
			return []Person{}, err
		}
		list = append(list, *ps)
	}

	if err := iter.Err(); err != nil {
		return []Person{}, err
	}
	return list, nil
}

func (p *PersonRepo) SavePerson(c *context.Context, person *PersonDTO) error {
	key := uuid.NewString()

	ps := Person{Id: key, Name: person.Name, Surname: person.Surname}
	val, err := json.Marshal(&ps)
	if err != nil {
		return fmt.Errorf("failed to save person: %s", err.Error())
	}

	if err := p.Db.Set(*c, key, val, 0).Err(); err != nil {
		return fmt.Errorf("failed to save person: %s", err.Error())
	}

	return nil
}
