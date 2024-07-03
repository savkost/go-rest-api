package middleware

import (
	"encoding/json"
	"errors"
	"go-essentials/go-mongodb-rest-api/models"
	"go-essentials/go-mongodb-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// This method is a MIDDLEWARE and checks if the request user in an ADMIN or SUPERADMIN user in order to
// continue with the request
func CheckAdminUser(ctx *gin.Context) {

	// Retrieve the user from the context request
	specUserJsonData := ctx.GetString("user")
	if specUserJsonData == "" {
		// Not existent user
		// ABORT NOW the current request
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized user.", errors.New("not authorized user").Error())
		return
	}

	// Unmarshal the JSON data
	var specUser models.User
	err := json.Unmarshal([]byte(specUserJsonData), &specUser)
	if err != nil {
		// ABORT NOW the current request
		utils.HandleError(ctx, http.StatusUnauthorized, "Error unmarshaling the user data.", err.Error())
		return
	}

	// Checking the user role and admin capabilities
	if utils.CheckStringNotEmpty(specUser.IsAdmin) && utils.CheckStringNotEmpty(specUser.Role) {
		if utils.CheckStringNotEmpty(specUser.IsActive) && specUser.IsActive == "1" {
			if specUser.IsAdmin == "1" && (specUser.Role == "admin" || specUser.Role == "superadmin") {
				// The user is an admin or a superadmin
				// Call the NEXT function to continue with the request
				ctx.Next()
				return
			}
		}
	}

	// The user is not an admin
	utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized user.", errors.New("not authorized user").Error())
}
