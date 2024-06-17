package service

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/hwebz/go-grpc-guide/pb"
	"github.com/jinzhu/copier"
	"log"
	"sync"
)

var ErrAlreadyExists = errors.New("Record already exists")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

type InMemoryLaptopStore struct {
	mutex sync.Mutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}

// Search laptops with filter, returns one by one via found function
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for _, laptop := range store.data {
		// FOR TESTING: heavy processing
		//time.Sleep(time.Second)
		log.Print("Checking laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled {
			log.Println("Context is cancelled")
			return errors.New("Context is cancelled")
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	log.Println("isQualified: ", filter, laptop)
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3 // value * 8
	case pb.Memory_KILOBYTE:
		return value << 13 // value * 8 * 1024
	case pb.Memory_MEGABYTE:
		return value << 23 // value * 8 * 1024 * 1024
	case pb.Memory_GIGABYTE:
		return value << 33 // value * 8 * 1024 * 1024 * 1024
	case pb.Memory_TERABYTE:
		return value << 43 // value * 8 * 1024 * 1024 * 1024 * 1024
	default:
		return value
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("Cannot copy laptop data: %w", err)
	}

	return other, nil
}
