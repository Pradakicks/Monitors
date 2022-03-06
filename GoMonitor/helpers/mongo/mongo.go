package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// Replace with your mongoDB Cluster URI
const uriLink = "mongodb+srv://testPradaMonitor:ThisWillBeOpenSourcedMongoPrada@cluster0.yikgy.mongodb.net/monitorDB?retryWrites=true&w=majority"
// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI(uriLink)
	// clientOptions := options.Client().ApplyURI("mongodb+srv://vibris-User:eIDpR4kttFu57FHE@vibris.jyxhh.mongodb.net/testing?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Shopify!")

	collection := client.Database("monitorDB").Collection("Shopify")
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return collection
}
func ConnectDBShopifyLink() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI(uriLink)
	// clientOptions := options.Client().ApplyURI("mongodb+srv://vibris-User:eIDpR4kttFu57FHE@vibris.jyxhh.mongodb.net/testing?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Shopify!")

	collection := client.Database("monitorDB").Collection("ShopifyLink")
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return collection
}
func ConnectDBMain() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI(uriLink)
	// clientOptions := options.Client().ApplyURI("mongodb+srv://vibris-User:eIDpR4kttFu57FHE@vibris.jyxhh.mongodb.net/testing?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Main!")

	collection := client.Database("monitorDB").Collection("Main")
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return collection
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}


// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
