package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-essentials/go-mongodb-rest-api/db"
	"go-essentials/go-mongodb-rest-api/models"
	"go-essentials/go-mongodb-rest-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// This method creates a new license
func CreateLicense(ctx *gin.Context) {

	// License Data
	var licenseData models.License
	err := ctx.ShouldBindJSON(&licenseData)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing license data.", err.Error())
		return
	}

	// Print the received license data
	fmt.Println("Received license data:", licenseData)

	// CONTROLS AND VARIOUS LICENSE CHECKS ------------------
	// ------------------------------------------------------

	// Expiration Date and user holder id
	if !utils.CheckStringNotEmpty(licenseData.Expiration_dt) || !utils.CheckStringNotEmpty(licenseData.UserHolderId) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide all the necessary data in order to create a new license.", errors.New("missing license data").Error())
		return
	}

	// User full name and category id
	if !utils.CheckStringNotEmpty(licenseData.UserFullName) || !utils.CheckStringNotEmpty(licenseData.CategoryId) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide all the necessary data in order to create a new license.", errors.New("missing license data").Error())
		return
	}

	// Category type and category title
	if !utils.CheckStringNotEmpty(licenseData.CategoryType) || !utils.CheckStringNotEmpty(licenseData.CategoryTitle) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide all the necessary data in order to create a new license.", errors.New("missing license data").Error())
		return
	}

	numCategoryTypeTransform, err := utils.TransformStringToInteger64(licenseData.CategoryType)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing category type to number.", errors.New("error integer transformation").Error())
		return
	}
	fmt.Println("Category Type: ", numCategoryTypeTransform)

	// Time span type
	if licenseData.TimeSpanType <= 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Check if time span type IN [1, 3, 6, 12]
	acceptTimeSpanType := false
	timeSpanAcceptedValue := []int64{1, 3, 6, 12}
	for _, value := range timeSpanAcceptedValue {
		if value == licenseData.TimeSpanType {
			acceptTimeSpanType = true
			break
		}
	}

	if !acceptTimeSpanType {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Check if the user has already active licenses
	// String id not object id
	oID, err := primitive.ObjectIDFromHex(licenseData.UserHolderId)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Get the type to retrieve
	var userRetrieve models.User

	// Check user existence in the database and find the user by id
	// Retrieve a specific document
	collectionUsers := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_USERS)

	// Remove password from the documents
	removePasswordOption := bson.M{"password": 0}
	err = collectionUsers.FindOne(context.TODO(), bson.M{"_id": oID}, options.FindOne().SetProjection(removePasswordOption)).Decode(&userRetrieve)

	// Check the error
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "cannot retrieve specific user.", err.Error())
		return
	}

	// Check the user exists in the database
	if !utils.CheckStringNotEmpty(userRetrieve.ID) || !utils.CheckStringNotEmpty(userRetrieve.Email) {
		utils.HandleError(ctx, http.StatusNotFound, "cannot retrieve specific user.", errors.New("cannot retrieve specific user").Error())
		return
	}

	// Count all the licenses for this user
	collectionLicenses := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_LICENCES)
	countDocumentsResult, err := collectionLicenses.CountDocuments(context.TODO(), bson.M{"userHolderId": oID})

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error counting the documents in the collection.", err.Error())
		return
	}
	fmt.Println("Number of Licenses for this user:", countDocumentsResult)

	// Check the count
	if countDocumentsResult > 0 {
		utils.HandleError(ctx, http.StatusInternalServerError, "User already has { "+fmt.Sprintf("%d", countDocumentsResult)+" } licenses.", errors.New("user has licenses").Error())
		return
	}

	// LICENSE KEY AND ACTIVATED ON DEVICE ------------------
	// ------------------------------------------------------

	// Struct Type for data hashing
	type DataForHashing struct {
		Timestamp   string `json:"timestamp"`
		RandomBytes string `json:"random_bytes"`
	}

	// Produce a new UNIQUE license key for this license
	licenseKeyGenerated := ""
	for {

		// Generate secure random bytes (32 bytes - 256 bit)
		secureRandomHexString, err := utils.GenerateSecureRandomBytes(32)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the license key.", err.Error())
			return
		}

		// Create a new license key
		contructedDataForHashing := DataForHashing{
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
			RandomBytes: secureRandomHexString,
		}

		// Marshal the above struct into JSON string
		jsonData, err := json.Marshal(contructedDataForHashing)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the license key.", err.Error())
			return
		}

		licenseKeyGeneratedLocal, err := utils.GenerateSHA512Key(string(jsonData))
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the license key.", err.Error())
			return
		}
		licenseKeyGenerated = licenseKeyGeneratedLocal

		// Check if the license key is UNIQUE inside the LICENSES collection
		countDocumentsWithSameLicenseKeyResult, err := collectionLicenses.CountDocuments(context.TODO(), bson.M{"licenseKey": licenseKeyGenerated})

		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error counting licenses with same license key in the collection.", err.Error())
			return
		}
		fmt.Println("Number of Licenses with the same license key:", countDocumentsWithSameLicenseKeyResult)

		// Check the count of licenses with the same license key
		if countDocumentsResult <= 0 {
			break
		}
	}

	// Produce a new Secure Random HASH for the activatedOnDevice tag of the specific license
	activateOnDeviceGenerated := ""
	for {

		// Generate secure random bytes (32 bytes - 256 bit)
		secureRandomHexString, err := utils.GenerateSecureRandomBytes(32)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the activatedOnDevice value.", err.Error())
			return
		}

		// Create a new hash key
		contructedDataForHashing := DataForHashing{
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
			RandomBytes: secureRandomHexString,
		}

		// Marshal the above struct into JSON string
		jsonData, err := json.Marshal(contructedDataForHashing)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the activatedOnDevice value.", err.Error())
			return
		}

		activatedOnDeviceGeneratedLocal, err := utils.GenerateSHA512Key(string(jsonData))
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error generating the activatedOnDevice value.", err.Error())
			return
		}
		activateOnDeviceGenerated = activatedOnDeviceGeneratedLocal

		// Check if the activatedOnDevice value is UNIQUE inside the LICENSES collection
		countDocumentsWithSameActivatedOnDeviceResult, err := collectionLicenses.CountDocuments(context.TODO(), bson.M{"activatedOnDevice": activateOnDeviceGenerated})

		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error counting licenses with same activatedOnDevice value in the collection.", err.Error())
			return
		}
		fmt.Println("Number of Licenses with the same activatedOnDevice value:", countDocumentsWithSameActivatedOnDeviceResult)

		// Check the count of licenses with the same activatedOnDevice value
		if countDocumentsWithSameActivatedOnDeviceResult <= 0 {
			break
		}
	}

	// Set the appropriate values on the license data
	licenseData.ActivatedOnDevice = activateOnDeviceGenerated
	licenseData.LicenseKey = licenseKeyGenerated

	// Created dt and last update dt (YYYY-MM-DD HH:MM:SS)
	NOW_TIME := time.Now().Format("2006-01-02 15:04:05")
	licenseData.CreatedDt = NOW_TIME
	licenseData.LastUpdatedDt = NOW_TIME
	licenseData.IsActive = "1"
	licenseData.IsExpired = "0"

	// Transform the category id to object id
	oCategoryID, err := utils.StringIDtoObjectID(licenseData.CategoryId)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error transforming the category ID to object ID.", err.Error())
		return
	}

	// Transform the user holder id to object id
	oUserHolderID, err := utils.StringIDtoObjectID(licenseData.UserHolderId)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error transforming the user holder ID to object ID.", err.Error())
		return
	}

	// Cast the license data for insertion
	type LicenseDataForInsertion struct {
		ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
		LicenseKey        string             `bson:"licenseKey,omitempty" json:"licenseKey"`
		Begin_dt          string             `bson:"begin_dt,omitempty" json:"begin_dt"`
		Expiration_dt     string             `bson:"expiration_dt,omitempty" json:"expiration_dt"`
		UserHolderId      primitive.ObjectID `bson:"userHolderId,omitempty" json:"userHolderId"`
		UserFullName      string             `bson:"userFullName,omitempty" json:"userFullName"`
		CategoryId        primitive.ObjectID `bson:"categoryId,omitempty" json:"categoryId"`
		CategoryType      string             `bson:"categoryType,omitempty" json:"categoryType"`
		CategoryTitle     string             `bson:"categoryTitle,omitempty" json:"categoryTitle"`
		ActivatedOnDevice string             `bson:"activatedOnDevice,omitempty" json:"activatedOnDevice"`
		TimeSpanType      int64              `bson:"timeSpanType,omitempty" json:"timeSpanType"`
		Comments          string             `bson:"comments,omitempty" json:"comments"`
		IsActive          string             `bson:"isActive,omitempty" json:"isActive"`
		IsExpired         string             `bson:"isExpired,omitempty" json:"isExpired"`
		CreatedDt         string             `bson:"createdDt,omitempty" json:"createdDt"`
		LastUpdatedDt     string             `bson:"lastUpdatedDt,omitempty" json:"lastUpdatedDt"`
	}

	// Set the license data for insertion
	licenseDataInsert := LicenseDataForInsertion{
		LicenseKey:        licenseData.LicenseKey,
		Begin_dt:          NOW_TIME,
		Expiration_dt:     licenseData.Expiration_dt,
		UserHolderId:      oUserHolderID,
		UserFullName:      licenseData.UserFullName,
		CategoryId:        oCategoryID,
		CategoryType:      licenseData.CategoryType,
		CategoryTitle:     licenseData.CategoryTitle,
		ActivatedOnDevice: licenseData.ActivatedOnDevice,
		TimeSpanType:      licenseData.TimeSpanType,
		Comments:          licenseData.Comments,
		IsActive:          licenseData.IsActive,
		IsExpired:         licenseData.IsExpired,
		CreatedDt:         licenseData.CreatedDt,
		LastUpdatedDt:     licenseData.LastUpdatedDt,
	}

	// Print the data to insert
	fmt.Println("License data to insert: ", licenseDataInsert)

	// Insert the license into the database
	result, err := collectionLicenses.InsertOne(context.TODO(), licenseDataInsert)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error storing the new license in the database.", err.Error())
		return
	}

	// Print the insert result
	fmt.Println("Insert Result:", result)

	// CREATE THE QR CODE WITH THE LICENSE KEY --------------
	// ------------------------------------------------------
	qrCodeDataInput := utils.QRCodeProduct{
		Content: licenseData.LicenseKey,
		Size:    256,
	}

	qrCodeBase64Data, err := qrCodeDataInput.GenerateQRCode()
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error creating the QR code of the license key.", err.Error())
		return
	}

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "License inserted successfully in the collection { " + db.DB_TABLE_LICENCES + " }.",
		"data": []map[string]any{
			{
				"_id":        result.InsertedID,
				"licenseKey": licenseKeyGenerated,
				"QRCode":     qrCodeBase64Data,
			},
		},
	})
}

