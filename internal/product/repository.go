package product

import (
	"errors"
	"go_web_ejercicio_dominios/internal/domain"
	"go_web_ejercicio_dominios/pkg/store"
)

type Repository interface {
	GetAll() ([]*domain.Product, error)
	GetById(id int) (p *domain.Product, err error)
	SearchPriceGt(price float64) (newProducts []*domain.Product, err error)
	Create(p *domain.Product) (err error)
	Update(id int, p *domain.Product) (err error)
	Delete(id int) (err error)
}

type repository struct {
	storage store.Store
}

var (
	ErrRepoInternal  = errors.New("internal error")
	ErrRepoNotUnique = errors.New("product already exists")
	ErrRepoNotFound  = errors.New("product not found")
)

func NewRepository(storage store.Store) Repository {
	return &repository{storage}
}

func (r *repository) GetAll() ([]*domain.Product, error) {

	products, err := r.storage.GetAll()
	if err != nil {
		return []*domain.Product{}, ErrRepoInternal
	}
	return products, nil
}

func (r *repository) GetById(id int) (p *domain.Product, err error) {
	product, err := r.storage.GetById(id)
	if err != nil {
		if err == store.ErrStoreInternal {
			err = ErrRepoInternal
		} else {
			err = ErrRepoNotFound
		}
		return
	}
	p = product
	return

}

func (r *repository) SearchPriceGt(price float64) (newProducts []*domain.Product, err error) {
	products, err := r.storage.GetAll()
	if err != nil {
		err = ErrRepoInternal
		return
	}

	for i := range products {
		if products[i].Price > price {
			newProducts = append(newProducts, products[i])
		}
	}
	return
}

func (r *repository) Create(p *domain.Product) (err error) {
	err = r.storage.Create(p)
	if err != nil {
		if err == store.ErrStoreInternal {
			err = ErrRepoInternal
		} else {
			err = ErrRepoNotUnique
		}
		return
	}
	return
}

func (r *repository) Update(id int, p *domain.Product) (err error) {
	err = r.storage.Update(id, p)
	if err != nil {
		if err == store.ErrStoreInternal {
			err = ErrRepoInternal
		} else {
			err = ErrRepoNotFound
		}
		return
	}
	return
}

func (r *repository) Delete(id int) (err error) {
	err = r.storage.Delete(id)
	if err != nil {
		if err == store.ErrStoreInternal {
			err = ErrRepoInternal
		} else {
			err = ErrRepoNotFound
		}
	}
	return
}
