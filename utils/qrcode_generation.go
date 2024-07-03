package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

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
