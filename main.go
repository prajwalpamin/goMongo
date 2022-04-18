package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prajwalpamin/goMongo/configs"
	_ "github.com/prajwalpamin/goMongo/docs"
	"github.com/prajwalpamin/goMongo/routes"
)

// @title           Swagger User API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	router := mux.NewRouter()
	configs.ConnectDB()
	routes.UserRoute(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
