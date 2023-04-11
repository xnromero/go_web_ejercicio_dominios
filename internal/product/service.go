package product

import (
	"errors"
	"go_web_ejercicio_dominios/internal/domain"
)

type Service interface {
	GetAll() ([]*domain.Product, error)
	GetById(id int) (p *domain.Product, err error)
	SearchPriceGt(price float64) (newProducts []*domain.Product, err error)
	Create(p *domain.Product) (err error)
	Update(id int, p *domain.Product) (err error)
	Delete(id int) (err error)
}

type service struct {
	r Repository
}

var (
	ErrServiceInternal  = errors.New("internal error")
	ErrServiceInvalid   = errors.New("invalid product")
	ErrServiceNotUnique = errors.New("product already exists")
	ErrServiceNotFound  = errors.New("product not found")
)

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAll() ([]*domain.Product, error) {
	products, err := s.r.GetAll()
	if err != nil {
		return []*domain.Product{}, ErrServiceInternal
	}
	return products, nil
}

func (s *service) GetById(id int) (p *domain.Product, err error) {
	p, error := s.r.GetById(id)
	if error != nil {
		if err == ErrRepoInternal {
			err = ErrServiceInternal
		} else {
			err = ErrServiceNotFound
		}
		return
	}
	return
}

func (s *service) SearchPriceGt(price float64) (newProducts []*domain.Product, err error) {
	newProducts, err = s.r.SearchPriceGt(price)
	if err != nil {
		err = ErrServiceInternal
		return
	}
	return
}

func (s *service) Create(p *domain.Product) (err error) {

	err = s.r.Create(p)
	if err != nil {
		if err == ErrRepoInternal {
			err = ErrServiceInternal
		} else {
			err = ErrServiceNotUnique
		}
		return
	}
	return

}

func (s *service) Update(id int, p *domain.Product) (err error) {

	error := s.r.Update(id, p)
	if error != nil {
		if err == ErrRepoInternal {
			err = ErrServiceInternal
		} else {
			err = ErrServiceNotFound
		}
		return
	}

	return
}

func (s *service) Delete(id int) (err error) {

	error := s.r.Delete(id)
	if error != nil {
		if err == ErrRepoInternal {
			err = ErrServiceInternal
		} else {
			err = ErrServiceNotUnique
		}
		return
	}
	return
}
