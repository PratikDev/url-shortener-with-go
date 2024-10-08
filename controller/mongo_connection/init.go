package mongo_connection

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	const db_name = "netflix-go"
	const collection_name = "watchList"
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
