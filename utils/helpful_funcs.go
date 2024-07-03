package utils

import (
	"encoding/json"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This method checks if the given string is not nil and not empty
func CheckStringNotEmpty(sGiven string) bool {
	return len(sGiven) > 0
}

// This method checks if the given string is an allowed user role
func CheckAllowedRole(roleGiven string) bool {

	// Allowed roles for the users
	ALLOWED_ROLES := []string{"user", "admin", "superadmin"}

	// Check the role of the user
	foundRole := false
	for _, v := range ALLOWED_ROLES {
		if v == roleGiven {
			foundRole = true
			break
		}
	}

	return foundRole
}

// This method transforms the given string to int64
func TransformStringToInteger64(s string) (int64, error) {

	// Attempt to parse the string
	numTransform, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return numTransform, nil
}

// This method transforms the given string to object ID
func StringIDtoObjectID(documentID string) (primitive.ObjectID, error) {

	// String id to object id
	oID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Return the object ID
	return oID, nil
}

// This method transforms the given data to JSON object
func ConvertToJSON(inputData interface{}) (string, error) {

	// Convert to JSON
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		return "", err
	}

	// Return the JSON data
	return string(jsonData), nil
}
