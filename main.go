package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	var base64Data = flag.String("data", "", "Base64 encoded Google Authenticator export data")
	flag.Parse()

	if *base64Data == "" {
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
	writer.Write([]string{"Title", "Username", "Password", "TOTP Secret"})

	for _, otp := range payload.OtpParameters {
		title := otp.GetIssuer()
		username := otp.GetName()
		secret := base32.StdEncoding.EncodeToString(otp.GetSecret())
		writer.Write([]string{title, username, "", secret})
	}
}
