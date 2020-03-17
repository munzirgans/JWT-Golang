package main

import (
	"fmt"
	"jwt/pkg/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", config.Login).Methods("POST")
	router.HandleFunc("/signup", config.CreateUser).Methods("POST")
	router.HandleFunc("/logout", config.Logout).Methods("POST")
	router.HandleFunc("/refresh", config.Refresh).Methods("GET")
	router.HandleFunc("/test", config.Test).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Terhubung dengan Port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))
}
