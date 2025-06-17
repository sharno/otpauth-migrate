# OTP Auth Migrate

Tool to export Google Authenticator OTP entries to CSV format compatible with KeePassXC.

Some of the goals of this tool:
- The code should be very simple to read to make sure you are not running something that might steal your TOTPs
- It should deal with QR images so that you can use the tool directly with the export of Google Authenticator
- Export in CSV format to import easily in KeePassXC
- A cli tool for maximum flexibility if you want to migrate a big number of TOTPs

## Usage

### From QR Code Image
```bash
go run . -image path/to/qrcode.png
```

### From Base64 String
```bash
go run . -data "your_base64_encoded_data"
```

### Interactive Mode
```bash
go run .
# Then paste the base64 string when prompted
```

## Building

```bash
go build -o otpauth
```

## Output

CSV format with columns: Title, Username, Password, TOTP

The output is written to stdout, so you can redirect it to a file:
```bash
go run . -image qrcode.png > output.csv
```

## Supported Image Formats

- PNG
- JPEG

## Requirements

- Go 1.19+
- Google Authenticator export QR code or base64 data 