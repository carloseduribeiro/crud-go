package main

import (
	"github.com/carloseduribeiro/crud-go/configs"
	_ "github.com/carloseduribeiro/crud-go/docs"
	"github.com/carloseduribeiro/crud-go/internal/entity"
	"github.com/carloseduribeiro/crud-go/internal/infra/database"
	"github.com/carloseduribeiro/crud-go/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with authentication developed during Full Cycle's Go Expert course
// @termsOfService  http://swagger.io/terms/

// @contact.name   Carlos Eduardo Ribeiro
// @contact.url    https://www.linkedin.com/in/carloseduardoribeiro96/
// @contact.email  carloseribeiro96@gmail.com

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apiKey ApiKeyAuth
// @in HEADER
// @name Authorization
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
	r.Use(middleware.Recoverer)
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Post("/generate_token", userHandler.GetJWT)
	})
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	if err = http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
