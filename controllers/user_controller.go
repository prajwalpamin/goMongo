package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prajwalpamin/goMongo/configs"
	_ "github.com/prajwalpamin/goMongo/docs"
	model "github.com/prajwalpamin/goMongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func ServeHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("welcome Home")
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input paylod
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body model.User true "Create user"
// @Success 200 {object} model.User
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	//if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	//to generate unique id
	rand.Seed(time.Now().UnixNano())
	user.Id = strconv.Itoa(rand.Intn(100))
	newUser := model.User{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	result, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User inserted with insertid %v", result.InsertedID)

}

// GetUsers godoc
// @Summary Get details of all users
// @Description Get details of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} model.User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	result, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close(context.Background())

	for result.Next(context.Background()) {
		var sinlgeUser model.User
		err := result.Decode(&sinlgeUser)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, sinlgeUser)
	}
	json.NewEncoder(w).Encode(users)
}
