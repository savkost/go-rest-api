package db

import "gopkg.in/mgo.v2/bson"

// This method creates the 'licensecategories' schema
func CreateLicenseCategoriesSchema() bson.M {

	// Create a "licensecategories" collection with a JSON schema validator
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"title", "description", "priceEurosMonthly", "priceEurosThreeMonths", "priceEurosSixMonths", "priceEurosTwelveMonths", "categoryType", "textsQntAllowed", "imagesQntAllowed", "createdDt", "lastUpdatedDt", "isActive"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"description": "the title of the license category, which is required and must be a string",
			},
			"description": bson.M{
				"bsonType":    "string",
				"description": "the description of the license category, which is required and must be a string",
			},
			"priceEurosMonthly": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the monthly price in euros of the license category, which is required and must be a integer",
			},
			"priceEurosThreeMonths": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the three months price in euros of the license category, which is required and must be a integer",
			},
			"priceEurosSixMonths": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the six months price in euros of the license category, which is required and must be a integer",
			},
			"priceEurosTwelveMonths": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the twelve months price in euros of the license category, which is required and must be a integer",
			},
			"categoryType": bson.M{
				"bsonType":    "string",
				"description": "the type of the license category, which is required and must be a string",
			},
			"isActive": bson.M{
				"bsonType":    "string",
				"description": "A string value indicating that the license category is active or not, which is required and must be a string",
			},
			"textsQntAllowed": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the allowed number of texts of the license category, which is required and must be a integer",
			},
			"imagesQntAllowed": bson.M{
				"bsonType":    "long",
				"minimum":     0,
				"description": "the allowed number of images of the license category, which is required and must be a integer",
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
		},
	}

	// Return the schema
	return jsonSchema
}
