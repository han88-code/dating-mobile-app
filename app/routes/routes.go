package routes

import (
	"dating-mobile-app/app/controllers"

	"github.com/gorilla/mux"
)

var UsersRoutes = func(router *mux.Router) {
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")
	router.HandleFunc("/home", controllers.Home).Methods("POST")
	router.HandleFunc("/swipe", controllers.Swipe).Methods("POST")
	router.HandleFunc("/verified", controllers.VerifiedUser).Methods("POST")
}
