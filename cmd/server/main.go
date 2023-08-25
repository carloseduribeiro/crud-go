package main

import (
	"github.com/carloseduribeiro/crud-go/configs"
	"github.com/carloseduribeiro/crud-go/internal/entity"
	"github.com/carloseduribeiro/crud-go/internal/infra/database"
	"github.com/carloseduribeiro/crud-go/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig(".")
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
	userHandler := handlers.NewUser(database.NewUser(db), config.TokenAuth, config.JwtExpiresIn)
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})
	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)
	if err = http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
