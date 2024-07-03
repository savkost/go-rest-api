package routes

import (
	"go-essentials/go-mongodb-rest-api/controllers"
	"go-essentials/go-mongodb-rest-api/middleware"
	"go-essentials/go-mongodb-rest-api/models"

	"github.com/gin-gonic/gin"
)

// ROUTING PATHS -------------------------
// ---------------------------------------

// GENERAL HELPFUL
const ID_OBJECT_BASIC = ":id"
const COUNT_URL = "count"
const DELETE_ALL_DOCUMENTS = "deleteAll"
const DELETE_MULTIPLE_DOCUMENTS = "deleteMultiple"
const GET_MOST_RECENT = "getMostRecent"

// USERS AND AUTHENTICATION
const SIGNUP_URL = "/register"
const LOGIN_URL = "/login"
const USERS_BASIC_URL = "/users/"
const CREATE_ADMIN_USER = "createAdmin"

// LICENSES CATEGORIES
const LICENSES_CATEGORIES_BASIC_URL = "/categoriesLicenses/"

// LICENSES
const LICENSES_BASIC_URL = "/licenses/"
const RENEW_LICENSE_URL = "renew/"
const UPGRADE_LICENSE_URL = "upgrade/"
const COUNT_PER_CATEGORY_URL = "countPerCategory"

// This method registers all possible and supported routes
func RegisterRoutes(server *gin.Engine) {

	// Not protected routes
	server.POST(SIGNUP_URL, controllers.Register)
	server.POST(LOGIN_URL, controllers.Login)

	// Protected routes
	protectedEventsRoutes := server.Group("/").Use(middleware.Authenticate)

	// PROTECTED ROUTES FOR GENERAL USAGE -------------------------------------------------
	// ------------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------------

	// USERS -------------------------------------------------
	// -------------------------------------------------------

	// Protected routes
	protectedEventsRoutes.GET(USERS_BASIC_URL+ID_OBJECT_BASIC, controllers.GetDocumentByID[models.User])
	protectedEventsRoutes.PATCH(USERS_BASIC_URL+ID_OBJECT_BASIC, controllers.UpdateDocument[models.User])
	protectedEventsRoutes.DELETE(USERS_BASIC_URL+ID_OBJECT_BASIC, controllers.DeleteDocument[models.User])

	// LICENSES CATEGORIES -----------------------------------
	// -------------------------------------------------------
	protectedEventsRoutes.GET(LICENSES_CATEGORIES_BASIC_URL, controllers.GetAllDocuments[models.LicenseCategory])
	protectedEventsRoutes.GET(LICENSES_CATEGORIES_BASIC_URL+GET_MOST_RECENT, controllers.GetLastXDocuments[models.LicenseCategory])
	protectedEventsRoutes.GET(LICENSES_CATEGORIES_BASIC_URL+ID_OBJECT_BASIC, controllers.GetDocumentByID[models.LicenseCategory])

	// LICENSES ----------------------------------------------
	// -------------------------------------------------------
	protectedEventsRoutes.POST(LICENSES_BASIC_URL, controllers.CreateLicense)
	protectedEventsRoutes.GET(LICENSES_BASIC_URL+ID_OBJECT_BASIC, controllers.GetDocumentByID[models.License])
	protectedEventsRoutes.PATCH(LICENSES_BASIC_URL+ID_OBJECT_BASIC, controllers.UpdateDocument[models.License])
	protectedEventsRoutes.DELETE(LICENSES_BASIC_URL+ID_OBJECT_BASIC, controllers.DeleteDocument[models.License])
	protectedEventsRoutes.PATCH(LICENSES_BASIC_URL+RENEW_LICENSE_URL+ID_OBJECT_BASIC, controllers.RenewLicense)
	protectedEventsRoutes.PATCH(LICENSES_BASIC_URL+UPGRADE_LICENSE_URL+ID_OBJECT_BASIC, controllers.UpgradeLicense)

	// PROTECTED ROUTES - FOR ADMINS|SUPERADMINS ONLY -------------------------------------
	// ------------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------------

	// Admin || SuperAdmin Routes - USERS
	protectedEventsRoutes.Use(middleware.CheckAdminUser).POST(USERS_BASIC_URL+CREATE_ADMIN_USER, controllers.CreateAdmin[models.User])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(USERS_BASIC_URL, controllers.GetAllDocuments[models.User])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(USERS_BASIC_URL+GET_MOST_RECENT, controllers.GetLastXDocuments[models.User])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(USERS_BASIC_URL+COUNT_URL, controllers.CountAllDocuments[models.User])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).DELETE(USERS_BASIC_URL+DELETE_ALL_DOCUMENTS, controllers.DeleteAllDocuments[models.User])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).POST(USERS_BASIC_URL+DELETE_MULTIPLE_DOCUMENTS, controllers.DeleteMultipleDocuments[models.User])

	// Admin || SuperAdmin Routes - LICENSE CATEGORIES
	protectedEventsRoutes.Use(middleware.CheckAdminUser).POST(LICENSES_CATEGORIES_BASIC_URL, controllers.CreateDocument[models.LicenseCategory])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).PATCH(LICENSES_CATEGORIES_BASIC_URL+ID_OBJECT_BASIC, controllers.UpdateDocument[models.LicenseCategory])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).DELETE(LICENSES_CATEGORIES_BASIC_URL+ID_OBJECT_BASIC, controllers.DeleteDocument[models.LicenseCategory])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(LICENSES_CATEGORIES_BASIC_URL+COUNT_URL, controllers.CountAllDocuments[models.LicenseCategory])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).DELETE(LICENSES_CATEGORIES_BASIC_URL+DELETE_ALL_DOCUMENTS, controllers.DeleteAllDocuments[models.LicenseCategory])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).POST(LICENSES_CATEGORIES_BASIC_URL+DELETE_MULTIPLE_DOCUMENTS, controllers.DeleteMultipleDocuments[models.LicenseCategory])

	// Admin || SuperAdmin Routes - LICENSES
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(LICENSES_BASIC_URL, controllers.GetAllDocuments[models.License])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(LICENSES_BASIC_URL+GET_MOST_RECENT, controllers.GetLastXDocuments[models.License])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(LICENSES_BASIC_URL+COUNT_URL, controllers.CountAllDocuments[models.License])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).GET(LICENSES_BASIC_URL+COUNT_PER_CATEGORY_URL, controllers.CountLicensesPerCategory)
	protectedEventsRoutes.Use(middleware.CheckAdminUser).DELETE(LICENSES_BASIC_URL+DELETE_ALL_DOCUMENTS, controllers.DeleteAllDocuments[models.License])
	protectedEventsRoutes.Use(middleware.CheckAdminUser).POST(LICENSES_BASIC_URL+DELETE_MULTIPLE_DOCUMENTS, controllers.DeleteMultipleDocuments[models.License])
}
