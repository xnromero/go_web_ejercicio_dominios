package handler

import (
	"bytes"
	"encoding/json"
	"go_web_ejercicio_dominios/internal/domain"
	"go_web_ejercicio_dominios/internal/product"
	"go_web_ejercicio_dominios/pkg/store"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func createServer(token string) *gin.Engine {

	if token != "" {
		err := os.Setenv("TOKEN", token)
		if err != nil {
			panic(err)
		}
	}

	storage := store.NewStore("./products_copy.json")
	repository := product.NewRepository(storage)
	service := product.NewService(repository)
	productHandler := NewProductHandler(service)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	group := router.Group("/products")
	{
		group.GET("", productHandler.GetAll())
		group.GET("/:id", productHandler.GetById())
		group.GET("/search", productHandler.SearchPriceGt())
		group.POST("", validateToken, productHandler.Create())
		group.PUT("/:id", validateToken, productHandler.Update())
		group.PATCH("/:id", validateToken, productHandler.UpdatePartial())
		group.DELETE("/:id", validateToken, productHandler.Delete())
	}

	return router
}

func validateToken(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token != os.Getenv("TOKEN") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}
	ctx.Next()

}

func createRequestTest(method string, url string, body string, token string) (*http.Request, *httptest.ResponseRecorder) {

	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("TOKEN", token)
	}
	return req, httptest.NewRecorder()
}

func loadProducts(path string) ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func writeProducts(path string, list []domain.Product) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func TestGetAll(t *testing.T) {

	//Arrange
	var expected = response{Data: []domain.Product{}}

	r := createServer("123456")
	request, response := createRequestTest(http.MethodGet, "/products", "", "123456")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}
	expected.Data = p
	actual := map[string][]domain.Product{}

	//Act
	r.ServeHTTP(response, request)

	//Asserts
	assert.Equal(t, http.StatusOK, response.Code)
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected.Data, actual["data"])
}

func TestGetById(t *testing.T) {
	//Arrange
	var expected = response{Data: domain.Product{
		Id:          1,
		Name:        "Oil - Margarine",
		Quantity:    439,
		CodeValue:   "S82254D",
		IsPublished: true,
		Expiration:  "15/12/2021",
		Price:       71.42,
	}}

	r := createServer("123456")
	request, response := createRequestTest(http.MethodGet, "/products/1", "", "123456")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}
	expected.Data = p[0]
	actual := map[string]domain.Product{}

	//Act
	r.ServeHTTP(response, request)

	//Asserts
	assert.Equal(t, http.StatusOK, response.Code)
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected.Data, actual["data"])
}

func TestCreate(t *testing.T) {
	//Arrange
	var expected = response{Data: domain.Product{
		Id:          500,
		Name:        "Coca Cola",
		Quantity:    500,
		CodeValue:   "test1234",
		IsPublished: true,
		Expiration:  "15/12/2024",
		Price:       21.00,
	}}

	product, err := json.Marshal(expected.Data)
	if err != nil {
		t.Fatal(err)
	}

	r := createServer("123456")
	request, response := createRequestTest(http.MethodPost, "/products", string(product), "123456")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}

	//Act
	r.ServeHTTP(response, request)
	actual := map[string]domain.Product{}

	err = writeProducts("./products_copy.json", p)
	if err != nil {
		t.Fatal(err)
	}

	//Assert
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, expected.Data, actual["data"])

}

func TestDelete(t *testing.T) {

	//Arrange
	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}

	r := createServer("123456")
	request, response := createRequestTest(http.MethodDelete, "/products/500", "", "123456")

	//Act
	r.ServeHTTP(response, request)

	err = writeProducts("./products_copy.json", p)
	if err != nil {
		t.Fatal(err)
	}

	//Assert
	assert.Nil(t, response.Body.Bytes())
	assert.Equal(t, http.StatusNoContent, response.Code)

}

func TestNotFound(t *testing.T) {
	//Arrange
	methods := []string{http.MethodGet, http.MethodPatch, http.MethodDelete}
	r := createServer("123456")

	//Act
	for _, method := range methods {
		request, response := createRequestTest(method, "/products/51231399", "{}", "123456")
		r.ServeHTTP(response, request)
		//Assert
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

func TestUnauthorized(t *testing.T){
	//Arrange
	methods := []string{http.MethodDelete}
	r := createServer("123456")

	//Act
	for _, method := range methods{
		request, response := createRequestTest(method, "/products/101323", "","12345678")
		r.ServeHTTP(response, request)
		//Assert
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}

}
