package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/juancruzestevez/go-gorm-restapi/db"
	"github.com/juancruzestevez/go-gorm-restapi/middleware"
	"github.com/juancruzestevez/go-gorm-restapi/models"
	"github.com/juancruzestevez/go-gorm-restapi/routes"
	"github.com/juancruzestevez/go-gorm-restapi/utils"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Puerto configurado:", port)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET no definido")
	}
	log.Println("JWT_SECRET cargado")

	utils.InitJWT(jwtSecret)

	db.DBConnection()

	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.Task{})

	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/auth/register", routes.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", routes.LoginHandler).Methods("POST")

	// User routes
	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods("DELETE")

	// Task routes
	r.HandleFunc("/tasks", routes.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", routes.GetTaskHandler).Methods("GET")
	
	r.Handle("/tasks", 
		middleware.AuthMiddleware(
			http.HandlerFunc(routes.CreateTaskHandler),
		)).Methods("POST")

	r.Handle("/tasks/{id}", 
		middleware.AuthMiddleware(
			http.HandlerFunc(routes.DeleteTaskHandler),
		)).Methods("DELETE")

	log.Println("Servidor escuchando en el puerto " + port)
	http.ListenAndServe(":"+port, r)
}