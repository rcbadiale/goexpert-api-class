package main

import (
	"goexpert-api/configs"
	_ "goexpert-api/docs"
	"goexpert-api/internal/entity"
	"goexpert-api/internal/infra/database"
	"goexpert-api/internal/infra/webserver/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     This is a simple API made for the Go Expert course.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Rafael Badiale
// @contact.email  nobody@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
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
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	// Creating services
	// Products
	productService := database.NewProductService(db)
	productHandler := handlers.NewProductHandler(productService)
	// User
	userService := database.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, config.TokenAuth, config.JWTExpiresIn)

	// Using Chi as router
	r := chi.NewRouter()

	// General middlewares
	// r.Use(middleware.Logger) // Chi Logger
	r.Use(LogRequest) // Custom Logger

	r.Route("/products", func(r chi.Router) {
		// Group middlewares
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		// Routes
		r.Get("/", productHandler.GetProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/user", func(r chi.Router) {
		// Routes
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	http.ListenAndServe(":8000", r)
}

// Custom Middleware
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
