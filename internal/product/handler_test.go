package product_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mime-rona/irp-app-from-template/internal/product"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// mockStore is a controllable Store for unit tests.
type mockStore struct {
	createFn func(product.Product) (product.Product, error)
	listFn   func() ([]product.Product, error)
	getFn    func(string) (product.Product, error)
	updateFn func(string, product.Product) (product.Product, error)
	deleteFn func(string) error
}

func (m *mockStore) Create(p product.Product) (product.Product, error) { return m.createFn(p) }
func (m *mockStore) List() ([]product.Product, error)                  { return m.listFn() }
func (m *mockStore) Get(id string) (product.Product, error)            { return m.getFn(id) }

func (m *mockStore) Update(id string, p product.Product) (product.Product, error) {
	return m.updateFn(id, p)
}

func (m *mockStore) Delete(id string) error { return m.deleteFn(id) }

func newTestRouter(store product.Store) *gin.Engine {
	r := gin.New()
	h := product.NewHandler(store)
	products := r.Group("/products")
	products.POST("", h.Create)
	products.GET("", h.List)
	products.GET("/:id", h.Get)
	products.PUT("/:id", h.Update)
	products.DELETE("/:id", h.Delete)
	return r
}

// --- Create ---

func TestHandler_Create_Success(t *testing.T) {
	store := &mockStore{
		createFn: func(p product.Product) (product.Product, error) {
			p.ID = "1"
			return p, nil
		},
	}
	router := newTestRouter(store)

	body := strings.NewReader(`{"name":"Widget","price":9.99}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/products", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"1"`)
}

func TestHandler_Create_BadRequest(t *testing.T) {
	router := newTestRouter(&mockStore{})

	body := strings.NewReader(`not-json`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/products", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_StoreError(t *testing.T) {
	store := &mockStore{
		createFn: func(_ product.Product) (product.Product, error) {
			return product.Product{}, errors.New("store failure")
		},
	}
	router := newTestRouter(store)

	body := strings.NewReader(`{"name":"X","price":1}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/products", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- List ---

func TestHandler_List_Success(t *testing.T) {
	store := &mockStore{
		listFn: func() ([]product.Product, error) {
			return []product.Product{{ID: "1", Name: "A", Price: 1}}, nil
		},
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/products", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"1"`)
}

func TestHandler_List_StoreError(t *testing.T) {
	store := &mockStore{
		listFn: func() ([]product.Product, error) {
			return nil, errors.New("store failure")
		},
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/products", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Get ---

func TestHandler_Get_Success(t *testing.T) {
	store := &mockStore{
		getFn: func(id string) (product.Product, error) {
			return product.Product{ID: id, Name: "Widget", Price: 9.99}, nil
		},
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/products/1", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	store := &mockStore{
		getFn: func(_ string) (product.Product, error) {
			return product.Product{}, product.ErrNotFound
		},
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/products/99", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Get_StoreError(t *testing.T) {
	store := &mockStore{
		getFn: func(_ string) (product.Product, error) {
			return product.Product{}, errors.New("store failure")
		},
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/products/1", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Update ---

func TestHandler_Update_Success(t *testing.T) {
	store := &mockStore{
		updateFn: func(id string, p product.Product) (product.Product, error) {
			p.ID = id
			return p, nil
		},
	}
	router := newTestRouter(store)

	body := strings.NewReader(`{"name":"Updated","price":5}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, "/products/1", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_BadRequest(t *testing.T) {
	router := newTestRouter(&mockStore{})

	body := strings.NewReader(`bad`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, "/products/1", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_NotFound(t *testing.T) {
	store := &mockStore{
		updateFn: func(_ string, _ product.Product) (product.Product, error) {
			return product.Product{}, product.ErrNotFound
		},
	}
	router := newTestRouter(store)

	body := strings.NewReader(`{"name":"X","price":1}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, "/products/99", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Update_StoreError(t *testing.T) {
	store := &mockStore{
		updateFn: func(_ string, _ product.Product) (product.Product, error) {
			return product.Product{}, errors.New("store failure")
		},
	}
	router := newTestRouter(store)

	body := strings.NewReader(`{"name":"X","price":1}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, "/products/1", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Delete ---

func TestHandler_Delete_Success(t *testing.T) {
	store := &mockStore{
		deleteFn: func(_ string) error { return nil },
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, "/products/1", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	store := &mockStore{
		deleteFn: func(_ string) error { return product.ErrNotFound },
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, "/products/99", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Delete_StoreError(t *testing.T) {
	store := &mockStore{
		deleteFn: func(_ string) error { return errors.New("store failure") },
	}
	router := newTestRouter(store)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, "/products/1", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
