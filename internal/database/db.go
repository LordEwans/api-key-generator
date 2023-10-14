package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/lordewans/api-key-generator/configs"
	"github.com/lordewans/api-key-generator/internal"
	"github.com/lordewans/api-key-generator/internal/models"
	"github.com/lordewans/api-key-generator/internal/responses"
	"github.com/lordewans/api-key-generator/pkg/generate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

type DB struct {
	client *mongo.Client
}

func ConnectDB() (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configs.EnvMongoURI()))
	internal.Handle(err)

	internal.Handle(err)

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to MongoDB!")
	return &DB{client: client}, err
}

func colHelper(db *DB, collectionName string) *mongo.Collection {
	return db.client.Database("APIKeys").Collection(collectionName)
}

func (db *DB) ctxDeferHelper(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc) {
	collection := colHelper(db, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	return collection, ctx, cancel
}

func (db *DB) resErrHelper(collectionName string, input any) (*mongo.InsertOneResult, context.CancelFunc, error) {
	collection, ctx, cancel := db.ctxDeferHelper(collectionName)

	res, err := collection.InsertOne(ctx, input)

	internal.Handle(err)

	return res, cancel, err
}

func (db *DB) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.User{
			Username: user.Username,
			APIKey:   generate.GenerateKey(),
			API:      user.API,
		}

		res, cancel, err := db.resErrHelper(user.API, newUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": map[string]interface{}{
			"key": newUser.APIKey,
			"id":  res,
		}}})
	}
}

func (db *DB) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		api, id := c.Param("api"), c.Param("_id")
		collection, ctx, cancel := db.ctxDeferHelper(api)
		defer cancel()

		ID, _ := primitive.ObjectIDFromHex(id)

		result, err := collection.DeleteOne(ctx, bson.M{"_id": ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}})
	}
}

func (db *DB) VerifyKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		api, apiKey := c.Param("api"), c.Param("api-key")
		collection, ctx, cancel := db.ctxDeferHelper(api)
		defer cancel()
		fmt.Printf("%s, %s\n", api, apiKey)

		if err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user); err != nil {
			fmt.Printf("Found 0 results for API Key: %s\n", apiKey)
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Authentication failed"})
			internal.Handle(err)
			return
		}
	}
}

func (db *DB) VerifyKeyAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		api, apiKey := "api-key-generator", c.Param("api-key")
		collection, ctx, cancel := db.ctxDeferHelper(api)
		defer cancel()

		if err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user); err != nil {
			fmt.Printf("Found 0 results for API Key: %s\n", apiKey)
			c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Authentication failed"})
			internal.Handle(err)
			return
		}
	}
}
