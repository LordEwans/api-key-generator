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

		res, cancel, err := db.resErrHelper("users", newUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": res}})
	}
}
