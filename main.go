package main

import (
	"go-essentials/go-mongodb-rest-api/db"
	"go-essentials/go-mongodb-rest-api/routes"

	"github.com/gin-gonic/gin"
)

// MAIN FUNC -----------------------------
// ---------------------------------------
func main() {

	// Create and initialize the MongoDB database
	db.InitDB()

	// Create and initialize the pre configured SERVER
	server := gin.Default()

	// Routing and Handling
	routes.RegisterRoutes(server)

	// Start the server in order to listen for incoming requests
	// localhost + :8082 (PORT) -> Development
	server.Run(":8082")
}
