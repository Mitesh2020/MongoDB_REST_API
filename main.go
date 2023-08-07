package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"mongo/controllers"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set the MongoDB connection URI
	uri := "mongodb://localhost:27017" // Replace with your MongoDB URI

	// Set options for the MongoDB client
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	} else {
		fmt.Println("Connected to MongoDB Server")
	}
	defer client.Disconnect(context.Background())

	r := httprouter.New()
	uc := controllers.NewUserController(client)
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	serverAddr := "localhost:8080"
	fmt.Printf("Server is listening on %s...\n", serverAddr)
	serverErr := http.ListenAndServe(serverAddr, r) // Change variable name to serverErr
	if serverErr != nil {
		log.Fatalf("Error starting the server: %v", serverErr)
	}
}

/* func getSession() *mgo.Session {
	// Set up the MongoDB connection string
	uri := "mongodb://localhost:27017"

	// Connect to MongoDB
	session, err := mgo.Dial(uri)
	if err != nil {
		if err.Error() == "no reachable servers" {
			fmt.Println("MongoDB server is not running.")
			return nil
		} else {
			log.Fatalf("Error connecting to MongoDB: %v", err)
		}
	}

	// Set mode and other session configurations if needed
	session.SetMode(mgo.Monotonic, true)

	fmt.Println("Connected to MongoDB!")

	return session
}
*/
