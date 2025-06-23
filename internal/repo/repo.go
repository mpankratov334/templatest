package repo

import (
	"Templatest/internal/config"
	"Templatest/internal/dto"
	"context"
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"time"
)

type Repository interface {
	Post(ctx context.Context, obj DataObject) (int, error)
	GetAll(ctx context.Context) (*[]DataObject, error)
	Get(ctx context.Context, id string) (*DataObject, error)
	Put(ctx context.Context, id string) (*DataObject, error)
	Delete(ctx context.Context, id string) (*DataObject, error)
}

type repository struct {
	_slice []DataObject
	cap    int
	mutex  *sync.Mutex
}

// NewRepository returns array-based in-memory repository
func NewRepository(cfg config.InMemory) (Repository, error) {
	slice := make([]DataObject, 0, cfg.Capacity)
	return &repository{
		_slice: slice,
		cap:    cfg.Capacity,
		mutex:  &sync.Mutex{},
	}, nil
}

func (r *repository) Post(ctx context.Context, obj DataObject) (int, error) {
	if len(r._slice) == r.cap {
		return 0, errors.New("unavailable for recording, memory is full")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	obj.ID = strconv.Itoa(len(r._slice) + 1)
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = time.Now()

	r._slice = append(r._slice, obj)
	return len(r._slice), nil
}

// GetAll return pointer to the original slice
// Be careful to change it
func (r *repository) GetAll(ctx context.Context) (*[]DataObject, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return &r._slice, nil
}

// Get returns pointer to in-time original object
// Consider this when change it
func (r *repository) Get(ctx context.Context, id string) (*DataObject, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {

		return nil, errors.Wrap(err, "error converting id to int")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if idInt > len(r._slice) || idInt < 1 {
		return nil, dto.ErrInvalidID
	}
	idx := idInt - 1

	return &r._slice[idx], nil
}

func (r *repository) Put(ctx context.Context, id string) (*DataObject, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrap(err, "error converting id to int")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if idInt > len(r._slice) || idInt < 1 {
		return nil, dto.ErrInvalidID
	}
	idx := idInt - 1

	r._slice[idx].UpdatedAt = time.Now()

	return &r._slice[idx], nil
}

func (r *repository) Delete(ctx context.Context, id string) (*DataObject, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrap(err, "error converting id to int")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if idInt > len(r._slice) || idInt < 1 {
		return nil, dto.ErrInvalidID
	}
	idx := idInt - 1
	deletedObj := r._slice[idx]

	r._slice = append(r._slice[:idx], r._slice[idx+1:]...)
	for i := idx; i < len(r._slice); i++ {
		r._slice[i].ID = strconv.Itoa(i + 1)
	}

	return &deletedObj, nil
}
