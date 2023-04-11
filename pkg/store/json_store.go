package store

import (
	"encoding/json"
	"errors"
	"go_web_ejercicio_dominios/internal/domain"
	"os"
)

type Store interface {
	load() ([]*domain.Product, error)
	save(products []*domain.Product) (err error)
	GetAll() (products []*domain.Product, err error)
	GetById(id int) (p *domain.Product, err error)
	Create(p *domain.Product) (err error)
	Update(id int, p *domain.Product) (err error)
	Delete(id int) (err error)
}

type jsonStore struct {
	path string
}

var (
	ErrStoreInternal  = errors.New("internal error")
	ErrStoreNotUnique = errors.New("product already exists")
	ErrStoreNotFound  = errors.New("product not found")
)

func NewStore(path string) Store {
	return &jsonStore{path: path}
}

// Cargar los productos desde archivo json
func (s *jsonStore) load() ([]*domain.Product, error) {
	var products []*domain.Product
	file, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *jsonStore) GetAll() (products []*domain.Product, err error) {

	products, err = s.load()
	if err != nil {
		err = ErrStoreInternal
		return
	}
	return
}

func (s *jsonStore) GetById(id int) (p *domain.Product, err error) {
	products, err := s.load()
	if err != nil {
		err = ErrStoreInternal
		return
	}

	for i := range products {
		if products[i].Id == id {
			p = products[i]
			return
		}
	}
	err = ErrStoreNotFound
	return

}

func (s *jsonStore) Create(p *domain.Product) (err error) {
	products, err := s.load()
	if err != nil {
		err = ErrStoreInternal
		return
	}

	ok := s.productExist(products, p.CodeValue)
	if ok {
		err = ErrStoreNotUnique
		return
	}

	p.Id = len(products) + 1
	products = append(products, p)
	return s.save(products)

}

// guardar los productos en un archivo json
func (s *jsonStore) save(products []*domain.Product) (err error) {
	bytes, errSave := json.Marshal(products)
	if errSave != nil {
		err = errSave
		return err
	}
	return os.WriteFile(s.path, bytes, 0644)
}

func (s *jsonStore) productExist(products []*domain.Product, codeValue string) bool {
	for i := range products {
		if products[i].CodeValue == codeValue {
			return true
		}
	}
	return false
}

func (s *jsonStore) Update(id int, p *domain.Product) (err error) {
	products, err := s.load()
	if err != nil {
		err = ErrStoreInternal
		return
	}

	for i := range products {
		if products[i].Id == id {
			products[i] = p
			products[i].Id = id
			return s.save(products)
		}
	}
	err = ErrStoreNotFound
	return
}

func (s *jsonStore) Delete(id int) (err error) {
	products, err := s.load()
	if err != nil {
		err = ErrStoreInternal
		return
	}

	for i := range products {
		if products[i].Id == id {
			products = append(products[:i], products[i+1:]...)
			return s.save(products)
		}
	}

	err = ErrStoreNotFound
	return
}
