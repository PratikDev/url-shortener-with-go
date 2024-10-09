package mongo_connection

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Collection *mongo.Collection

// initiating mongodb connection
func init() {
	// loading env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	const db_name = "short-url"
	const collection_name = "urls"
	connection_string := os.Getenv("MONGODB_CONNECTION_STRING")

	// client options
	clientOption := options.Client().ApplyURI(connection_string)

	// database options
	dbOption := options.Database()

	// connect to db
	client, err := mongo.Connect(clientOption)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(db_name, dbOption).Collection(collection_name)

	fmt.Println("Collection instance is ready")
}

// Get URL from id
func GetURLFromId(id string) (string, error) {
	urlObjectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return "", &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Invalid ID"}
	}

	filter := bson.M{"_id": urlObjectId}
	result := Collection.FindOne(context.Background(), filter)

	var urlObject bson.M
	result.Decode(&urlObject)

	// check if url key exists in the object
	if url, ok := urlObject["url"].(string); ok {
		return url, nil
	}

	return "", &customErrors.CustomError{Code: http.StatusNotFound, Message: "URL not found"}

}
