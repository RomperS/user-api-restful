package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	httpHandler "user-api-restful/cmd/api/http"
	"user-api-restful/internal/application"
	"user-api-restful/internal/persistence/database"
	"user-api-restful/internal/persistence/entity"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("PASSWORD_")
	dbHost := os.Getenv("DB_HOST")
	dbName := "users_db"
	dbPort := "5432"

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&entity.UserEntity{})

	if err != nil {
		log.Fatal("failed to auto migrate users: ", err)
	}

	userRepository := database.NewPostgresRepository(db)

	userService := application.NewUserServiceImpl(userRepository, userRepository)

	userHandler := httpHandler.NewUserHandler(userService)

	router := chi.NewRouter()

	router.Use(httpHandler.AuthAndLoggingMiddleware)

	router.Route("/users", func(r chi.Router) {
		// POST /users - Create a new user
		r.Post("/", httpHandler.ErrorHandlerWrapper(userHandler.CreateUser))

		// GET /users - Retrieve all users (FindAll)
		r.Get("/", httpHandler.ErrorHandlerWrapper(userHandler.FindAll))

		// PUT /users - Update an existing user (Update)
		// Common pattern: Use PUT to replace the entire resource, often including the ID in the body.
		r.Put("/", httpHandler.ErrorHandlerWrapper(userHandler.Update))

		// GET /users/{id} - Retrieve a specific user by ID (FindById)
		// The '{id}' is a URL parameter that userHandler.FindById needs to extract.
		r.Get("/{id}", httpHandler.ErrorHandlerWrapper(userHandler.FindById))

		// DELETE /users/{id} - Delete a specific user by ID (Delete)
		r.Delete("/{id}", httpHandler.ErrorHandlerWrapper(userHandler.Delete))
	})

	log.Printf("Server starting on port :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