// This method renews an existent license
func RenewLicense(ctx *gin.Context) {

	// Retrieve the license id from the parameters
	licenseID := ctx.Param("id")
	if !utils.CheckStringNotEmpty(licenseID) {
		utils.HandleError(ctx, http.StatusBadRequest, "Error finding license ID.", errors.New("error finding license ID").Error())
		return
	}

	// String id not object id
	licenseIDObject, err := primitive.ObjectIDFromHex(licenseID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Console the body for renewing the license
	// License Data
	var licenseData models.License
	err = ctx.ShouldBindJSON(&licenseData)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing license data.", err.Error())
		return
	}

	// Print the received license data
	fmt.Println("Received license data:", licenseData)

	// CONTROLS AND VARIOUS LICENSE CHECKS ------------------
	// ------------------------------------------------------

	// Expiration Date and time span type
	if !utils.CheckStringNotEmpty(licenseData.Expiration_dt) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide all the necessary data in order to renew the license.", errors.New("missing license data").Error())
		return
	}

	// Time span type
	if licenseData.TimeSpanType <= 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Check if time span type IN [1, 3, 6, 12]
	acceptTimeSpanType := false
	timeSpanAcceptedValue := []int64{1, 3, 6, 12}
	for _, value := range timeSpanAcceptedValue {
		if value == licenseData.TimeSpanType {
			acceptTimeSpanType = true
			break
		}
	}

	if !acceptTimeSpanType {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Set the filter for _id and the last updated date
	filter := bson.M{"_id": licenseIDObject}

	// Set the proper fields for update
	licenseData.IsActive = "1"
	licenseData.IsExpired = "0"
	licenseData.LastUpdatedDt = time.Now().Format("2006-01-02 15:04:05")

	// 'Cast' the request body to bson.M Map
	update := bson.M{"$set": licenseData}

	// Find the license and update the data
	if len(update) > 0 {

		// Execute the statement and return the UPDATED license document after the update
		collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_LICENCES)
		result := models.License{}
		err := collection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "renew of the selected license failed", err.Error())
			return
		}

		// All successful
		fmt.Println("License Renew successful: ", result)

		// Return response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "License renewed successfully.",
			"data":    []models.License{result},
		})

	} else {
		utils.HandleError(ctx, http.StatusInternalServerError, "no available data to renew license", errors.New("no available data to renew").Error())
		return
	}
}

