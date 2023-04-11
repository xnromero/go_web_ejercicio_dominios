package main

import (
	"go_web_ejercicio_dominios/cmd/server/handler"
	"go_web_ejercicio_dominios/internal/product"
	"go_web_ejercicio_dominios/pkg/store"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}

	storage := store.NewStore("../../products.json")
	repository := product.NewRepository(storage)
	service := product.NewService(repository)
	productHandler := handler.NewProductHandler(service)

	router := gin.Default()

	group := router.Group("/products")
	{
		group.GET("", productHandler.GetAll())
		group.GET("/:id", productHandler.GetById())
		group.GET("/search", productHandler.SearchPriceGt())
		group.POST("", validateToken, productHandler.Create())
		group.PUT("/:id", validateToken, productHandler.Update())
		group.PATCH("/:id", validateToken, productHandler.UpdatePartial())
		group.DELETE("/:id",validateToken, productHandler.Delete())
	}
	router.Run()
}

func validateToken(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token != os.Getenv("TOKEN") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}
	ctx.Next()

}
