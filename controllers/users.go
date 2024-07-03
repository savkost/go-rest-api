package controllers

import (
	"context"
	"errors"
	"fmt"
	"go-essentials/go-mongodb-rest-api/db"
	"go-essentials/go-mongodb-rest-api/models"
	"go-essentials/go-mongodb-rest-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// This method signups the user
func Register(ctx *gin.Context) {

	// User Data
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing user data.", err.Error())
		return
	}

	// Hash the user password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error hashing the user password.", err.Error())
		return
	}

	// Update the user password
	user.Password = hashedPassword

	// Created dt and last update dt (YYYY-MM-DD HH:MM:SS)
	NOW_TIME := time.Now().Format("2006-01-02 15:04:05")
	user.CreatedDt = NOW_TIME
	user.LastUpdatedDt = NOW_TIME

	// Set the fields: isAdmin, isActive, role
	// All the users that are using this API are created as normal users
	user.Role = "user"
	user.IsAdmin = "0"
	user.IsActive = "1"

	// Print the data to insert
	fmt.Println("User data to insert: ", user)

	// Insert the user into the database
	collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_USERS)
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error storing the new user in the database.", err.Error())
		return
	}

	// Print the insert result
	fmt.Println("Insert Result:", result)

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully.",
		"data": []map[string]any{
			{
				"_id": result.InsertedID,
			},
		},
	})
}

// This method logs in the user
func Login(ctx *gin.Context) {

	// User Data for login
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing user data for login.", err.Error())
		return
	}

	// Login try
	collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_USERS)
	var result models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&result)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "invalid credentials", err.Error())
		return
	}

	// Set the retrieved password
	var retrievedPassword string = result.Password

	// Compare the stored hashed password with the given password
	validPassword := utils.CheckPasswordHash(user.Password, retrievedPassword)
	if !validPassword {
		utils.HandleError(ctx, http.StatusUnauthorized, "invalid credentials", err.Error())
		return
	}

	// Set the user data
	user = result

	// Create the JWT token for the user
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error logging in the user.", err.Error())
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successfull.",
		"token":   token,
	})
}

// This method creates and ADMIN or SUPERADMIN user
func CreateAdmin[T any](ctx *gin.Context) {

	// User Data
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing user data.", err.Error())
		return
	}

	// Check the role of the user
	if !utils.CheckAllowedRole(user.Role) {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error creating user.", errors.New("the provided role is not supported").Error())
		return
	}

	// Hash the user password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error hashing the user password.", err.Error())
		return
	}

	// Update the user password
	user.Password = hashedPassword

	// Created dt and last update dt (YYYY-MM-DD HH:MM:SS)
	NOW_TIME := time.Now().Format("2006-01-02 15:04:05")
	user.CreatedDt = NOW_TIME
	user.LastUpdatedDt = NOW_TIME

	// Set the fields: isAdmin, isActive
	user.IsAdmin = "1"
	user.IsActive = "1"

	// Print the data to insert
	fmt.Println("User admin data to insert: ", user)

	// Insert the user into the database
	collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_USERS)
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error storing the new admin user in the database.", err.Error())
		return
	}

	// Print the insert result
	fmt.Println("Insert Result:", result)

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Admin User created successfully.",
		"data": []map[string]any{
			{
				"_id": result.InsertedID,
			},
		},
	})
}
