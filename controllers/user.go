package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"mongo/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

type UserController struct {
	client *mongo.Client
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.client.Database("go-mongo").Collection("users")
	filter := bson.M{"_id": oid}

	var u models.User
	err = collection.FindOne(context.Background(), filter).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := uc.client.Database("go-mongo").Collection("users")
	u.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.client.Database("go-mongo").Collection("users")
	filter := bson.M{"_id": oid}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user", oid, "\n")
}
