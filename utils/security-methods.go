package utils

import (
	"crypto"
	"encoding/hex"
	"fmt"

	"github.com/theckman/go-securerandom"
	"golang.org/x/crypto/bcrypt"
)

// This method hashes the given password
// Input: Plaintext password
// Output: Hashed password
func HashPassword(password string) (string, error) {

	// Password bytes slice, Cost = 12
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	// Convert the bytes slice to a string and return the hashed password
	hashedPassword := string(hashedPasswordBytes)
	return hashedPassword, nil
}

// This method checks and compares a hashed password against the raw password
func CheckPasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// This method generates secure random bytes
// Input: size of random bytes to create
func GenerateSecureRandomBytes(sizeRandomBytesGenerate int) (string, error) {

	// 1. Create the secure random bytes
	randomBytes, err := securerandom.Bytes(sizeRandomBytesGenerate)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Secure random bytes generated:", randomBytes)

	// 2. Transform the secure random bytes to HEX string
	outputHexBytes := make([]byte, hex.EncodedLen(len(randomBytes)))
	intHexBytes := hex.Encode(outputHexBytes, randomBytes)
	outputHexString := string(outputHexBytes)
	fmt.Println("Secure random bytes HEX STRING:", outputHexString)
	fmt.Println("Size:", intHexBytes)

	// 3. Return the HEX string
	return outputHexString, nil
}

// This method generates a SHA512 key HASH
// Input: data input to hash
func GenerateSHA512Key(dataInputToHash string) (string, error) {

	// 0. Console the given input
	fmt.Println("Input to HASH:", dataInputToHash)

	// Generate the SHA512 HASH
	hashFunc := crypto.SHA512.New()
	_, err := hashFunc.Write([]byte(dataInputToHash))
	if err != nil {
		return "", err
	}

	// Return the output hash
	outputHexHash := make([]byte, hex.EncodedLen(len(hashFunc.Sum(nil))))
	nBytesHash := hex.Encode(outputHexHash, hashFunc.Sum(nil))
	fmt.Println("Output SHA512 HASH:", string(outputHexHash))
	fmt.Println("Size:", nBytesHash)
	return string(outputHexHash), nil
}

// This method generates a SHA256 key HASH
// Input: data input to hash
func GenerateSHA256Key(dataInputToHash string) (string, error) {

	// 0. Console the given input
	fmt.Println("Input to HASH:", dataInputToHash)

	// Generate the SHA256 HASH
	hashFunc := crypto.SHA256.New()
	_, err := hashFunc.Write([]byte(dataInputToHash))
	if err != nil {
		return "", err
	}

	// Return the output hash
	outputHexHash := make([]byte, hex.EncodedLen(len(hashFunc.Sum(nil))))
	nBytesHash := hex.Encode(outputHexHash, hashFunc.Sum(nil))
	fmt.Println("Output SHA256 HASH:", string(outputHexHash))
	fmt.Println("Size:", nBytesHash)
	return string(outputHexHash), nil
}

// This method transposes the input using the Caesar Cryptosystem
func TransposeUsingCaesarCryptosystem(inputToTranspose string, placesNumberTranspose int64, isReverseCaesar bool) (string, error) {

	// Console the input data
	fmt.Println("Input to transpose:", inputToTranspose)

	// Store the output
	outputResult := ""

	// Transpose each character
	for _, valRune := range inputToTranspose {

		// Check if it is isReverseCaesar
		asciiCodeOfChar := valRune
		if isReverseCaesar {
			asciiCodeOfChar -= rune(placesNumberTranspose)
		} else {
			asciiCodeOfChar += rune(placesNumberTranspose)
		}

		// Set the output result
		outputResult += fmt.Sprintf("%c", asciiCodeOfChar)
	}

	fmt.Println("Caesar Cryptosystem Output:", outputResult)
	return outputResult, nil
}
