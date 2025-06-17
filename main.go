package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"google.golang.org/protobuf/proto"
)

func readQRFromImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

func main() {
	var base64Data = flag.String("data", "", "Base64 encoded Google Authenticator export data")
	var imagePath = flag.String("image", "", "Path to QR code image file")
	flag.Parse()

	if *imagePath != "" {
		qrText, err := readQRFromImage(*imagePath)
		if err != nil {
			log.Fatalf("Failed to read QR code from image: %v", err)
		}
		
		// Extract base64 data from QR code URL
		if strings.HasPrefix(qrText, "otpauth-migration://offline?data=") {
			*base64Data = strings.TrimPrefix(qrText, "otpauth-migration://offline?data=")
		} else {
			log.Fatalf("Invalid QR code format. Expected otpauth-migration URL, got: %s", qrText)
		}
	} else if *base64Data == "" {
		fmt.Print("Paste Base64 string: ")
		fmt.Scanln(base64Data)
	}

	// URL decode first in case it's URL encoded
	decoded, err := url.QueryUnescape(*base64Data)
	if err != nil {
		decoded = *base64Data // fallback to original if URL decode fails
	}

	data, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		log.Fatalf("Base64 decode failed: %v", err)
	}

	var payload MigrationPayload
	if err := proto.Unmarshal(data, &payload); err != nil {
		log.Fatalf("Protobuf parse failed: %v", err)
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// CSV header for KeePassXC
	writer.Write([]string{"Title", "Username", "Password", "TOTP"})

	for _, otp := range payload.OtpParameters {
		title := otp.GetIssuer()
		username := otp.GetName()
		secret := base32.StdEncoding.EncodeToString(otp.GetSecret())
		writer.Write([]string{title, username, "", secret})
	}
}
