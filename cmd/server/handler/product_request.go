package handler

import (
	"go_web_ejercicio_dominios/internal/domain"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

func (pr ProductRequest) ToDomain() *domain.Product {
	return &domain.Product{
		Name:        pr.Name,
		Quantity:    pr.Quantity,
		CodeValue:   pr.CodeValue,
		IsPublished: pr.IsPublished,
		Expiration:  pr.Expiration,
		Price:       pr.Price,
	}

}
