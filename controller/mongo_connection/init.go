package mongo_connection

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
	"github.com/pratikdev/url-shortner-with-go/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var Collection *mongo.Collection
var UserCollection *mongo.Collection

// initiating mongodb connection
func init() {
	// loading env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	const db_name = "short-url"
	const collection_name = "urls"
	const users_collection_name = "users"
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
	UserCollection = client.Database(db_name, dbOption).Collection(users_collection_name)

	fmt.Println("Collection instances are ready")
}

// HashPassword hashes a plaintext password using bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a hashed password with a plaintext password
func checkPasswordHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// login user
func LoginUser(ld models.LoginDetails) (models.User, error) {
	filter := bson.M{"username": ld.Username}
	result := UserCollection.FindOne(context.Background(), filter)

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, &customErrors.CustomError{Code: http.StatusNotFound, Message: "No user found with the given name"}
		}

		return models.User{}, &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if !checkPasswordHash(user.Password, ld.Password) {
		return user, &customErrors.CustomError{Code: http.StatusUnauthorized, Message: "Invalid credentials"}
	}

	return user, nil
}

// register user
func RegisterUser(user models.LoginDetails) error {
	// check if user already exists
	filter := bson.M{"username": user.Username}
	existingUserResult := UserCollection.FindOne(context.Background(), filter)

	var existingUser models.User
	if err := existingUserResult.Decode(&existingUser); err == nil {
		return &customErrors.CustomError{Code: http.StatusConflict, Message: "Username taken"}
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	_, docCreationErr := UserCollection.InsertOne(context.Background(), models.User{
		Username: user.Username,
		Password: hashedPassword,
	})
	if docCreationErr != nil {
		return &customErrors.CustomError{Code: http.StatusInternalServerError, Message: docCreationErr.Error()}
	}

	return nil
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
	if _, ok := urlObject["url"].(string); ok {
		return urlObject["url"].(string), nil
	}

	return "", &customErrors.CustomError{Code: http.StatusNotFound, Message: "URL not found"}
}

// Create new URL
func CreateNewURL(data models.NewURL) (string, error) {
	result, err := Collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

// Get all URLs
func GetAllURLs(authorId string) ([]models.URL, error) {
	authorObjectId, err := bson.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Invalid ID"}
	}

	filter := bson.M{"author": authorObjectId}
	cursor, err := Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	var urls []models.URL
	if err = cursor.All(context.Background(), &urls); err != nil {
		return nil, &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return urls, nil
}
