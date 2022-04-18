package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/prajwalpamin/goMongo/configs"
	_ "github.com/prajwalpamin/goMongo/docs"
	model "github.com/prajwalpamin/goMongo/models"
	"github.com/prajwalpamin/goMongo/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

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
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	newUser := model.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
	}
	result, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User inserted with insertid %v ", result.InsertedID)

}

// GetUser godoc
// @Summary Get details of  user
// @Description Get details of  user
// @Tags user
// @Accept  json
// @Produce  json
// @Param userId   path int true "User ID"
// @Success 200 {array} model.User
// @Router /user/{userId} [get]

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(r)
	userId := params["userId"]
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusBadRequest, Message: "success", Data: map[string]interface{}{"data": user}}
	json.NewEncoder(w).Encode(response)
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

	result, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Panic(err)
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

// UpdateUser godoc
// @Summary Update a user
// @Description Update user with the input paylod
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body model.User true "Update user"
// @Success 200 {object} model.User
// @Router /user/{userId} [put]

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vars := mux.Vars(r)
	userId := vars["userId"]
	user := model.User{}
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	if validationErr := validate.Struct(&user); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	update := bson.M{"name": user.Name, "email": user.Email, "password": hash}
	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	updatedUser := model.User{}
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
	json.NewEncoder(w).Encode(response)
	return

}

// @Summary delete a user by ID
// @ID delete-user
// @Produce json
// @Param userId path string true "user ID"
// @Success 200 {object} model.User
// @Router /user/{userId} [delete]

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(r)
	userId := params["userId"]
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "user not found"}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "user successfully deleted"}}
	json.NewEncoder(w).Encode(response)
}
