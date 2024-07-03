package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// DB NAME: test-GO-DB
const DB_NAME = "test-GO-DB"

// TABLE NAMES
const DB_TABLE_USERS = "users"
const DB_TABLE_LICENSES_CATEGORIES = "licensecategories"
const DB_TABLE_LICENCES = "licenses"

// Global Database Client for all the DB Interactions
var MongoClient *mongo.Client

// This method initializes a connection to the MONGO DB and returns a new session
func InitDB() *mongo.Client {

	// Load the DOTENV data
	err := godotenv.Load()
	if err != nil {
		panic("Could not load DOTENV data: " + err.Error())
	}

	// Fetch the MONGO DB connection
	MONGO_DB_CONNECTION_PASSWORD := os.Getenv("MONGO_DB_CONNECTION_PASSWORD")
	MONGO_DB_STRING_1 := os.Getenv("MONGO_DB_CONNECTION_STRING_1")
	MONGO_DB_STRING_2 := os.Getenv("MONGO_DB_CONNECTION_STRING_2")

	// Connect to our CLOUD mongodb
	// DB: test-GO-DB
	MONGODBG_URI := MONGO_DB_STRING_1 + MONGO_DB_CONNECTION_PASSWORD + MONGO_DB_STRING_2
	MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODBG_URI))

	// Check if connection is established
	if err != nil {
		panic("Could not connect to the database: " + err.Error())
	}

	// Retrieve all the collection names in the database
	cNames, err := MongoClient.Database(DB_NAME).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		panic("Could not list all the collections in the database: " + err.Error())
	}

	// Create all the necessary collections
	err = CreateAllCollections(cNames)
	if err != nil {
		panic("Could not create all necessary collections in the database: " + err.Error())
	}

	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Connection with MongoDB successful")
	fmt.Println("Collections in DB: ", cNames)
	fmt.Println("----------------------------------------------------------------")

	// Return the mongodb client
	return MongoClient
}

// This method creates all the necessary collections in the database
func CreateAllCollections(alreadyExistentCollectionNames []string) error {

	// 1. Create 'users' Collection
	err := CreateUsersCollection(alreadyExistentCollectionNames)

	if err != nil {
		fmt.Println("Error creating 'users' collection: ", err)
		return err
	}

	// 2. Create 'licensecategories' Collection
	err = CreateLicenseCategoriesCollection(alreadyExistentCollectionNames)

	if err != nil {
		fmt.Println("Error creating 'licensecategories' collection: ", err)
		return err
	}

	// 3. Create 'licenses' Collection
	err = CreateLicensesCollection(alreadyExistentCollectionNames)

	if err != nil {
		fmt.Println("Error creating 'licenses' collection: ", err)
		return err
	}

	// All operations successfull
	return nil
}

// This method creates the 'users' collection in the database
func CreateUsersCollection(alreadyExistentCollectionNames []string) error {

	// Check if 'users' collection is already created and if YES then return nil error
	for _, value := range alreadyExistentCollectionNames {
		if value == DB_TABLE_USERS {
			return nil
		}
	}

	// Create a "users" collection with a JSON schema validator
	jsonSchema := CreateUsersSchema()

	// Set the validator
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)

	// Create the collection
	err := MongoClient.Database(DB_NAME).CreateCollection(context.TODO(), DB_TABLE_USERS, opts)
	if err != nil {
		return err
	}

	// Create the unique 'email' index on the collection
	collection := MongoClient.Database(DB_NAME).Collection(DB_TABLE_USERS)
	keysMap := make(map[string]int64, 1)
	keysMap["email"] = 1
	indexModel := mongo.IndexModel{
		Keys:    keysMap,
		Options: options.Index().SetUnique(true),
	}

	ind, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	fmt.Println("Index on email created:", ind)
	return err
}

// This method creates the 'licensecategories' collection in the database
func CreateLicenseCategoriesCollection(alreadyExistentCollectionNames []string) error {

	// Check if 'licensecategories' collection is already created and if YES then return nil error
	for _, value := range alreadyExistentCollectionNames {
		if value == DB_TABLE_LICENSES_CATEGORIES {
			return nil
		}
	}

	// Create a "licensecategories" collection with a JSON schema validator
	jsonSchema := CreateLicenseCategoriesSchema()

	// Set the validator
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)

	// Create the collection
	err := MongoClient.Database(DB_NAME).CreateCollection(context.TODO(), DB_TABLE_LICENSES_CATEGORIES, opts)
	if err != nil {
		return err
	}

	// Create the unique 'categoryType' index on the collection
	collection := MongoClient.Database(DB_NAME).Collection(DB_TABLE_LICENSES_CATEGORIES)
	keysMap := make(map[string]int64, 1)
	keysMap["categoryType"] = 1
	indexModel := mongo.IndexModel{
		Keys:    keysMap,
		Options: options.Index().SetUnique(true),
	}

	ind, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	fmt.Println("Index on categoryType created:", ind)
	return err
}

// This method creates the 'licenses' collection in the database
func CreateLicensesCollection(alreadyExistentCollectionNames []string) error {

	// Check if 'licenses' collection is already created and if YES then return nil error
	for _, value := range alreadyExistentCollectionNames {
		if value == DB_TABLE_LICENCES {
			return nil
		}
	}

	// Create a "licenses" collection with a JSON schema validator
	jsonSchema := CreateLicensesSchema()

	// Set the validator
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)
	opts.SetValidationAction("error")

	// Create the collection
	err := MongoClient.Database(DB_NAME).CreateCollection(context.TODO(), DB_TABLE_LICENCES, opts)
	if err != nil {
		return err
	}

	// Create the unique 'licenseKey' index on the collection
	collection := MongoClient.Database(DB_NAME).Collection(DB_TABLE_LICENCES)
	keysMap := make(map[string]int64, 1)
	keysMap["licenseKey"] = 1
	indexModel := mongo.IndexModel{
		Keys:    keysMap,
		Options: options.Index().SetUnique(true),
	}

	ind, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	fmt.Println("Index on licenseKey created:", ind)
	return err
}
