package product

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNotFound is returned when a product cannot be found.
var ErrNotFound = errors.New("product not found")

// Store defines the interface for product persistence.
type Store interface {
	Create(p Product) (Product, error)
	List() ([]Product, error)
	Get(id string) (Product, error)
	Update(id string, p Product) (Product, error)
	Delete(id string) error
}

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu      sync.RWMutex
	items   map[string]Product
	counter int64
}

// NewMemoryStore returns a new empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[string]Product)}
}

// Create adds a new product and returns it with its assigned ID.
func (s *MemoryStore) Create(p Product) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	p.ID = fmt.Sprintf("%d", s.counter)
	s.items[p.ID] = p
	return p, nil
}

// List returns all products.
func (s *MemoryStore) List() ([]Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Product, 0, len(s.items))
	for _, p := range s.items {
		result = append(result, p)
	}
	return result, nil
}

// Get returns the product with the given id.
func (s *MemoryStore) Get(id string) (Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

// Update replaces the product with the given id.
func (s *MemoryStore) Update(id string, p Product) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return Product{}, ErrNotFound
	}
	p.ID = id
	s.items[id] = p
	return p, nil
}

// Delete removes the product with the given id.
func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
