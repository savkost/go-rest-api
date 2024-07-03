package db

import "gopkg.in/mgo.v2/bson"

// This method creates the 'users' schema
func CreateUsersSchema() bson.M {

	// Create a "users" collection with a JSON schema validator
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"firstName", "lastName", "role", "isAdmin", "isActive", "email", "password", "createdDt", "lastUpdatedDt"},
		"properties": bson.M{
			"firstName": bson.M{
				"bsonType":    "string",
				"description": "the first name of the user, which is required and must be a string",
			},
			"lastName": bson.M{
				"bsonType":    "string",
				"description": "the last name of the user, which is required and must be a string",
			},
			"role": bson.M{
				"bsonType":    "string",
				"description": "the role of the user, which is required and must be a string",
			},
			"isAdmin": bson.M{
				"bsonType":    "string",
				"description": "A 0 | 1 value indicating whether the user is admin or not, which is required and must be a string",
			},
			"isActive": bson.M{
				"bsonType":    "string",
				"description": "A 0 | 1 value indicating whether the user is active or not, which is required and must be a string",
			},
			"email": bson.M{
				"bsonType":    "string",
				"description": "the email of the user, which is required and must be a string",
			},
			"password": bson.M{
				"bsonType":    "string",
				"description": "the password of the user, which is required and must be a string",
			},
			"createdDt": bson.M{
				"bsonType":    "string",
				"description": "the created date of the user, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
			"lastUpdatedDt": bson.M{
				"bsonType":    "string",
				"description": "the last updated date of the user, which is required and must be a string (yyyy-MM-dd HH:mm:ss)",
			},
		},
	}

	// Return the schema
	return jsonSchema
}