// This method upgrades an existent license
func UpgradeLicense(ctx *gin.Context) {

	// Retrieve the license id from the parameters
	licenseID := ctx.Param("id")
	if !utils.CheckStringNotEmpty(licenseID) {
		utils.HandleError(ctx, http.StatusBadRequest, "Error finding license ID.", errors.New("error finding license ID").Error())
		return
	}

	// String id not object id
	licenseIDObject, err := primitive.ObjectIDFromHex(licenseID)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "cannot convert hex id to bson id.", errors.New("cannot convert hex id to bson id").Error())
		return
	}

	// Console the body for renewing the license
	// License Data
	var licenseData models.License
	err = ctx.ShouldBindJSON(&licenseData)

	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Error parsing license data.", err.Error())
		return
	}

	// Print the received license data
	fmt.Println("Received license data:", licenseData)

	// CONTROLS AND VARIOUS LICENSE CHECKS ------------------
	// ------------------------------------------------------

	// Expiration Date and begin_dt
	if !utils.CheckStringNotEmpty(licenseData.Expiration_dt) || !utils.CheckStringNotEmpty(licenseData.Begin_dt) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide the expiration date and the begin date in order to renew the license.", errors.New("missing license data").Error())
		return
	}

	// categoryId, categoryType and categoryTitle
	if !utils.CheckStringNotEmpty(licenseData.CategoryId) || !utils.CheckStringNotEmpty(licenseData.CategoryType) || !utils.CheckStringNotEmpty(licenseData.CategoryTitle) {
		utils.HandleError(ctx, http.StatusBadRequest, "Please provide the category data (id, type, title) in order to renew the license.", errors.New("missing license data").Error())
		return
	}

	// Time span type
	if licenseData.TimeSpanType <= 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Check if time span type IN [1, 3, 6, 12]
	acceptTimeSpanType := false
	timeSpanAcceptedValue := []int64{1, 3, 6, 12}
	for _, value := range timeSpanAcceptedValue {
		if value == licenseData.TimeSpanType {
			acceptTimeSpanType = true
			break
		}
	}

	if !acceptTimeSpanType {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid TimeSpanType field value", errors.New("invalid value").Error())
		return
	}

	// Set the filter for _id and the last updated date
	filter := bson.M{"_id": licenseIDObject}

	// Set the proper fields for upgrade
	licenseData.IsActive = "1"
	licenseData.IsExpired = "0"
	licenseData.LastUpdatedDt = time.Now().Format("2006-01-02 15:04:05")

	// Transform the category id to object id
	oCategoryID, err := utils.StringIDtoObjectID(licenseData.CategoryId)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error transforming the category ID to object ID.", err.Error())
		return
	}

	// Cast the license data for upgrade
	type LicenseDataForUpgrade struct {
		Begin_dt      string             `bson:"begin_dt,omitempty" json:"begin_dt"`
		Expiration_dt string             `bson:"expiration_dt,omitempty" json:"expiration_dt"`
		CategoryId    primitive.ObjectID `bson:"categoryId,omitempty" json:"categoryId"`
		CategoryType  string             `bson:"categoryType,omitempty" json:"categoryType"`
		CategoryTitle string             `bson:"categoryTitle,omitempty" json:"categoryTitle"`
		TimeSpanType  int64              `bson:"timeSpanType,omitempty" json:"timeSpanType"`
		IsActive      string             `bson:"isActive,omitempty" json:"isActive"`
		IsExpired     string             `bson:"isExpired,omitempty" json:"isExpired"`
		LastUpdatedDt string             `bson:"lastUpdatedDt,omitempty" json:"lastUpdatedDt"`
	}

	// Set the license data for upgrade
	licenseDataUpgrade := LicenseDataForUpgrade{
		Begin_dt:      licenseData.Begin_dt,
		Expiration_dt: licenseData.Expiration_dt,
		CategoryId:    oCategoryID,
		CategoryType:  licenseData.CategoryType,
		CategoryTitle: licenseData.CategoryTitle,
		TimeSpanType:  licenseData.TimeSpanType,
		IsActive:      licenseData.IsActive,
		IsExpired:     licenseData.IsExpired,
		LastUpdatedDt: licenseData.LastUpdatedDt,
	}

	// Print the data to insert
	fmt.Println("License data to upgrade: ", licenseDataUpgrade)

	// 'Cast' the request body to bson.M Map
	upgradeStatement := bson.M{"$set": licenseDataUpgrade}

	// Find the license and upgrade the data
	if len(upgradeStatement) > 0 {

		// Execute the statement and return the UPGRADED license document after the upgrade
		collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_LICENCES)
		result := models.License{}
		err := collection.FindOneAndUpdate(context.TODO(), filter, upgradeStatement, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "upgrade of the selected license failed", err.Error())
			return
		}

		// All successful
		fmt.Println("License Upgrade successful: ", result)

		// Return response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "License upgraded successfully.",
			"data":    []models.License{result},
		})

	} else {
		utils.HandleError(ctx, http.StatusInternalServerError, "no available data to upgrade license", errors.New("no available data to renew").Error())
		return
	}
}

