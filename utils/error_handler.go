package utils

import (
	"github.com/gin-gonic/gin"
)

// This method handles the errors on the API requests
// Inputs:
//  1. context of request
//  2. status of the response
//  3. Any message to add to the error response
//  4. Error string
func HandleError(context *gin.Context, statusReturn int, messageReturn string, errorHappened string) {

	// Handle error and abort this API request
	context.AbortWithStatusJSON(statusReturn, gin.H{
		"message": messageReturn,
		"error":   errorHappened,
	})

}
