# go-rest-api

go-rest-api is a complete API for three entities, Users, Licenses and Licenses Categories. This API supports user authentication and authorization with the generation of the appropriate JWT tokens and user roles. 

If you find the plugin helpful, please consider [Supporting the project](https://github.com/sponsors/savkost).

<img src="https://github.com/savkost/go-rest-api/blob/main/screenshots/SKA Logo – 1.png" alt="Screenshot" width="238px" style="max-width: 100%" />

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Authentication, Authorization](#auth-authorize)
- [Generics for Models API](#generic-api-models)
- [MongoDB as Data Storage](#mongo-db)
- [CRUD and Various API Operations](#crud-and-actions)
- [Encryption Methods and Hashing](#encr-hashing)
- [Paging, Filtering, Max Rows](#paging-filter-rows)
- [QRCode Generation](#qrcode-license)
- [Global Error Handling](#glob-error-handling)
- [API Testing](#api-testing)
- [Further Help, Links, LinkedIn](#links)

## Features

- Complete API for three entities, Users, Licenses and Licenses Categories.
- Authentication and Authorization with JWT and roles.
- MongoDB for Data storage.
- CRUD operations for all models. Including count, delete many and delete all.
- QR Code Generation for every license creation.
- Complete testing API.
- Encryption methods and Hashing.

## Getting Started

You can run and initiate the server on localhost by typing the following inside a terminal (on the master directory):

```go
go run main.go
```

The on the terminal all the incoming requests will start to show up as they arrive.

## Authentication, Authorization

The API supports authentication and authorization through JWT and roles. The hashing of the user passwords and the comparison of them is implemented with the use of bcrypt library (golang.org/x/crypto/bcrypt). Furthermore, some APIs and routes and only accessed if the current user possesses specific roles.

## Generics for Models API

This API supports Generics for all the major actions on the models. By setting the generic type [T any] we can then implement all basic routes and actions such as CRUD for any given model. In conjuction to this, the API sets a standard and constant form of responses in order to be more friendly to the developer and user. As shown in the example below, let's see together how we can count the documents of any collection of the system.

```go
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
```

With the use of modelToCollectionName() we can retrieve the model name of the required model and then we can implement any action.

## MongoDB as Data Storage

As previously mentioned, this project and API uses the MongoDB as the data storage for persistent storage of data. Specifically, the API uses the `go.mongodb.org/mongo-driver/mongo` driver. In the mongo-db.go file all the initial actions take place, such as database connection and creation of all the collections.

## CRUD and Various API Operations

For each model (collection) in the system, the API supports the following actions:
  - Create document
  - Update document
  - Delete document
  - Count documents of the collection
  - Retrieve all document of the collection
  - Filtered search of the documents (filters, page, size).
  - Retrieve document by id.
  - Get last X documents of the collection
  - Delete multiple documents by providing a list of IDs.
  - Delete all documents from the collection.

Simultaneously, with the use of generics the creation of all the above actions is fairly easy. Setting the routes, the table name, the configuration of the collection and then adding the table name into the modelToCollectionName function and that's it!

## Encryption Methods and Hashing

This API also implements some encryption and hashing methods such as bcrypt, sha256 and sha512.

## Paging, Filtering, Max Rows

This API supports paging, filtering and limiting the number of results. The user and the developer can add `page` and `limit` in the search query in order to apply the referred actions. Then the result matches the following form:

```go
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
```

## QRCode Generation

On every new license, the API generates a new QR Code that includes the license key. The API uses the [QR Code](https://github.com/skip2/go-qrcode) in order to generate the appropriate QR codes. For example in the following statements we can see the generation of a QR Code.

```go
// Local Struct for QR Code Generation
type QRCodeProduct struct {
	Content string `json:"content"`
	Size    int    `json:"size"`
}

// This method generates QR Code
func (qrData QRCodeProduct) GenerateQRCode() (string, error) {

	// Print the given QR data
	fmt.Println("Qr Data Input:", qrData)

	// Generate the QR Code
	// Input 1: the content string in the QR Code
	// Input 2: the error recovery percentage
	// Input 3: size of the QR Code (image width and height the same = square)
	// OUTPUT: byte slice with the bytes of the PNG image holding the QR Code
	qrCodeResult, err := qrcode.Encode(qrData.Content, qrcode.High, qrData.Size)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Transform to base64 encoding
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeResult)
	fmt.Println("Base64 QR Code:", qrCodeBase64)

	// Success create the QR Code
	return qrCodeBase64, nil
}

// Then simply create a qr code by
qrCodeBase64Data, err := qrCodeDataInput.GenerateQRCode()
if err != nil {
  utils.HandleError(ctx, http.StatusInternalServerError, "Error creating the QR code of the license key.", err.Error())
  return
}
```

## Global Error Handling

The API supports global error handling and it is implemented at the `error_handler.go` file. As an example, we can handle an API error as follows:

```go
utils.HandleError(ctx, http.StatusBadRequest, "Error finding license ID.", errors.New("error finding license ID").Error())
```

## API Testing

The `api-test` folder contains a wide variety of ready tests for every API action. Navigate there and simply call them in order to test the API.

## Further Help, Links, LinkedIn

To get more help on the go-rest-api, feel free to send me any questions at: [Savvas Kostoudas](mailto:savkostoudas@gmail.com)

Let's connect on LinkedIn: [LinkedIn](https://www.linkedin.com/in/savvas-kostoudas-6897491b1/)

If you find this API helpful, please consider [Supporting the project](https://github.com/sponsors/savkost).
