package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-essentials/go-mongodb-rest-api/db"
	"go-essentials/go-mongodb-rest-api/utils"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Private
// This method converts the model name to the corresponding collection name
func modelToCollectionName(modelName string) (string, error) {
	switch modelName {
	case "models.User":
		// USERS
		return db.DB_TABLE_USERS, nil

	case "models.LicenseCategory":
		// LICENSE CATEGORIES
		return db.DB_TABLE_LICENSES_CATEGORIES, nil

	case "models.License":
		// LICENSES
		return db.DB_TABLE_LICENCES, nil

	default:
		return "", errors.New("not found corresponding model to collection name: " + modelName)
	}
}

// CREATE DOCUMENT -------------
// -----------------------------
func CreateDocument[T any](ctx *gin.Context) {

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Terminate the create document action if it is the 'users' API
	if collectionName == db.DB_TABLE_USERS {
		messageShowError := "create document is not allowed for 'users'. Use register API or createAdmin."
		utils.HandleError(ctx, http.StatusMethodNotAllowed, messageShowError, messageShowError)
		return
	}

	// Retrieve and read the request body
	requestBody := make(map[string]interface{})
	decoderData := json.NewDecoder(ctx.Request.Body)

	// Cast all integers or floats as Number
	decoderData.UseNumber()

	// Decode the data
	err = decoderData.Decode(&requestBody)
	fmt.Println("Request body:", requestBody)
	fmt.Printf("Type:%T\n", requestBody)

	// Check if error OR the request map is empty
	if err != nil || len(requestBody) == 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing and decoding document data.", err.Error())
		return
	}

	// Created dt and last update dt (YYYY-MM-DD HH:MM:SS)
	NOW_TIME := time.Now().Format("2006-01-02 15:04:05")
	requestBody["createdDt"] = NOW_TIME
	requestBody["lastUpdatedDt"] = NOW_TIME
	requestBody["isActive"] = "1"

	// Print the data to insert
	fmt.Println("Document data to insert: ", requestBody)

	// Insert the document into the database
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)
	result, err := collection.InsertOne(context.TODO(), requestBody)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error storing the new document in the database.", err.Error())
		return
	}

	// Print the insert result
	fmt.Println("Insert Result:", result)

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Document inserted successfully in the collection { " + collectionName + " }.",
		"data": []map[string]any{
			{
				"_id": result.InsertedID,
			},
		},
	})
}

