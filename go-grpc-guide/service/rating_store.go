package service

import "sync"

// RatingStore is an interface to store laptop ratings
type RatingStore interface {
	// Add new laptop score to the store and returns its rating
	Add(laptopID string, score float64) (*Rating, error)
}

// Rating contains the rating informatio of a laptop
type Rating struct {
	Count uint32
	Sum   float64
}

// ImMemoryRatingStore stores laptop ratings in memory
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

// NewInMemoryRatingStore returns a new InMemoryRatingStore
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

// Add new laptop score to the store and returns its rating
func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}

	store.rating[laptopID] = rating
	return rating, nil
}
