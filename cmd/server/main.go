package main

import (
	"github.com/carloseduribeiro/crud-go/configs"
	"github.com/carloseduribeiro/crud-go/internal/entity"
	"github.com/carloseduribeiro/crud-go/internal/infra/database"
	"github.com/carloseduribeiro/crud-go/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&entity.Product{}, &entity.User{}); err != nil {
		panic(err)
	}
	productHandler := handlers.NewProductHandler(database.NewProduct(db))
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Post("/products", productHandler.CreateProduct)
	if err = http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
