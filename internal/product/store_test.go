package product_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mime-rona/irp-app-from-template/internal/product"
)

func TestMemoryStore_CreateAndGet(t *testing.T) {
	s := product.NewMemoryStore()

	p, err := s.Create(product.Product{Name: "Widget", Price: 9.99})
	require.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Widget", p.Name)
	assert.Equal(t, 9.99, p.Price)

	got, err := s.Get(p.ID)
	require.NoError(t, err)
	assert.Equal(t, p, got)
}

func TestMemoryStore_Get_NotFound(t *testing.T) {
	s := product.NewMemoryStore()

	_, err := s.Get("missing")
	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestMemoryStore_List(t *testing.T) {
	s := product.NewMemoryStore()

	_, err := s.Create(product.Product{Name: "A", Price: 1.0})
	require.NoError(t, err)
	_, err = s.Create(product.Product{Name: "B", Price: 2.0})
	require.NoError(t, err)

	all, err := s.List()
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestMemoryStore_Update(t *testing.T) {
	s := product.NewMemoryStore()

	p, err := s.Create(product.Product{Name: "Old", Price: 1.0})
	require.NoError(t, err)

	updated, err := s.Update(p.ID, product.Product{Name: "New", Price: 2.0})
	require.NoError(t, err)
	assert.Equal(t, p.ID, updated.ID)
	assert.Equal(t, "New", updated.Name)
}

func TestMemoryStore_Update_NotFound(t *testing.T) {
	s := product.NewMemoryStore()

	_, err := s.Update("missing", product.Product{Name: "X"})
	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestMemoryStore_Delete(t *testing.T) {
	s := product.NewMemoryStore()

	p, err := s.Create(product.Product{Name: "ToDelete", Price: 0.5})
	require.NoError(t, err)

	err = s.Delete(p.ID)
	require.NoError(t, err)

	_, err = s.Get(p.ID)
	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestMemoryStore_Delete_NotFound(t *testing.T) {
	s := product.NewMemoryStore()

	err := s.Delete("missing")
	assert.ErrorIs(t, err, product.ErrNotFound)
}