// GET ALL DOCUMENTS -----------
// -----------------------------
func GetAllDocuments[T any](ctx *gin.Context) {

	// Console the Query
	fmt.Println("Request Query:", ctx.Request.URL)

	// Console all the query parameters
	queryGivenParams := ctx.Request.URL.Query()
	fmt.Println("Query Parameters:", queryGivenParams)

	// PAGING - SKIP
	pageNumberString := ctx.Query("page")
	fmt.Println("Page requested (string):", pageNumberString)

	var pageRequested int64 = 1
	if utils.CheckStringNotEmpty(pageNumberString) {
		pageNumberLocal, err := strconv.ParseInt(pageNumberString, 10, 64)
		if err != nil || pageNumberLocal <= 0 {
			errorMsg := "Error transforming page number to int64 OR page number is lower/equal to ZERO."
			utils.HandleError(ctx, http.StatusBadRequest, errorMsg, errorMsg)
			return
		}
		fmt.Println("Page number (INT PARSED):", pageNumberLocal)
		pageRequested = pageNumberLocal
	}

	// LIMIT RESULTS - LIMIT
	limitResultString := ctx.Query("limit")
	fmt.Println("Limit Results (string):", limitResultString)

	var limit int64 = 100
	if utils.CheckStringNotEmpty(limitResultString) {
		limitLocal, err := strconv.ParseInt(limitResultString, 10, 64)
		if err != nil || limitLocal <= 0 {
			errorMsg := "Error transforming limit to int64 OR limit is lower/equal to ZERO."
			utils.HandleError(ctx, http.StatusBadRequest, errorMsg, errorMsg)
			return
		}
		fmt.Println("Limit Results (INT PARSED):", limitLocal)
		limit = limitLocal
	}

	// Set the query search options
	searchOpts := options.Find()
	skipNumberValues := (pageRequested - 1) * limit
	searchOpts.SetSkip(skipNumberValues).SetLimit(limit)

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Retrieve all the collection documents
	var rows *mongo.Cursor
	var err error

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)

	// Find the total number of documents in the requested collection
	countDocumentsResult, err := collection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error counting the documents in the collection.", err.Error())
		return
	}

	// Check the requested page number if it is viable
	if skipNumberValues > countDocumentsResult {
		utils.HandleError(ctx, http.StatusInternalServerError, "The requested page is not existent. Page: "+pageNumberString, "The requested page is not existent.")
		return
	}

	// Remove password from the documents if is is the "users" collection
	if collectionName == db.DB_TABLE_USERS {
		rows, err = collection.Find(context.TODO(), bson.M{}, searchOpts.SetProjection(bson.M{"password": 0}))
	} else {
		rows, err = collection.Find(context.TODO(), bson.M{}, searchOpts)
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error retrieving documents list.", err.Error())
		return
	}

	// Return all documents data
	documentsList := []T{}
	if err = rows.All(context.TODO(), &documentsList); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error retrieving documents list.", err.Error())
		return
	}

	// Decode all the received documents
	for rows.Next(context.TODO()) {

		// Local Variable
		var resDoc T

		err = rows.Decode(&resDoc)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error decoding documents list.", err.Error())
			return
		}
		documentsList = append(documentsList, resDoc)
	}

	// Calculate the total number of pages remaining in the requested collection
	totalPages := math.Floor(float64(countDocumentsResult) / float64(limit))
	if math.Ceil(float64(totalPages*float64(limit))) < float64(countDocumentsResult) {
		totalPages += 1
	}
	fmt.Println("Total Pages:", totalPages)

	// Send the response with all the documents
	// The GIN package will automatically encode the response in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"message":               "Retrieved documents successfully. Retrieved {" + fmt.Sprint(len(documentsList)) + "} Documents from Collection { " + collectionName + " }.",
		"data":                  documentsList,
		"rows":                  len(documentsList),
		"currentPage":           pageRequested,
		"totalPages":            totalPages,
		"totalNumbersDocuments": countDocumentsResult,
	})
}

// GET A DOCUMENT BY ID --------
// -----------------------------
func GetDocumentByID[T any](ctx *gin.Context) {

	// Error variable
	var err error

	// Retrieve the document id from the parameters
	documentID := ctx.Param("id")
	if documentID == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "Error finding document ID.", errors.New("error finding document ID").Error())
		return
	}

	// String id not object id
	oID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Retrieve a specific document
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)

	// Remove password from the documents if is is the "users" collection
	if collectionName == db.DB_TABLE_USERS {
		err = collection.FindOne(context.TODO(), bson.M{"_id": oID}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&docRetrieve)
	} else {
		err = collection.FindOne(context.TODO(), bson.M{"_id": oID}).Decode(&docRetrieve)
	}

	// Check the error
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "cannot retrieve specific document.", err.Error())
		return
	}

	// Return the found document
	// The GIN package will automatically encode the response in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Document retrieved successfully from Collection { " + collectionName + " }.",
		"data":    docRetrieve,
	})
}

