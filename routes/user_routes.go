package routes

import (
	"github.com/gorilla/mux"
	controller "github.com/prajwalpamin/goMongo/controllers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/", controller.ServeHome).Methods("GET")
	router.HandleFunc("/users", controller.GetUsers).Methods("GET")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}
