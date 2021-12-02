package main

import (
	"business/config"
	"business/handlers"
	"business/storage"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var cl storage.Client
var stf storage.Staff

func init() {
	_ = godotenv.Load()
	host := os.Getenv("DATABASE_HOST")
	dname := os.Getenv("DATABASE_NAME")
	dbPort := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USERNAME")

	c := config.Config{
		DatabaseHost:     host,
		DatabaseName:     dname,
		DatabasePort:     dbPort,
		DatabaseUsername: user,
	}
	cl = storage.NewClient(c)
	stf = storage.NewStaff(cl)

}

func main() {
	sh := handlers.NewStaffHandler(stf)

	router := mux.NewRouter()
	router.HandleFunc("/staff/table", sh.AutoMigrate).Methods("POST")
	router.HandleFunc("/staff/create", sh.SignUp).Methods("POST")
	router.HandleFunc("/product/create", sh.Save).Methods("POST")
	router.HandleFunc("/staff/login", sh.SignIn).Methods("GET")
	router.HandleFunc("/product/product", sh.Product).Methods("GET")
	router.HandleFunc("/product/products", sh.Products).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}