// GET LAST X DOCUMENTS --------
// -----------------------------
func GetLastXDocuments[T any](ctx *gin.Context) {

	// Retrieve all the URL Parameters ------------
	// --------------------------------------------
	fmt.Println("Request Query:", ctx.Request.URL)

	// Console all the query parameters
	queryGivenParams := ctx.Request.URL.Query()
	fmt.Println("Query Parameters:", queryGivenParams)

	// PAGING - SKIP
	pageNumberString := ctx.Query("page")
	fmt.Println("Page requested (string):", pageNumberString)

	var pageRequested int64 = 1
	if utils.CheckStringNotEmpty(pageNumberString) {
		pageNumberLocal, err := strconv.ParseInt(pageNumberString, 10, 64)
		if err != nil || pageNumberLocal <= 0 {
			errorMsg := "Error transforming page number to int64 OR page number is lower/equal to ZERO."
			utils.HandleError(ctx, http.StatusBadRequest, errorMsg, errorMsg)
			return
		}
		fmt.Println("Page number (INT PARSED):", pageNumberLocal)
		pageRequested = pageNumberLocal
	}

	// LIMIT RESULTS - LIMIT
	limitResultString := ctx.Query("limit")
	fmt.Println("Limit Results (string):", limitResultString)

	var limit int64 = 100
	if utils.CheckStringNotEmpty(limitResultString) {
		limitLocal, err := strconv.ParseInt(limitResultString, 10, 64)
		if err != nil || limitLocal <= 0 {
			errorMsg := "Error transforming limit to int64 OR limit is lower/equal to ZERO."
			utils.HandleError(ctx, http.StatusBadRequest, errorMsg, errorMsg)
			return
		}
		fmt.Println("Limit Results (INT PARSED):", limitLocal)
		limit = limitLocal
	}

	// Set the query search options
	skipNumberValues := (pageRequested - 1) * limit

	// Created Date From
	createdDateFrom := ctx.Query("created_dtFrom")
	fmt.Println("Created Date From:", createdDateFrom)

	// Created Date To
	createdDateTo := ctx.Query("created_dtTo")
	fmt.Println("Created Date To:", createdDateTo)

	// Retrieve all the collection documents
	var rows *mongo.Cursor
	var err error

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)

	// Find the total number of documents in the requested collection
	countDocumentsResult, err := collection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error counting the documents in the collection.", err.Error())
		return
	}

	// Check the requested page number if it is viable
	if skipNumberValues > countDocumentsResult {
		utils.HandleError(ctx, http.StatusInternalServerError, "The requested page is not existent. Page: "+pageNumberString, "The requested page is not existent.")
		return
	}

	// Search with the given filters
	// If not given any filters, then the API returns all the users
	// Construct the filter query and the options
	filterObject := bson.M{}
	opts := options.Find()
	opts.SetSkip(skipNumberValues).SetLimit(limit)

	// Check all the fields --------
	// -----------------------------

	// Created Date From, To
	if utils.CheckStringNotEmpty(createdDateFrom) && utils.CheckStringNotEmpty(createdDateTo) {
		filterObject["createdDt"] = bson.M{
			"$gte": createdDateFrom,
			"$lte": createdDateTo,
		}
	} else {
		if utils.CheckStringNotEmpty(createdDateFrom) {
			filterObject["createdDt"] = bson.M{"$gte": createdDateFrom}
		} else if utils.CheckStringNotEmpty(createdDateTo) {
			filterObject["createdDt"] = bson.M{"$lte": createdDateTo}
		}
	}

	// Retrieve last X documents from the database
	// Remove password from the documents if is is the "users" collection
	if collectionName == db.DB_TABLE_USERS {
		rows, err = collection.Find(context.TODO(), filterObject, opts.SetProjection(bson.M{"password": 0}))
	} else {
		rows, err = collection.Find(context.TODO(), filterObject, opts)
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Error retrieving the last X documents.", err.Error())
		return
	}

	// Return all documents data
	documentsList := []T{}
	if err = rows.All(context.TODO(), &documentsList); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error retrieving documents list.", err.Error())
		return
	}

	// Decode all the received documents
	for rows.Next(context.TODO()) {

		// Local Variable
		var resDoc T

		err = rows.Decode(&resDoc)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error decoding documents list.", err.Error())
			return
		}
		documentsList = append(documentsList, resDoc)
	}

	// Calculate the total number of pages remaining in the requested collection
	totalPages := math.Floor(float64(countDocumentsResult) / float64(limit))
	if math.Ceil(float64(totalPages*float64(limit))) < float64(countDocumentsResult) {
		totalPages += 1
	}
	fmt.Println("Total Pages:", totalPages)

	// Send the response with all the documents
	// The GIN package will automatically encode the response in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"message":               "Retrieved documents successfully. Retrieved {" + fmt.Sprint(len(documentsList)) + "} Documents from Collection { " + collectionName + " }.",
		"data":                  documentsList,
		"rows":                  len(documentsList),
		"currentPage":           pageRequested,
		"totalPages":            totalPages,
		"totalNumbersDocuments": countDocumentsResult,
	})
}

