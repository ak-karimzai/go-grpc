package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/jinzhu/copier"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Fetch(id string) (*pb.Laptop, error)
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: map[string]*pb.Laptop{},
	}
}

func (store *InMemoryLaptopStore) Save(
	laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %v", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Fetch(
	id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot copy laptop data: %v", err)
	}

	return other, nil
}

// type DBLaptopStore struct {
// }
