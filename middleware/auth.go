package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-essentials/go-mongodb-rest-api/db"
	"go-essentials/go-mongodb-rest-api/models"
	"go-essentials/go-mongodb-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// This method authenticates the token and the user
func Authenticate(ctx *gin.Context) {

	// Extract the token and check authorization
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {

		// Not existent token
		// ABORT NOW the current request
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized user.", errors.New("not authorized user").Error())
		return
	}

	// Check the received token
	userID, err := utils.VerifyToken(token)
	if err != nil {
		// ABORT NOW the current request
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized user.", errors.New("not authorized user").Error())
		return
	}

	fmt.Println("User ID: ", userID)

	// String id not object id
	oID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Get the type to retrieve
	var userRetrieve models.User

	// Check user existence in the database and find the user by id
	// Retrieve a specific document
	collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_USERS)

	// Remove password from the documents if is is the "users" collection
	removePasswordOption := bson.M{"password": 0}
	err = collection.FindOne(context.TODO(), bson.M{"_id": oID}, options.FindOne().SetProjection(removePasswordOption)).Decode(&userRetrieve)

	// Check the error
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "cannot retrieve specific user.", err.Error())
		return
	}

	// Add the user id to the request context for NEXT method
	ctx.Set("userId", userID)

	// Add all the user data to the request context
	jsonData, err := json.Marshal(userRetrieve)
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Error converting to JSON the user data.", err.Error())
		return
	}
	ctx.Set("user", string(jsonData))

	// Call the NEXT function to continue with the request
	ctx.Next()
}