// COUNT ALL DOCUMENTS -----
// -------------------------
func CountAllDocuments[T any](ctx *gin.Context) {

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Count all documents in the provided collection
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)
	countDocumentsResult, err := collection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error counting the documents in the collection.", err.Error())
		return
	}

	// Send the response with the count of documents
	// The GIN package will automatically encode the response in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Counted {" + fmt.Sprint(countDocumentsResult) + "} Documents in the Collection { " + collectionName + " }.",
		"rows":    countDocumentsResult,
	})
}

// DELETE DOCUMENT ---------
// -------------------------
func DeleteDocument[T any](ctx *gin.Context) {

	// Retrieve the user id from the parameters
	documentID := ctx.Param("id")
	if documentID == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "Error finding document ID.", errors.New("error finding document ID").Error())
		return
	}

	// String id not object id
	oID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Count all documents in the provided collection
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)

	// Delete the specific document from the collection
	// Set the delete statement
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": oID})
	if err != nil || result.DeletedCount <= 0 {
		utils.HandleError(ctx, http.StatusInternalServerError, "deletion of the document data failed", errors.New("deletion of the document data failed").Error())
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully.",
	})
}

// DELETE MULTIPLE DOCUMENTS
// -------------------------
func DeleteMultipleDocuments[T any](ctx *gin.Context) {

	// Documents IDs Struct
	type usersIds struct {
		ListOfIds []string `json:"listOfIds"`
	}

	// Extract the list of documents IDs to delete from the request body
	var documentsIdsListObject usersIds

	decoderData := json.NewDecoder(ctx.Request.Body)

	// Cast all integers or floats as Number
	decoderData.UseNumber()

	// Decode the data
	err := decoderData.Decode(&documentsIdsListObject)
	fmt.Println("List of documents IDs:", documentsIdsListObject.ListOfIds)

	// Check the list of documents IDs
	if len(documentsIdsListObject.ListOfIds) <= 0 || err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "No documents IDs found to delete.", errors.New("no documents IDs found to delete").Error())
		return
	}

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Convert every id to Object ID
	objectIDs := []primitive.ObjectID{}
	for _, val := range documentsIdsListObject.ListOfIds {

		// Convert the ID to oid
		oID, err := primitive.ObjectIDFromHex(val)
		if err != nil {
			utils.HandleError(ctx, http.StatusNotFound, "cannot convert hex id to bson id", err.Error())
			return
		}

		// Add it to the list of OIDs
		objectIDs = append(objectIDs, oID)
	}

	// Delete multiple users in the database
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)
	deleteMultipleUsersResult, err := collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": objectIDs}})

	if err != nil || deleteMultipleUsersResult.DeletedCount < 0 {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error deleting multiple users.", err.Error())
		return
	}

	// Send the response with the deleted users count
	// The GIN package will automatically encode the response in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"message":          "Deleted multiple documents successfully. Deleted {" + fmt.Sprint(deleteMultipleUsersResult.DeletedCount) + "} Documents for Collection { " + collectionName + " }.",
		"documentsDeleted": deleteMultipleUsersResult.DeletedCount,
	})
}

