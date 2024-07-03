package db

import "gopkg.in/mgo.v2/bson"

// This method creates the 'licenses' schema
func CreateLicensesSchema() bson.M {

	// Create a "licenses" collection with a JSON schema validator
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"licenseKey", "expiration_dt", "begin_dt", "userHolderId", "userFullName", "categoryId", "categoryType", "categoryTitle", "activatedOnDevice", "timeSpanType", "createdDt", "lastUpdatedDt", "isActive", "isExpired"},
		"properties": bson.M{
			"licenseKey": bson.M{
				"bsonType":    "string",
				"description": "the license κey of the license, which is required and must be a string",
			},
			"begin_dt": bson.M{
				"bsonType":    "string",
				"description": "the begin date of the license, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
			"expiration_dt": bson.M{
				"bsonType":    "string",
				"description": "the expiration date of the license, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
			"userHolderId": bson.M{
				"bsonType":    "objectId",
				"description": "the id of the holder of the license, which is required and must be an object ID",
			},
			"userFullName": bson.M{
				"bsonType":    "string",
				"description": "the fullname of the user of the license, which is required and must be a integer",
			},
			"categoryId": bson.M{
				"bsonType":    "objectId",
				"description": "the id of the category of the license, which is required and must be an object ID",
			},
			"categoryType": bson.M{
				"bsonType":    "string",
				"description": "the category type of the license, which is required and must be a string",
			},
			"categoryTitle": bson.M{
				"bsonType":    "string",
				"description": "The title of the license category, which is required and must be a string",
			},
			"activatedOnDevice": bson.M{
				"bsonType":    "string",
				"description": "a key for the activated on device tag, which is required and must be a string",
			},
			"timeSpanType": bson.M{
				"bsonType":    "long",
				"enum":        []int64{1, 3, 6, 12},
				"description": "the time span type of the license, which is required and must be IN the [1, 3, 6, 12] range",
			},
			"createdDt": bson.M{
				"bsonType":    "string",
				"description": "the created date of the license category, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
			"lastUpdatedDt": bson.M{
				"bsonType":    "string",
				"description": "the last updated date of the license category, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
			"comments": bson.M{
				"bsonType":    "string",
				"description": "(Optional) comments of the license category",
			},
			"isActive": bson.M{
				"bsonType":    "string",
				"description": "A string value indicating that the license category is ACTIVE or not, which is required and must be a string",
			},
			"isExpired": bson.M{
				"bsonType":    "string",
				"description": "A string value indicating that the license category is EXPIRED or not, which is required and must be a string",
			},
		},
	}

	// Return the schema
	return jsonSchema
}
