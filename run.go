package cmd

import (
	"fmt"
	"golang_hws/db"
	"golang_hws/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	db.Connect()
	db.Migrate()

	r := mux.NewRouter()
	r.HandleFunc("/transactions", handlers.HandleTransactions)
	r.HandleFunc("/transactions", handlers.Authenticate(handlers.HandleTransactions)).Methods("GET", "POST")
	r.HandleFunc("/transactions/{id}", handlers.HandleTransactions)
	r.HandleFunc("/users", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", handlers.LoginUser).Methods("POST")
	http.Handle("/", r)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