// DELETE ALL DOCUMENTS ----
// -------------------------
func DeleteAllDocuments[T any](ctx *gin.Context) {

	// Delete all the users
	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Delete all users in the database
	collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)
	deleteAllUsersResult, err := collection.DeleteMany(context.TODO(), bson.M{})

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Check the deleteCount field
	if deleteAllUsersResult.DeletedCount >= 0 {
		// Send the response with the deleted users count
		// The GIN package will automatically encode the response in JSON format
		ctx.JSON(http.StatusOK, gin.H{
			"message":          "Deleted all documents successfully. Deleted {" + fmt.Sprint(deleteAllUsersResult.DeletedCount) + "} Documents for Collection { " + collectionName + " }.",
			"documentsDeleted": deleteAllUsersResult.DeletedCount,
		})

	} else {
		// Handle the error
		utils.HandleError(ctx, http.StatusMethodNotAllowed, "Deletion is not allowed for this user.", errors.New("error happened in deleting all the documents").Error())
		return
	}
}

// UPDATE DOCUMENT ----------
// --------------------------
func UpdateDocument[T any](ctx *gin.Context) {

	// Retrieve the user id from the parameters
	documentID := ctx.Param("id")
	if documentID == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "Error finding document ID.", errors.New("error finding document ID").Error())
		return
	}

	// Retrieve and read the request body
	requestBody := make(map[string]interface{})

	decoderData := json.NewDecoder(ctx.Request.Body)

	// Cast all integers or floats as Number
	decoderData.UseNumber()

	// Decode the data
	err := decoderData.Decode(&requestBody)
	fmt.Println("Request body:", requestBody)
	fmt.Printf("Type:%T\n", requestBody)

	// Check if error OR the request map is empty
	if err != nil || len(requestBody) == 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing and decoding document data.", err.Error())
		return
	}

	// Update the specific document
	// Convert the ID to oid
	oID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id", err.Error())
		return
	}

	// Set the filter for _id and the last updated date
	filter := bson.M{"_id": oID}
	requestBody["lastUpdatedDt"] = time.Now().Format("2006-01-02 15:04:05")

	// 'Cast' the request body to bson.M Map
	update := bson.M{"$set": requestBody}

	// Get the type to retrieve
	var docRetrieve T
	fmt.Printf("Models Type: %T\n", docRetrieve)

	// Set the collection
	collectionName, err := modelToCollectionName(fmt.Sprintf("%T", docRetrieve))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	// Check if provided password in the request body
	if requestBody["password"] != nil && utils.CheckStringNotEmpty(requestBody["password"].(string)) {
		utils.HandleError(ctx, http.StatusBadRequest, "field 'password' in not allowed in the request body.", errors.New("not allowed editing of field 'password'").Error())
		return
	}

	// Role - ONLY FOR "users" COLLECTION
	if collectionName == db.DB_TABLE_USERS && requestBody["role"] != nil && utils.CheckStringNotEmpty(requestBody["role"].(string)) {
		// Allowed roles for the users
		if !utils.CheckAllowedRole(requestBody["role"].(string)) {
			utils.HandleError(ctx, http.StatusBadRequest, "the provided role is not supported", errors.New("the provided role is not supported").Error())
			return
		}
	}

	if len(update) > 0 {

		// Execute the statement
		collection := db.MongoClient.Database(db.DB_NAME).Collection(collectionName)
		result, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "update of the user data failed", err.Error())
			return
		}

		// All successful
		fmt.Println("Update successful: ", result)

		// Return response
		ctx.JSON(http.StatusOK, gin.H{
			"message":          "Document updated successfully.",
			"documentsUpdated": result.ModifiedCount,
		})

	} else {
		utils.HandleError(ctx, http.StatusInternalServerError, "no available data to update", errors.New("no available data to update").Error())
		return
	}
}
