package main

import (
	"log"
	"net/http"

	"dating-mobile-app/app/routes"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	r := mux.NewRouter()
	routes.UsersRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":77", r))
}
