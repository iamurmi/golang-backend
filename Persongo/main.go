package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitmpty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content_type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("user").Collection("peoples")
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("user").Collection("peoples")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(people)
}
func GetPeople(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	var person Person
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("user").Collection("peoples")
	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(person)

}
func UpdatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	var person Person
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("user").Collection("peoples")
	filter := bson.M{"_id": id}
	upd := bson.D{
		{"$set", bson.D{
			{"firstname", person.Firstname},
			{"lastname", person.Lastname},
		}},
	}
	collection.FindOneAndUpdate(ctx, filter, upd).Decode(&person)
	json.NewEncoder(response).Encode(person)

}
func DeletePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("user").Collection("peoples")
	filter := bson.M{"_id": id}
	res, _ := collection.DeleteOne(ctx, filter)
	json.NewEncoder(response).Encode(res)
}

var ctx, _ = context.WithTimeout(context.Background(), 1560*time.Second)

func main() {
	fmt.Println("Starting the application")
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://urmikewat:U2Ld5FQvHv3fP5sH@cluster0.fxxvz.mongodb.net/user?retryWrites=true&w=majority"))
	client.Connect(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPeople).Methods("GET")
	router.HandleFunc("/updateperson/{id}", UpdatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/deleteperson/{id}", DeletePersonEndpoint).Methods("DELETE")
	http.ListenAndServe(":12345", router)
	defer client.Disconnect(ctx)

}