// This method counts licenses per category
// (Optional) FROM - TO Begin_dt parameters
func CountLicensesPerCategory(ctx *gin.Context) {

	// Console the Query
	fmt.Println("Request Query:", ctx.Request.URL)

	// Created Date From
	createdDateFrom := ctx.Query("created_dtFrom")
	fmt.Println("Created Date From:", createdDateFrom)

	// Created Date To
	createdDateTo := ctx.Query("created_dtTo")
	fmt.Println("Created Date To:", createdDateTo)

	collection := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_LICENCES)

	// Search with the given filters
	// If not given any filters, then the API returns the counting for all licenses
	// Construct the filter query
	matchStage := bson.M{}

	// Check all the fields --------
	// -----------------------------

	// Created Date From, To
	if utils.CheckStringNotEmpty(createdDateFrom) && utils.CheckStringNotEmpty(createdDateTo) {
		matchStage = bson.M{"$match": bson.M{"createdDt": bson.M{"$gte": createdDateFrom, "$lte": createdDateTo}}}

	} else {
		if utils.CheckStringNotEmpty(createdDateFrom) {
			matchStage = bson.M{"$match": bson.M{"createdDt": bson.M{"$gte": createdDateFrom}}}
		} else if utils.CheckStringNotEmpty(createdDateTo) {
			matchStage = bson.M{"$match": bson.M{"createdDt": bson.M{"$lte": createdDateFrom}}}
		}
	}

	// Set the group stage
	groupStage := bson.M{
		"$group": bson.M{
			"_id":           "$categoryId",
			"countLicenses": bson.M{"$sum": 1},
		},
	}

	// Set the group aggregation pipe stage
	fmt.Println("Match Stage:", matchStage)
	fmt.Println("Group Stage:", groupStage)

	// Pipeline Filter
	pipelineFilter := []bson.M{}
	if len(matchStage) > 0 {
		pipelineFilter = append(pipelineFilter, matchStage, groupStage)
	} else {
		pipelineFilter = append(pipelineFilter, groupStage)
	}

	fmt.Println("Pipeline Filter:", pipelineFilter)

	// Aggregate and filter
	rows, err := collection.Aggregate(context.TODO(), pipelineFilter)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "error happened in licenses counting per category", err.Error())
		return
	}

	// Return all counting data
	licensesCountingResult := []bson.M{}
	if err = rows.All(context.TODO(), &licensesCountingResult); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error counting licenses per category.", err.Error())
		return
	}

	// All successful
	fmt.Println("License counting per category successful: ", licensesCountingResult)

	// Retrieve all the license categories
	if len(licensesCountingResult) > 0 {
		arrayWithOIDs := []primitive.ObjectID{}
		for _, val := range licensesCountingResult {
			if val["_id"] != nil && val["_id"] != "" {

				// Append the oID in the arrayWithOIDs array
				arrayWithOIDs = append(arrayWithOIDs, val["_id"].(primitive.ObjectID))
			}
		}

		fmt.Println("Array with OIDs: ", arrayWithOIDs)

		filterQuery := bson.M{"_id": bson.M{"$in": arrayWithOIDs}}
		rowsLicenseCategories, err := db.MongoClient.Database(db.DB_NAME).Collection(db.DB_TABLE_LICENSES_CATEGORIES).Find(context.TODO(), filterQuery)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error retrieving license categories.", err.Error())
			return
		}

		// Receive all the categories data and place them in an array
		licenseCategories := []models.LicenseCategory{}
		if err = rowsLicenseCategories.All(context.TODO(), &licenseCategories); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Error retrieving license categories.", err.Error())
			return
		}

		// Decode all the received license categories
		for rowsLicenseCategories.Next(context.TODO()) {

			// Local Variable
			var resLicenseCat models.LicenseCategory

			err = rowsLicenseCategories.Decode(&resLicenseCat)
			if err != nil {
				utils.HandleError(ctx, http.StatusInternalServerError, "Error decoding license categories.", err.Error())
				return
			}
			licenseCategories = append(licenseCategories, resLicenseCat)
		}

		// For every license category found, place the info in the main result
		for _, licenseCategory := range licenseCategories {
			for _, valRes := range licensesCountingResult {
				if valRes["_id"].(primitive.ObjectID).Hex() != "" {
					if valRes["_id"].(primitive.ObjectID).Hex() == licenseCategory.ID {
						valRes["title"] = licenseCategory.Title
						valRes["categoryType"] = licenseCategory.CategoryType
						break
					}
				}
			}
		}
	}

	// Find the total counting of licenses
	var totalLicensesCount int32 = 0
	if len(licensesCountingResult) > 0 {
		for _, value := range licensesCountingResult {
			if value["countLicenses"] != nil {
				totalLicensesCount += value["countLicenses"].(int32)
			}
		}
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "License counting per category successful.",
		"data":          licensesCountingResult,
		"licensesCount": totalLicensesCount,
	})
}
