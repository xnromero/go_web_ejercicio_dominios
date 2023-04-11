package handler

import (
	"encoding/json"
	"errors"
	"go_web_ejercicio_dominios/internal/product"
	"go_web_ejercicio_dominios/pkg/web"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	s product.Service
}

var (
	ErrHandlerInternal       = errors.New("internal error")
	ErrHandlerInvalid        = errors.New("invalid product")
	ErrHandlerNotUnique      = errors.New("product already exists")
	ErrHandlerNotFound       = errors.New("product not found")
	ErrHandlerInvalidRequest = errors.New("invalid request")
)

func NewProductHandler(s product.Service) *productHandler {
	return &productHandler{s}
}

func (h *productHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := h.s.GetAll()
		if err != nil {
			web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
		}
		web.Success(ctx, http.StatusOK, products)
	}
}

func (h *productHandler) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			return
		}

		prod, err := h.s.GetById(id)
		if err != nil {
			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusNotFound, ErrHandlerNotFound)
			}
			return
		}
		web.Success(ctx, http.StatusOK, prod)
	}
}

func (h *productHandler) SearchPriceGt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
		}
		newProducts, errSearch := h.s.SearchPriceGt(price)
		if errSearch != nil {
			web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			return
		}
		web.Success(ctx, http.StatusOK, newProducts)
	}
}

func (h *productHandler) Create() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var req ProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			log.Println("log: ", err)
			return
		}

		prod := req.ToDomain()

		err := h.s.Create(prod)
		if err != nil {
			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusConflict, ErrHandlerNotUnique)
			}
			return
		}
		web.Success(ctx, http.StatusCreated, prod)

	}
}

func (h *productHandler) Update() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			return
		}

		var req ProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			log.Println("log: ", err)
			return
		}

		prod := req.ToDomain()

		err = h.s.Update(id, prod)
		if err != nil {

			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusNotFound, ErrHandlerNotFound)
			}
			return
		}
		web.Success(ctx, http.StatusOK, prod)

	}
}

func (h *productHandler) UpdatePartial() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			return
		}

		p, err := h.s.GetById(id)
		if err != nil {
			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusNotFound, ErrHandlerNotFound)
			}
			return
		}

		if err := json.NewDecoder(ctx.Request.Body).Decode(&p); err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			return
		}
		p.Id = id

		if err := h.s.Update(id, p); err != nil {
			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusConflict, ErrHandlerNotFound)
			}
			return
		}
		web.Success(ctx, http.StatusOK, p)
	}
}

func (h *productHandler) Delete() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, ErrHandlerInvalidRequest)
			return
		}

		err = h.s.Delete(id)
		if err != nil {
			if err == product.ErrServiceInternal {
				web.Failure(ctx, http.StatusInternalServerError, ErrHandlerInternal)
			} else {
				web.Failure(ctx, http.StatusNotFound, ErrHandlerNotFound)
			}
			return
		}
		web.Success(ctx, http.StatusNoContent, "removed product")

	}
}

